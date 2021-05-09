package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/lambdas/helpers"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/user"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type DeleteResponse struct {
	Error string
	Users *[]shared.User
}

type ListResponse struct {
	Users []shared.User
}

type CreateUserBody struct {
	Email string
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
		return createUser(ctx, req)
	default:
		return shared.NewAPIResponse(http.StatusNotImplemented, "not implemented")
	}
}

func delete(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	id := strings.Replace(req.Path, "/users/", "", 1)
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return shared.NewAPIResponse(http.StatusBadRequest, DeleteResponse{Error: "unable to parse id"})
	}

	u, ok, err := user.NewUserService().Get(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("unable to find user")
	}

	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting context data: %s", err.Error())
	}

	// Can't delete yourself.
	if u.Email == cd.User.Email {
		return shared.NewAPIResponse(http.StatusBadRequest, DeleteResponse{Error: "can't delete oneself"})
	}

	if err := user.NewUserService().Delete(ctx, u.ID); err != nil {
		return nil, fmt.Errorf("error deleting user: %s", err.Error())
	}

	users, err := user.NewUserService().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, DeleteResponse{Users: &users})
}

func list(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	users, err := user.NewUserService().List(ctx)
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

	if _, err := user.NewUserService().Put(ctx, shared.User{
		ID:    uuid.New(),
		Email: body.Email,
	}); err != nil {
		return nil, fmt.Errorf("error inserting user: %s", err.Error())
	}

	users, err := user.NewUserService().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, ListResponse{Users: users})
}
