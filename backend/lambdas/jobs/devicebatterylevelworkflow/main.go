package main

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/sqs"
	"mlock/lambdas/shared/workflows/devicebatterylevel"
	mshared "mlock/shared"
	"strings"
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
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	ctx = shared.CreateContextData(ctx)

	if err := mshared.LoadConfig(); err != nil {
		return Response{}, fmt.Errorf("error loading config: %s", err.Error())
	}

	queuePrefixes, err := mshared.GetConfig("AWS_SQS_QUEUE_PREFIXES")
	if err != nil {
		return Response{}, fmt.Errorf("error getting queue prefixes: %s", err.Error())
	}

	// TODO: move this to properties.
	queueNames := []string{}
	for _, n := range strings.Split(queuePrefixes, ",") {
		queueNames = append(queueNames, n+"-in.fifo")
	}

	s, err := sqs.GetClient(ctx)
	if err != nil {
		return Response{}, fmt.Errorf("error getting sqs client: %s", err.Error())

	}

	if err := devicebatterylevel.KickOffUpdates(ctx, s, queueNames); err != nil {
		return Response{}, fmt.Errorf("error kicking off updates: %s", err.Error())
	}

	return Response{
		Messages: []string{},
	}, nil
}
