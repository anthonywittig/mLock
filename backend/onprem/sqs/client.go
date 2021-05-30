package sqs

import (
	"context"
	"fmt"
	"mlock/shared"
	msqs "mlock/shared/sqs"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Client struct {
	sqsClient      *msqs.Client
	readQueueName  string
	writeQueueName string
}

func New(ctx context.Context) (*Client, error) {
	s, err := msqs.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting sqs client: %s", err.Error())
	}

	queue, err := shared.GetConfig("AWS_SQS_QUEUE_PREFIX")
	if err != nil {
		return nil, fmt.Errorf("error getting queue config: %s", err.Error())
	}

	return &Client{
		sqsClient:      s,
		readQueueName:  queue + "-in.fifo",
		writeQueueName: queue + "-out.fifo",
	}, nil
}

func (c *Client) GetMessages(ctx context.Context) ([]types.Message, error) {
	return c.sqsClient.GetMessages(ctx, c.readQueueName)
}

func (c *Client) AcknowledgeMessage(ctx context.Context, message types.Message) error {
	return c.sqsClient.AcknowledgeMessage(ctx, c.readQueueName, message)
}

func (c *Client) SendMessage(ctx context.Context, message protoreflect.ProtoMessage) error {
	return c.sqsClient.SendMessage(ctx, c.writeQueueName, message)
}
