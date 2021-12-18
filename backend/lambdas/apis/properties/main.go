package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/lambdas/helpers"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/property"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type DeleteResponse struct {
	Error string `json:"error"`
}

type ListResponse struct {
	Entities []shared.Property `json:"entities"`
}

type CreateRequest struct {
	Name string `json:"name"`
}

type CreateResponse struct {
	Entity shared.Property `json:"entity"`
}

func main() {
	helpers.StartAPILambda(HandleRequest, []string{helpers.MiddlewareAuth})
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
	id := strings.Replace(req.Path, "/properties/", "", 1)
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return shared.NewAPIResponse(http.StatusBadRequest, DeleteResponse{Error: "unable to parse id"})
	}

	entity, ok, err := property.NewRepository().Get(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("error getting entity: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("unable to find entity: %s", parsedID)
	}

	// TODO: Can't delete a property with existing units.

	if err := property.NewRepository().Delete(ctx, entity.ID); err != nil {
		return nil, fmt.Errorf("error deleting entity: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, DeleteResponse{})
}

func list(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	entities, err := property.NewRepository().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting entities: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, ListResponse{Entities: entities})
}

func create(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	var body CreateRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	entity, err := property.NewRepository().Put(ctx, shared.Property{
		ID:   uuid.New(),
		Name: body.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("error inserting entity: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, CreateResponse{Entity: entity})
}
