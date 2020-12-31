package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/shared"
	"mlock/shared/datastore"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ListResponse struct {
	Users []datastore.User
}

type CreateUserBody struct {
	Email string
}

func main() {
	shared.StartAPILambda(HandleRequest, []string{shared.MiddlewareAuth})
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return list(ctx, req)
	case "POST":
		return createUser(ctx, req)
	default:
		return shared.NewAPIResponse(http.StatusNotImplemented, "not implemented")
	}
}

func list(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	users, err := datastore.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, ListResponse{Users: users})
}

func createUser(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	var body CreateUserBody
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	if err := datastore.InsertUser(ctx, body.Email); err != nil {
		return nil, fmt.Errorf("error inserting user: %s", err.Error())
	}

	users, err := datastore.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, ListResponse{Users: users})
}
