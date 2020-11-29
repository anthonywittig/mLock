package main

import (
	"context"
	"fmt"
	"mlock/shared"
	"mlock/shared/datastore"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
}

type Response struct {
	Users []datastore.User
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, event MyEvent) (Response, error) {
	if err := shared.LoadConfig(); err != nil {
		return Response{}, fmt.Errorf("error loading config: %s", err.Error())
	}

	users, err := datastore.GetUsers()
	if err != nil {
		return Response{}, fmt.Errorf("error getting users: %s", err.Error())
	}

	return Response{
		Users: users,
	}, nil
}
