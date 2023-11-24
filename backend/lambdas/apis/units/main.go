package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/lambdas/helpers"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/device"
	"mlock/lambdas/shared/dynamo/property"
	"mlock/lambdas/shared/dynamo/unit"
	"mlock/lambdas/shared/hostaway"
	mshared "mlock/shared"
	"net/http"
	"regexp"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type DeleteResponse struct {
	Error string `json:"error"`
}

type DetailResponse struct {
	Entity shared.Unit   `json:"entity"`
	Extra  ExtraEntities `json:"extra"`
}

type ListResponse struct {
	Entities []ListResponseEntity `json:"entities"`
	Extra    ExtraEntities        `json:"extra"`
}

type ListResponseDevice struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ListResponseEntity struct {
	Devices    []ListResponseDevice `json:"devices"`
	ID         uuid.UUID            `json:"id"`
	Name       string               `json:"name"`
	PropertyID uuid.UUID            `json:"propertyId"`
	UpdatedBy  string               `json:"updatedBy"`
}

type CreateBody struct {
	Name       string    `json:"name"`
	PropertyID uuid.UUID `json:"propertyId"`
}

type CreateResponse struct {
	Entity shared.Unit `json:"entity"`
}

type UpdateBody struct {
	Name              string    `json:"name"`
	PropertyID        uuid.UUID `json:"propertyId"`
	RemotePropertyURL string    `json:"remotePropertyUrl"`
}

type UpdateResponse struct {
	Entity shared.Unit `json:"entity"`
	Error  string      `json:"error"`
}

type ExtraEntities struct {
	Devices      []shared.Device      `json:"devices"`
	Properties   []shared.Property    `json:"properties"`
	Reservations []shared.Reservation `json:"reservations"`
}

var unitsRegex = regexp.MustCompile(`/units/?`)

func main() {
	helpers.StartAPILambda(HandleRequest, []string{helpers.MiddlewareAuth})
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	switch req.HTTPMethod {
	case "DELETE":
		return delete(ctx, req)
	case "GET":
		return get(ctx, req)
	case "POST":
		return create(ctx, req)
	case "PUT":
		return update(ctx, req)
	default:
		return shared.NewAPIResponse(http.StatusNotImplemented, "not implemented")
	}
}

func delete(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	id := unitsRegex.ReplaceAllString(req.Path, "")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return shared.NewAPIResponse(http.StatusBadRequest, DeleteResponse{Error: "unable to parse id"})
	}

	entity, ok, err := unit.NewRepository().Get(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("error getting entity: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("unable to find entity")
	}

	if err := unit.NewRepository().Delete(ctx, entity.ID); err != nil {
		return nil, fmt.Errorf("error deleting entity: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, DeleteResponse{})
}

func get(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	id := unitsRegex.ReplaceAllString(req.Path, "")
	if id != "" {
		return detail(ctx, req, id)
	}
	return list(ctx, req)
}

func list(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	entities, err := unit.NewRepository().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting entities: %s", err.Error())
	}

	devices, err := device.NewRepository().ListByUnit(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting devices: %s", err.Error())
	}

	properties, err := property.NewRepository().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting properties: %s", err.Error())
	}

	ds := map[uuid.UUID][]ListResponseDevice{}
	for id, devicess := range devices {
		dss := []ListResponseDevice{}
		for _, device := range devicess {
			dss = append(
				dss,
				ListResponseDevice{
					ID:   device.ID,
					Name: device.RawDevice.Name,
				},
			)
		}
		ds[id] = dss
	}

	es := []ListResponseEntity{}
	for _, e := range entities {
		esd, ok := ds[e.ID]
		if !ok {
			esd = []ListResponseDevice{}
		}

		es = append(
			es,
			ListResponseEntity{
				Devices:    esd,
				ID:         e.ID,
				Name:       e.Name,
				PropertyID: e.PropertyID,
				UpdatedBy:  e.UpdatedBy,
			},
		)
	}

	return shared.NewAPIResponse(http.StatusOK, ListResponse{
		Entities: es,
		Extra: ExtraEntities{
			Properties: properties,
		},
	})
}

func detail(ctx context.Context, req events.APIGatewayProxyRequest, id string) (*shared.APIResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("error parsing id: %s", err.Error())
	}

	entity, ok, err := unit.NewRepository().Get(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("error getting entity: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("entity not found: %s", parsedID)
	}

	tzName, err := mshared.GetConfig("TIME_ZONE")
	if err != nil {
		return nil, fmt.Errorf("error getting time zone name: %s", err.Error())
	}

	tz, err := time.LoadLocation(tzName)
	if err != nil {
		return nil, fmt.Errorf("error getting time zone %s", err.Error())
	}

	reservations := []shared.Reservation{}
	if entity.RemotePropertyURL != "" {
		hostawayReservations, err := hostaway.NewRepository(tz, "").Get(ctx, entity)
		if err != nil {
			return nil, fmt.Errorf("error getting reservation items: %s", err.Error())
		}
		reservations = append(reservations, hostawayReservations...)
	}

	properties, err := property.NewRepository().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting properties: %s", err.Error())
	}

	devices, err := device.NewRepository().ListForUnit(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("error getting devices: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, DetailResponse{
		Entity: entity,
		Extra: ExtraEntities{
			Devices:      devices,
			Properties:   properties,
			Reservations: reservations,
		},
	})
}

func update(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	id := unitsRegex.ReplaceAllString(req.Path, "")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("error parsing id: %s", err.Error())
	}

	var body UpdateBody
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	entity, ok, err := unit.NewRepository().Get(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("error getting entity: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("entity not found: %s", parsedID)
	}

	entity.Name = body.Name
	entity.PropertyID = body.PropertyID
	entity.RemotePropertyURL = body.RemotePropertyURL

	entity, err = unit.NewRepository().Put(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("error updating entity: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, UpdateResponse{Entity: entity})
}

func create(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	var body CreateBody
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	entity, err := unit.NewRepository().Put(ctx, shared.Unit{
		ID:         uuid.New(),
		Name:       body.Name,
		PropertyID: body.PropertyID,
	})
	if err != nil {
		return nil, fmt.Errorf("error inserting entity: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, CreateResponse{Entity: entity})
}
