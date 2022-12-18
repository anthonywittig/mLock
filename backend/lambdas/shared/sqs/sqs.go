package sqs

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"
	mshared "mlock/shared"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSService struct {
	c *sqs.Client
}

func NewSQSService(ctx context.Context) (*SQSService, error) {
	c, err := getClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %s", err.Error())
	}

	return &SQSService{
		c: c,
	}, nil
}

func getClient(ctx context.Context) (*sqs.Client, error) {
	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting context data: %s", err.Error())
	}

	if cd.SQS != nil {
		return cd.SQS, nil
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-1"))
	if err != nil {
		return nil, fmt.Errorf("error getting aws config: %s", err.Error())
	}

	cd.SQS = sqs.NewFromConfig(cfg)

	return cd.SQS, nil
}

func (s *SQSService) SendBlankMessageToPollSchedulesQueue(ctx context.Context) error {
	queueURL, err := mshared.GetConfig("POLL_SCHEDULES_QUEUE_URL")
	if err != nil {
		return fmt.Errorf("empty queue url in config")
	}

	sMInput := &sqs.SendMessageInput{
		MessageBody: aws.String("sent from SQSService"),
		QueueUrl:    aws.String(queueURL),
	}
	_, err = s.c.SendMessage(ctx, sMInput)
	if err != nil {
		return fmt.Errorf("error sending sqs message: %s", err.Error())
	}

	return nil
}
