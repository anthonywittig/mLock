package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/shared"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type DeleteResponse struct {
	Error    string
	Entities *[]shared.Property
}

type ListResponse struct {
	Entities []shared.Property
}

type CreateRequest struct {
	Name string
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
	id := strings.Replace(req.Path, "/properties/", "", 1)
	if id == "" {
		return shared.NewAPIResponse(http.StatusBadRequest, DeleteResponse{Error: "unable to parse id"})
	}

	entity, ok, err := shared.GetPropertyByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting entity: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("unable to find entity")
	}

	// TODO: Can't delete a property with existing units.

	if err := shared.DeleteProperty(ctx, entity.ID); err != nil {
		return nil, fmt.Errorf("error deleting entity: %s", err.Error())
	}

	entities, err := shared.GetProperties(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting entities: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, DeleteResponse{Entities: &entities})
}

func list(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	/*
		entities, err := shared.GetProperties(ctx)
		if err != nil {
			return nil, fmt.Errorf("error getting entities: %s", err.Error())
		}
	*/

	return shared.NewAPIResponse(http.StatusOK, ListResponse{Entities: []shared.Property{}})
}

func create(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	var body CreateRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	if err := shared.InsertProperty(ctx, body.Name); err != nil {
		return nil, fmt.Errorf("error inserting entity: %s", err.Error())
	}

	entities, err := shared.GetProperties(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting entities: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, ListResponse{Entities: entities})
}
