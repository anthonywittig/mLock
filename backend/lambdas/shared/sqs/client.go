package sqs

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"
	"mlock/shared/sqs"
)

func GetClient(ctx context.Context) (*sqs.Client, error) {
	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting context data: %s", err.Error())
	}

	if cd.SQS != nil {
		return cd.SQS, nil
	}

	c, err := sqs.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %s", err.Error())
	}
	cd.SQS = c
	return cd.SQS, nil
}
