package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
}

type Response struct {
	Users []string `json:"Users"`
}

func HandleRequest(ctx context.Context, event MyEvent) (Response, error) {
	return Response{
		Users: []string{
			"hello world",
			"joe.smith@gmail.com",
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
