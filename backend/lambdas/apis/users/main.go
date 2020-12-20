package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/shared"
	"mlock/shared/datastore"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ListResponse struct {
	Users []datastore.User
}

type CreateUserBody struct {
	Email string
}

func main() {
	if err := shared.LoadConfig(); err != nil {
		fmt.Printf("ERROR loading config: %s\n", err.Error())
		return
	}
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	switch req.HTTPMethod {
	case "GET":
		return list(ctx, req)
	case "POST":
		return createUser(ctx, req)
	default:
		return shared.APIResponse(http.StatusNotImplemented, "not implemented")
	}
}

func list(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	users, err := datastore.GetUsers(nil)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %s", err.Error())
	}

	return shared.APIResponse(http.StatusOK, ListResponse{Users: users})
}

func createUser(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var body CreateUserBody
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	db, err := datastore.GetDB()
	if err != nil {
		return nil, fmt.Errorf("error getting DB: %s", err.Error())
	}

	if err := datastore.InsertUser(db, body.Email); err != nil {
		return nil, fmt.Errorf("error inserting user: %s", err.Error())
	}

	users, err := datastore.GetUsers(db)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %s", err.Error())
	}

	return shared.APIResponse(http.StatusOK, ListResponse{Users: users})
}
