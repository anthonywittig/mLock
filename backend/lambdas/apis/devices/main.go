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
	"mlock/lambdas/shared/dynamo/unit"
	"mlock/lambdas/shared/ezlo"
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
		return shared.NewAPIResponse(http.StatusBadRequest, ErrorResponse{Error: "unable to parse request"})
	}
	if match {
		return lockcodes.HandleRequest(ctx, req)
	}

	match, err = regexp.MatchString(`^/devices/[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}/reboot-controller/`, req.Path)
	if err != nil {
		return shared.NewAPIResponse(http.StatusBadRequest, ErrorResponse{Error: "unable to parse request"})
	}
	if match && req.HTTPMethod == "POST" {
		return rebootController(ctx, req)
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
			Error: "device can't be deleted because it was recently refreshed and/or the device status",
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

	units, err := unit.NewRepository().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting units: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, ListResponse{
		Entities: entities,
		Extra: ExtraEntities{
			Units: units,
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

	auditLog, found, err := auditlog.Get(ctx, entity.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting audit logs: %s", err.Error())
	}
	if !found {
		auditLog = shared.AuditLog{Entries: []shared.AuditLogEntry{}}
	}
	if len(auditLog.Entries) > 100 {
		auditLog.Entries = auditLog.Entries[len(auditLog.Entries)-100:]
	}

	// Reverse the entries so that the newer items are first.
	// https://github.com/golang/go/wiki/SliceTricks#reversing
	for i := len(auditLog.Entries)/2 - 1; i >= 0; i-- {
		opp := len(auditLog.Entries) - 1 - i
		auditLog.Entries[i], auditLog.Entries[opp] = auditLog.Entries[opp], auditLog.Entries[i]
	}

	units, err := unit.NewRepository().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting units: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, DetailResponse{
		Entity: entity,
		Extra: ExtraEntities{
			AuditLog:           auditLog,
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

	entity.UnitID = body.UnitID

	entity, err = device.NewRepository().Put(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("error updating entity: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, UpdateResponse{Entity: entity})
}

func rebootController(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	id := entityRegex.ReplaceAllString(req.Path, "")
	id = strings.Replace(id, "/reboot-controller/", "", 1)
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

	connectionPool := ezlo.NewConnectionPool()
	defer connectionPool.Close()

	deviceController := ezlo.NewDeviceController(connectionPool)
	deviceController.RebootController(ctx, entity)

	return shared.NewAPIResponse(http.StatusOK, ErrorResponse{Error: ""})
}
