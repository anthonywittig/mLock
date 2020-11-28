package main

import (
	"context"
	"errors"
	"os"
	"path"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

type MyEvent struct {
}

type Response struct {
	Messages []string `json:"Messages"`
}

func HandleRequest(ctx context.Context, event MyEvent) (Response, error) {
	ex, err := os.Executable()
	if err != nil {
		return Response{}, err
	}
	dir := path.Dir(ex)
	if err := godotenv.Load(dir + "/.env"); err != nil {
		return Response{}, errors.New("Error loading .env file")
	}

	return Response{
		Messages: []string{
			"hi",
			"bye",
			os.Getenv("TEST"),
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
