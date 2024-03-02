package main

import (
	"context"
	"fmt"
	"log"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/miscellaneous"
	"time"

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
	/*
		return Response{
			Messages: []string{"no migrations run; we should create a migration tracking table and only run them once, probably. And we should run the migrations on every deploy, probably."},
		}, nil
	*/

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	ctx = shared.CreateContextData(ctx)

	log.Printf("migrating miscellaneous...\n")
	if err := miscellaneous.Migrate(ctx); err != nil {
		return Response{}, fmt.Errorf("error migrating miscellaneous: %s", err.Error())
	}
	log.Printf("migrated miscellaneous\n")

	return Response{Messages: []string{"success!"}}, nil

	// Old code as a reference to what we once did:
	/*
		log.Printf("migrating climatecontrol...\n")
		if err := climatecontrol.Migrate(ctx); err != nil {
			return Response{}, fmt.Errorf("error migrating climatecontrol: %s", err.Error())
		}
		log.Printf("migrated climatecontrol\n")

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
			return Response{}, fmt.Errorf("error migrating dynamo units: %s", err.Error())
		}

		if err := auditlog.Migrate(ctx); err != nil {
			return Response{}, fmt.Errorf("error migrating dynamo audit log: %s", err.Error())
		}
	*/
}
