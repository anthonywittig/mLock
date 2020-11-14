package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"name"`
}

type Response struct {
	Message string `json:"Answer:"`
}

func HandleRequest(ctx context.Context, name MyEvent) (Response, error) {
	return Response{Message: "hello world"}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
