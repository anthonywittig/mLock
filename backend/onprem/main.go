package main

import (
	"context"
	"fmt"
	"log"
	"mlock/onprem/app"
	"mlock/onprem/sqs"
	"mlock/shared"
	"time"
)

func main() {
	for {
		if err := run(); err != nil {
			log.Fatalf("error running application: %s", err.Error())
		}
		log.Println("(sleeping for a minute before trying again)")
		time.Sleep(1 * time.Minute)
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
