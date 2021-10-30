package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
}

type Response struct {
	Messages []string `json:"Messages"`
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, event MyEvent) (Response, error) {
	return Response{
		Messages: []string{"no migrations run - we should create a migration tracking table and only run them once, probably. And we should run the migrations on every deploy, probably."},
	}, nil

	// Old code as a reference to what we once did:
	/*
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		log.Printf("starting migrations\n")
		if err := mshared.LoadConfig(); err != nil {
			return Response{}, fmt.Errorf("error loading config: %s", err.Error())
		}

		ctx = shared.CreateContextData(ctx)

		log.Printf("migrating devices...\n")
		if err := device.Migrate(ctx); err != nil {
			return Response{}, fmt.Errorf("error migrating dynamo devices: %s", err.Error())
		}
		log.Printf("migrated devices\n")

		log.Printf("migrating user...\n")
		if err := user.Migrate(ctx); err != nil {
			return Response{}, fmt.Errorf("error migrating dynamo users: %s", err.Error())
		}
		log.Printf("migrated user\n")

		if err := property.Migrate(ctx); err != nil {
			return Response{}, fmt.Errorf("error migrating dynamo properties: %s", err.Error())
		}

		if err := unit.Migrate(ctx); err != nil {
			return Response{}, fmt.Errorf("error migrating dynamo properties: %s", err.Error())
		}

	*/
}
