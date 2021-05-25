package sqs

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"
	mshared "mlock/shared"
	"mlock/shared/protos/messaging"

	"github.com/aws/aws-sdk-go-v2/config"
	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
)

func GetClient(ctx context.Context) (*awssqs.Client, error) {
	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting context data: %s", err.Error())
	}

	if cd.SQS != nil {
		return cd.SQS, nil
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-1"))
	if err != nil {
		return nil, fmt.Errorf("error loading config: %s", err.Error())
	}

	cd.SQS = awssqs.NewFromConfig(cfg)
	return cd.SQS, nil
}

func SendMessage(ctx context.Context, queuePrefix string, message *messaging.HabCommand) error {
	s, err := GetClient(ctx)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err.Error())
	}

	queueName := queuePrefix + "-in.fifo"
	getQueueUrlResult, err := s.GetQueueUrl(ctx, &awssqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		return fmt.Errorf("error getting queue URL: %s", err.Error())
	}

	encodedMessage, err := mshared.EncodeMessage(message)
	if err != nil {
		return fmt.Errorf("error encoding message: %s", err.Error())
	}

	s.SendMessage(ctx, &awssqs.SendMessageInput{
		MessageBody: &encodedMessage,
		QueueUrl:    getQueueUrlResult.QueueUrl,
		MessageAttributes: map[string]types.MessageAttributeValue{
			mshared.SQSProtoMessageTypeKey: {
				DataType:    aws.String("String"),
				StringValue: aws.String("messaging.HabCommand"),
			},
		},
		MessageDeduplicationId:  aws.String(uuid.New().String()), // TODO: Probably want this based on a DB ID.
		MessageGroupId:          aws.String("default"),           // This _could_ be based on the lock, with a default for everything else.
		MessageSystemAttributes: map[string]types.MessageSystemAttributeValue{},
	})

	return nil
}
