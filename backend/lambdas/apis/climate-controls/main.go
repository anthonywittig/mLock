package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/lambdas/helpers"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/auditlog"
	"mlock/lambdas/shared/dynamo/climatecontrol"
	"mlock/lambdas/shared/dynamo/device"
	"mlock/lambdas/shared/dynamo/miscellaneous"
	"mlock/lambdas/shared/dynamo/unit"
	"mlock/lambdas/shared/sqs"
	mshared "mlock/shared"
	"net/http"
	"regexp"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type ClimateControlEntity struct {
	ClimateControl shared.ClimateControl `json:"climateControl"`
	Unit           shared.Unit           `json:"unit"`
}

type DeleteResponse struct {
	Error string `json:"error"`
}

type DetailResponse struct {
	Entity ClimateControlEntity `json:"entity"`
	Extra  ExtraEntities        `json:"extra"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ExtraEntities struct {
	AuditLog              shared.AuditLog              `json:"auditLog"`
	UnitOccupancyStatuses []shared.UnitOccupancyStatus `json:"unitOccupancyStatuses"`
}

type ListResponse struct {
	Entities                       []ClimateControlEntity        `json:"entities"`
	ClimateControlOccupiedSettings shared.ClimateControlSettings `json:"climateControlOccupiedSettings"`
	ClimateControlVacantSettings   shared.ClimateControlSettings `json:"climateControlVacantSettings"`
}

type UpdateBody struct {
	SyncWithReservations bool `json:"syncWithReservations"`
}

type SettingsUpdateRequest struct {
	ClimateControlOccupiedSettings shared.ClimateControlSettings `json:"climateControlOccupiedSettings"`
	ClimateControlVacantSettings   shared.ClimateControlSettings `json:"climateControlVacantSettings"`
}

var entityRegex = regexp.MustCompile(`/climate-controls/?`)

func main() {
	helpers.StartAPILambda(HandleRequest, []string{helpers.MiddlewareAuth})
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	match, err := regexp.MatchString(`^/climate-controls/settings`, req.Path)
	if err != nil {
		return shared.NewAPIResponse(http.StatusBadRequest, ErrorResponse{Error: "unable to parse request"})
	}
	if match {
		return settingsHandleRequest(ctx, req)
	}

	if id := entityRegex.ReplaceAllString(req.Path, ""); id != "" {
		parsedID, err := uuid.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("error parsing id: %s", err.Error())
		}

		entity, ok, err := climatecontrol.NewRepository().Get(ctx, parsedID)
		if err != nil {
			return nil, fmt.Errorf("error getting entity: %s", err.Error())
		}
		if !ok {
			return nil, fmt.Errorf("entity not found: %s", parsedID)
		}

		switch req.HTTPMethod {
		case "DELETE":
			return delete(ctx, req, entity)
		case "GET":
			return detail(ctx, req, entity)
		case "PUT":
			return update(ctx, req, entity)
		default:
			return shared.NewAPIResponse(http.StatusNotImplemented, "not implemented")
		}

	}

	switch req.HTTPMethod {
	case "GET":
		return list(ctx, req)
	default:
		return shared.NewAPIResponse(http.StatusNotImplemented, "not implemented")
	}
}

func delete(ctx context.Context, req events.APIGatewayProxyRequest, entity shared.ClimateControl) (*shared.APIResponse, error) {
	ok := false

	awhileAgo := time.Now().Add(-2 * time.Hour)
	if entity.LastRefreshedAt.Before(awhileAgo) {
		ok = true
	}

	if !ok {
		return shared.NewAPIResponse(http.StatusBadRequest, DeleteResponse{
			Error: "climate control can't be deleted because it was recently refreshed",
		})
	}

	if err := climatecontrol.NewRepository().Delete(ctx, entity.ID); err != nil {
		return nil, fmt.Errorf("error deleting entity: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, DeleteResponse{})
}

func detail(ctx context.Context, req events.APIGatewayProxyRequest, entity shared.ClimateControl) (*shared.APIResponse, error) {
	auditLog, found, err := auditlog.Get(ctx, entity.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting audit logs: %s", err.Error())
	}
	if !found {
		auditLog = shared.AuditLog{Entries: []shared.AuditLogEntry{}}
	}

	// We could probably kill this off mid-2025 since we limit it on update.
	if len(auditLog.Entries) > 100 {
		auditLog.Entries = auditLog.Entries[len(auditLog.Entries)-100:]
	}

	// Reverse the entries so that the newer items are first.
	// https://github.com/golang/go/wiki/SliceTricks#reversing
	for i := len(auditLog.Entries)/2 - 1; i >= 0; i-- {
		opp := len(auditLog.Entries) - 1 - i
		auditLog.Entries[i], auditLog.Entries[opp] = auditLog.Entries[opp], auditLog.Entries[i]
	}

	units, err := unit.NewRepository().ListByName(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting units: %s", err.Error())
	}
	unit := units[entity.GetFriendlyNamePrefix()]

	devices, err := device.NewRepository().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting devices: %s", err.Error())
	}

	tzName, err := mshared.GetConfig("TIME_ZONE")
	if err != nil {
		return nil, fmt.Errorf("error getting time zone name: %s", err.Error())
	}

	tz, err := time.LoadLocation(tzName)
	if err != nil {
		return nil, fmt.Errorf("error getting time zone %s", err.Error())
	}

	unitOccupancyStatuses := []shared.UnitOccupancyStatus{}
	for i := 0; i < 7; i++ {
		now := time.Now().In(tz).AddDate(0, 0, i)
		year, month, day := now.Date()
		date := time.Date(year, month, day, 0, 0, 0, 0, tz)

		occupiedStatusForDay := unit.OccupancyStatusForDay(devices, date)
		unitOccupancyStatuses = append(unitOccupancyStatuses, occupiedStatusForDay)
	}

	return shared.NewAPIResponse(http.StatusOK, DetailResponse{
		Entity: ClimateControlEntity{
			ClimateControl: entity,
			Unit:           unit,
		},
		Extra: ExtraEntities{
			AuditLog:              auditLog,
			UnitOccupancyStatuses: unitOccupancyStatuses,
		},
	})
}

func list(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	climateControls, err := climatecontrol.NewRepository().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting entities: %s", err.Error())
	}

	miscellaneous, ok, err := miscellaneous.NewRepository().Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting miscellaneous: %s", err.Error())
	}
	if !ok {
		return shared.NewAPIResponse(http.StatusNotFound, "miscellaneous not found")
	}

	units, err := unit.NewRepository().ListByName(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting units: %s", err.Error())
	}

	entities := make([]ClimateControlEntity, 0, len(climateControls))
	for _, climateControl := range climateControls {
		unit, _ := units[climateControl.GetFriendlyNamePrefix()]
		entities = append(entities, ClimateControlEntity{
			ClimateControl: climateControl,
			Unit:           unit,
		})
	}

	return shared.NewAPIResponse(
		http.StatusOK, ListResponse{
			Entities:                       entities,
			ClimateControlOccupiedSettings: miscellaneous.ClimateControlOccupiedSettings,
			ClimateControlVacantSettings:   miscellaneous.ClimateControlVacantSettings,
		})
}

func settingsHandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	switch req.HTTPMethod {
	case "PUT":
		return updateSettings(ctx, req)
	default:
		return shared.NewAPIResponse(http.StatusNotImplemented, "not implemented")
	}
}

func update(ctx context.Context, req events.APIGatewayProxyRequest, entity shared.ClimateControl) (*shared.APIResponse, error) {
	var body UpdateBody
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	ccRepo := climatecontrol.NewRepository()

	entity.SyncWithReservations = body.SyncWithReservations
	if err := ccRepo.AppendToAuditLog(ctx, entity, fmt.Sprintf("Updating sync with reservations to %t", entity.SyncWithReservations)); err != nil {
		return nil, fmt.Errorf("error appending to audit log: %s", err.Error())
	}

	entity, err := ccRepo.Put(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("error putting entity: %s", err.Error())
	}

	return detail(ctx, req, entity)
}

func updateSettings(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	var body SettingsUpdateRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	queue, err := sqs.NewSQSService(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating sqs service: %s", err.Error())
	}

	miscellaneousRepository := miscellaneous.NewRepository()

	miscellaneous, ok, err := miscellaneousRepository.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting miscellaneous: %s", err.Error())
	}
	if !ok {
		return shared.NewAPIResponse(http.StatusNotFound, "miscellaneous not found")
	}

	miscellaneous.ClimateControlOccupiedSettings = body.ClimateControlOccupiedSettings
	miscellaneous.ClimateControlVacantSettings = body.ClimateControlVacantSettings

	if _, err := miscellaneousRepository.Put(ctx, miscellaneous); err != nil {
		return nil, fmt.Errorf("error putting miscellaneous: %s", err.Error())
	}

	if err := queue.SendBlankMessageToManageClimateControlsQueue(ctx); err != nil {
		return nil, fmt.Errorf("error sending message to manage climate controls queue: %s", err.Error())
	}

	return list(ctx, events.APIGatewayProxyRequest{})
}
