package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/lambdas/helpers"
	"mlock/shared"
	"mlock/shared/postgres/property"
	"mlock/shared/postgres/unit"
	"net/http"
	"regexp"

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
	Entities []shared.Unit `json:"entities"`
	Extra    ExtraEntities `json:"extra"`
}

type CreateBody struct {
	Name       string    `json:"name"`
	PropertyID uuid.UUID `json:"PropertyId"`
}
type UpdateBody = CreateBody

type CreateResponse struct {
	Entity shared.Unit `json:"entity"`
}

type UpdateResponse struct {
	Entity shared.Unit `json:"entity"`
	Error  string      `json:"error"`
}

type ExtraEntities struct {
	Properties []shared.Property `json:"properties"`
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
	if id == "" {
		return shared.NewAPIResponse(http.StatusBadRequest, DeleteResponse{Error: "unable to parse id"})
	}

	entity, ok, err := unit.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting entity: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("unable to find entity")
	}

	if err := unit.Delete(ctx, entity.ID); err != nil {
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
	entities, err := unit.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting entities: %s", err.Error())
	}

	properties, err := property.GetAll(ctx)
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

func detail(ctx context.Context, req events.APIGatewayProxyRequest, id string) (*shared.APIResponse, error) {
	entity, ok, err := unit.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting entity: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("entity not found: %s", id)
	}

	properties, err := property.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting properties: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, DetailResponse{
		Entity: entity,
		Extra: ExtraEntities{
			Properties: properties,
		},
	})
}

func update(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	id := unitsRegex.ReplaceAllString(req.Path, "")
	if id == "" {
		return shared.NewAPIResponse(http.StatusBadRequest, UpdateResponse{Error: "unable to parse id"})
	}

	var body UpdateBody
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	entity, ok, err := unit.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting entity: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("entity not found: %s", id)
	}

	entity.Name = body.Name
	entity.PropertyID = body.PropertyID

	entity, err = unit.Update(ctx, entity)
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

	entity, err := unit.Insert(ctx, body.Name, body.PropertyID)
	if err != nil {
		return nil, fmt.Errorf("error inserting entity: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, CreateResponse{Entity: entity})
}
