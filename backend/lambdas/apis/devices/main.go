package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/lambdas/apis/devices/lockcodes"
	"mlock/lambdas/helpers"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/auditlog"
	"mlock/lambdas/shared/dynamo/device"
	"mlock/lambdas/shared/dynamo/property"
	"mlock/lambdas/shared/dynamo/unit"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type DeleteResponse struct {
	Error string `json:"error"`
}

type DetailResponse struct {
	Entity shared.Device `json:"entity"`
	Extra  ExtraEntities `json:"extra"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ListResponse struct {
	Entities []shared.Device `json:"entities"`
	Extra    ExtraEntities   `json:"extra"`
}

type UpdateBody struct {
	UnitID *uuid.UUID `json:"unitId"`
}

type UpdateResponse struct {
	Entity shared.Device `json:"entity"`
	Error  string        `json:"error"`
}

type ExtraEntities struct {
	AuditLog           shared.AuditLog            `json:"auditLog"`
	Properties         []shared.Property          `json:"properties"`
	Units              []shared.Unit              `json:"units"`
	UnmanagedLockCodes []shared.RawDeviceLockCode `json:"unmanagedLockCodes"`
}

var entityRegex = regexp.MustCompile(`/devices/?`)

func main() {
	helpers.StartAPILambda(HandleRequest, []string{helpers.MiddlewareAuth})
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	match, err := regexp.MatchString(`^/devices/[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}/lock-codes/`, req.Path)
	if err != nil {
		return shared.NewAPIResponse(http.StatusBadRequest, ErrorResponse{Error: "unable to url"})
	}
	if match {
		return lockcodes.HandleRequest(ctx, req)
	}

	switch req.HTTPMethod {
	case "DELETE":
		return delete(ctx, req)
	case "GET":
		return get(ctx, req)
	case "PUT":
		return update(ctx, req)
	default:
		return shared.NewAPIResponse(http.StatusNotImplemented, "not implemented")
	}
}

func delete(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	id := strings.Replace(req.Path, "/devices/", "", 1)
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return shared.NewAPIResponse(http.StatusBadRequest, DeleteResponse{Error: "unable to parse id"})
	}

	entity, ok, err := device.NewRepository().Get(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("error getting entity: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("unable to find entity: %s", parsedID)
	}

	// TODO: Can't delete a device that's being used. For now we'll require the status to be in an allow list.

	ok = false

	awhileAgo := time.Now().Add(-2 * time.Hour)
	if entity.LastRefreshedAt.Before(awhileAgo) {
		ok = true
	}

	if !ok {
		for _, s := range []string{shared.DeviceStatusOffline} {
			if s == entity.RawDevice.Status {
				ok = true
			}
		}
	}

	if !ok {
		return shared.NewAPIResponse(http.StatusBadRequest, DeleteResponse{
			Error: fmt.Sprintf("device can't be deleted because it was recently refreshed and/or the device status"),
		})
	}

	if err := device.NewRepository().Delete(ctx, entity.ID); err != nil {
		return nil, fmt.Errorf("error deleting entity: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, DeleteResponse{})
}

func get(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	id := entityRegex.ReplaceAllString(req.Path, "")
	if id != "" {
		return detail(ctx, req, id)
	}
	return list(ctx, req)
}

func list(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	entities, err := device.NewRepository().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting entities: %s", err.Error())
	}

	properties, err := property.NewRepository().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting properties: %s", err.Error())
	}

	units, err := unit.NewRepository().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting units: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, ListResponse{
		Entities: entities,
		Extra: ExtraEntities{
			Properties: properties,
			Units:      units,
		},
	})
}

func detail(ctx context.Context, req events.APIGatewayProxyRequest, id string) (*shared.APIResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("error parsing id: %s", err.Error())
	}

	entity, ok, err := device.NewRepository().Get(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("error getting entity: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("entity not found: %s", parsedID)
	}
	entity.SortManagedLockCodes()

	auditLog, found, err := auditlog.Get(ctx, entity.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting audit logs: %s", err.Error())
	}
	if !found {
		auditLog = shared.AuditLog{Entries: []shared.AuditLogEntry{}}
	}
	if len(auditLog.Entries) > 30 {
		auditLog.Entries = auditLog.Entries[len(auditLog.Entries)-30:]
	}

	properties, err := property.NewRepository().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting properties: %s", err.Error())
	}

	units, err := unit.NewRepository().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting units: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, DetailResponse{
		Entity: entity,
		Extra: ExtraEntities{
			AuditLog:           auditLog,
			Properties:         properties,
			Units:              units,
			UnmanagedLockCodes: entity.GenerateUnmanagedLockCodes(),
		},
	})
}

func update(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	id := entityRegex.ReplaceAllString(req.Path, "")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("error parsing id: %s", err.Error())
	}

	var body UpdateBody
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	entity, ok, err := device.NewRepository().Get(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("error getting entity: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("entity not found: %s", parsedID)
	}

	if body.UnitID != nil {
		// Just verify that it exists and has the right property ID (TODO: don't allow a unit to change its property ID if a device is assigned to it).
		unit, ok, err := unit.NewRepository().Get(ctx, *body.UnitID)
		if err != nil {
			return nil, fmt.Errorf("error getting unit: %s", err.Error())
		}
		if !ok {
			return nil, fmt.Errorf("unit not found: %s", parsedID)
		}

		if unit.PropertyID != entity.PropertyID {
			return nil, fmt.Errorf("property IDs don't match")
		}
	}

	entity.UnitID = body.UnitID

	entity, err = device.NewRepository().Put(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("error updating entity: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, UpdateResponse{Entity: entity})
}
