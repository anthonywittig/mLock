package sqs

import (
	"context"
	"errors"
	"fmt"
	"mlock/shared"

	"github.com/aws/aws-sdk-go-v2/config"
	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Client struct {
	sqsClient                  *awssqs.Client
	visibilityTimeoutInSeconds int32
}

func New(ctx context.Context) (*Client, error) {
	svc, err := sqsClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting sqs client: %s", err.Error())
	}

	return &Client{
		sqsClient:                  svc,
		visibilityTimeoutInSeconds: 60,
	}, nil
}

func (c *Client) GetMessages(ctx context.Context, queueName string) ([]types.Message, error) {
	qURL, err := c.getQueueURL(ctx, queueName)
	if err != nil {
		return nil, fmt.Errorf("error getting the queue URL: %s", err.Error())
	}

	msgResult, err := c.sqsClient.ReceiveMessage(ctx, &awssqs.ReceiveMessageInput{
		AttributeNames: []types.QueueAttributeName{
			"All",
		},
		MessageAttributeNames: []string{
			"All",
		},
		QueueUrl:            qURL,
		MaxNumberOfMessages: 1,
		VisibilityTimeout:   c.visibilityTimeoutInSeconds,
	})
	if err != nil {
		return nil, err
	}

	return msgResult.Messages, nil
}

func (c *Client) AcknowledgeMessage(ctx context.Context, queueName string, message types.Message) error {
	qURL, err := c.getQueueURL(ctx, queueName)
	if err != nil {
		return fmt.Errorf("error getting the queue URL: %s", err.Error())
	}

	_, err = c.sqsClient.DeleteMessage(ctx, &awssqs.DeleteMessageInput{
		QueueUrl:      qURL,
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		return fmt.Errorf("error deleting message: %s", err.Error())
	}

	return nil
}

func (c *Client) SendMessage(ctx context.Context, queueName string, message protoreflect.ProtoMessage) error {
	getQueueUrlResult, err := c.sqsClient.GetQueueUrl(ctx, &awssqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		return fmt.Errorf("error getting queue URL: %s", err.Error())
	}

	encodedMessage, err := shared.EncodeMessage(message)
	if err != nil {
		return fmt.Errorf("error encoding message: %s", err.Error())
	}

	c.sqsClient.SendMessage(ctx, &awssqs.SendMessageInput{
		MessageBody: &encodedMessage,
		QueueUrl:    getQueueUrlResult.QueueUrl,
		MessageAttributes: map[string]types.MessageAttributeValue{
			shared.SQSProtoMessageTypeKey: {
				DataType:    aws.String("String"),
				StringValue: aws.String(string(message.ProtoReflect().Descriptor().FullName())),
			},
		},
		MessageDeduplicationId:  aws.String(uuid.New().String()), // TODO: Probably want this based on a DB ID.
		MessageGroupId:          aws.String("default"),           // This _could_ be based on the lock, with a default for everything else.
		MessageSystemAttributes: map[string]types.MessageSystemAttributeValue{},
	})

	return nil
}

func (c *Client) getQueueURL(ctx context.Context, queue string) (*string, error) {
	urlResult, err := c.sqsClient.GetQueueUrl(ctx, &awssqs.GetQueueUrlInput{
		QueueName: &queue,
	})
	if err != nil {
		return nil, err
	}

	return urlResult.QueueUrl, nil
}

func sqsClient(ctx context.Context) (*awssqs.Client, error) {
	region := config.WithRegion(shared.GetConfigUnsafe("AWS_REGION"))
	cfg, err := config.LoadDefaultConfig(ctx, region)
	profile := shared.GetConfigUnsafe("AWS_PROFILE")
	if profile != "" {
		cfg, err = config.LoadDefaultConfig(ctx, region, config.WithSharedConfigProfile(profile))
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting aws config: %s", err.Error()))
	}

	return awssqs.NewFromConfig(cfg), nil
}
