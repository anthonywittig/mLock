package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/shared"
	"mlock/shared/models/units"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type DeleteResponse struct {
	Error string `json:"error"`
}

type ListResponse struct {
	Entities []units.Entity `json:"entities"`
	Extra    ExtraEntities  `json:"extra"`
}

type CreateRequest struct {
	Name       string    `json:"name"`
	PropertyID uuid.UUID `json:"PropertyId"`
}

type CreateResponse struct {
	Entity units.Entity `json:"entity"`
}

type ExtraEntities struct {
	Properties []shared.Property `json:"properties"`
}

func main() {
	shared.StartAPILambda(HandleRequest, []string{shared.MiddlewareAuth})
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	switch req.HTTPMethod {
	case "DELETE":
		return delete(ctx, req)
	case "GET":
		return list(ctx, req)
	case "POST":
		return create(ctx, req)
	default:
		return shared.NewAPIResponse(http.StatusNotImplemented, "not implemented")
	}
}

func delete(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	id := strings.Replace(req.Path, "/units/", "", 1)
	if id == "" {
		return shared.NewAPIResponse(http.StatusBadRequest, DeleteResponse{Error: "unable to parse id"})
	}

	entity, ok, err := units.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting entity: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("unable to find entity")
	}

	if err := units.Delete(ctx, entity.ID); err != nil {
		return nil, fmt.Errorf("error deleting entity: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, DeleteResponse{})
}

func list(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	entities, err := units.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting entities: %s", err.Error())
	}

	properties, err := shared.GetProperties(ctx)
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

func create(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	var body CreateRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	entity, err := units.Insert(ctx, body.Name, body.PropertyID)
	if err != nil {
		return nil, fmt.Errorf("error inserting entity: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, CreateResponse{Entity: entity})
}
