package main

import (
	"context"
	"fmt"
	"log"
	"mlock/onprem/app"
	"mlock/onprem/sqs"
	"mlock/shared"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("error running application: %s", err.Error())
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	if err := shared.LoadConfig(); err != nil {
		return fmt.Errorf("error loading config: %s", err.Error())
	}

	sqsClient, err := sqs.New(ctx)
	if err != nil {
		return fmt.Errorf("error getting new sqs: %s", err.Error())
	}

	a, err := app.NewApp(sqsClient)
	if err != nil {
		return fmt.Errorf("error getting new app: %s", err.Error())
	}

	return a.Run(ctx)
}
