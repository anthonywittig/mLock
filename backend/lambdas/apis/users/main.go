package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/shared"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type DeleteResponse struct {
	Message string
}

type ListResponse struct {
	Users []shared.User
}

type CreateUserBody struct {
	Email string
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
		return createUser(ctx, req)
	default:
		return shared.NewAPIResponse(http.StatusNotImplemented, "not implemented")
	}
}

func delete(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {

	// Can't delete yourself.

	return shared.NewAPIResponse(http.StatusOK, DeleteResponse{Message: req.Path})

	/*
		users, err := shared.GetUsers(ctx)
		if err != nil {
			return nil, fmt.Errorf("error getting users: %s", err.Error())
		}

		return shared.NewAPIResponse(http.StatusOK, ListResponse{Users: users})
	*/
}

func list(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	users, err := shared.GetUsers(ctx)
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

	if err := shared.InsertUser(ctx, body.Email); err != nil {
		return nil, fmt.Errorf("error inserting user: %s", err.Error())
	}

	users, err := shared.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, ListResponse{Users: users})
}
