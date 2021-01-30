package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/lambdas/helpers"
	"mlock/shared"
	"mlock/shared/dynamo/property"
	"mlock/shared/dynamo/unit"
	"mlock/shared/ical"
	"net/http"
	"net/url"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

type DeleteResponse struct {
	Error string `json:"error"`
}

type DetailResponse struct {
	Entity shared.Unit   `json:"entity"`
	Extra  ExtraEntities `json:"extra"`
}

type ListResponse struct {
	Entities []shared.Unit `json:"entities"`
	Extra    ExtraEntities `json:"extra"`
}

type CreateBody struct {
	Name         string `json:"name"`
	PropertyName string `json:"propertyName"`
}

type CreateResponse struct {
	Entity shared.Unit `json:"entity"`
}

type UpdateResponse struct {
	Entity shared.Unit `json:"entity"`
	Error  string      `json:"error"`
}

type UpdateBody struct {
	Name         string `json:"name"`
	PropertyName string `json:"propertyName"`
	CalendarURL  string `json:"calendarUrl"`
}

type ExtraEntities struct {
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
	name, err := url.QueryUnescape(unitsRegex.ReplaceAllString(req.Path, ""))
	if err != nil {
		return nil, fmt.Errorf("error unescaping name: %s", err.Error())
	}
	if name == "" {
		return shared.NewAPIResponse(http.StatusBadRequest, DeleteResponse{Error: "unable to parse name"})
	}

	entity, ok, err := unit.Get(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("error getting entity: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("unable to find entity")
	}

	if err := unit.Delete(ctx, entity.Name); err != nil {
		return nil, fmt.Errorf("error deleting entity: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, DeleteResponse{})
}

func get(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	name, err := url.QueryUnescape(unitsRegex.ReplaceAllString(req.Path, ""))
	if err != nil {
		return nil, fmt.Errorf("error unescaping name: %s", err.Error())
	}
	if name != "" {
		return detail(ctx, req, name)
	}
	return list(ctx, req)
}

func list(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	entities, err := unit.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting entities: %s", err.Error())
	}

	properties, err := property.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting properties: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, ListResponse{
		Entities: entities,
		Extra: ExtraEntities{
			Properties: properties,
		},
	})
}

func detail(ctx context.Context, req events.APIGatewayProxyRequest, name string) (*shared.APIResponse, error) {
	entity, ok, err := unit.Get(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("error getting entity: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("entity not found: %s", name)
	}

	var reservations []shared.Reservation
	if entity.CalendarURL != "" {
		var err error
		// TODO: cache this.
		reservations, err = ical.Get(context.Background(), entity.CalendarURL)
		if err != nil {
			return nil, fmt.Errorf("error getting calendar items: %s", err.Error())
		}
	}

	properties, err := property.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting properties: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, DetailResponse{
		Entity: entity,
		Extra: ExtraEntities{
			Properties:   properties,
			Reservations: reservations,
		},
	})
}

func update(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	name, err := url.QueryUnescape(unitsRegex.ReplaceAllString(req.Path, ""))
	if err != nil {
		return nil, fmt.Errorf("error unescaping name: %s", err.Error())
	}
	if name == "" {
		return shared.NewAPIResponse(http.StatusBadRequest, UpdateResponse{Error: "unable to parse name"})
	}

	var body UpdateBody
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	entity, ok, err := unit.Get(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("error getting entity: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("entity not found: _%s_", name)
	}

	entity.Name = body.Name
	entity.PropertyName = body.PropertyName
	entity.CalendarURL = body.CalendarURL

	entity, err = unit.Put(ctx, name, entity)
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

	entity, err := unit.Put(ctx, "", shared.Unit{
		Name:         body.Name,
		PropertyName: body.PropertyName,
	})
	if err != nil {
		return nil, fmt.Errorf("error inserting entity: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, CreateResponse{Entity: entity})
}
