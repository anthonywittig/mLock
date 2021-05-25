package sqs

import (
	"context"
	"errors"
	"fmt"
	"mlock/shared"

	"github.com/aws/aws-sdk-go-v2/config"
	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Client struct {
	sqsClient                  *awssqs.Client
	queueURL                   *string
	visibilityTimeoutInSeconds int32
}

func New(ctx context.Context) (*Client, error) {
	queue := shared.GetConfig("AWS_SQS_QUEUE_PREFIX") + "-in.fifo"

	svc, err := sqsClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting sqs client: %s", err.Error())
	}

	urlResult, err := getQueueURL(ctx, svc, queue)
	if err != nil {
		return nil, fmt.Errorf("error getting the queue URL: %s", err.Error())
	}

	queueURL := urlResult.QueueUrl

	/*
		msgResult, err := getMessages(sess, queueURL, timeout)
		if err != nil {
			fmt.Println("Got an error receiving messages:")
			fmt.Println(err)
			return
		}
		fmt.Println("Message ID:     " + *msgResult.Messages[0].MessageId)
		fmt.Println("Message Handle: " + *msgResult.Messages[0].ReceiptHandle)
	*/

	return &Client{
		sqsClient:                  svc,
		queueURL:                   queueURL,
		visibilityTimeoutInSeconds: 60,
	}, nil
}

func (c *Client) GetMessages(ctx context.Context) ([]types.Message, error) {
	msgResult, err := c.sqsClient.ReceiveMessage(ctx, &awssqs.ReceiveMessageInput{
		AttributeNames: []types.QueueAttributeName{
			"All",
		},
		MessageAttributeNames: []string{
			"All",
		},
		QueueUrl:            c.queueURL,
		MaxNumberOfMessages: 1,
		VisibilityTimeout:   c.visibilityTimeoutInSeconds,
	})
	if err != nil {
		return nil, err
	}

	return msgResult.Messages, nil
}

func (c *Client) AcknowledgeMessage(ctx context.Context, message types.Message) error {
	_, err := c.sqsClient.DeleteMessage(ctx, &awssqs.DeleteMessageInput{
		QueueUrl:      c.queueURL,
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		return fmt.Errorf("error deleting message: %s", err.Error())
	}

	return nil
}

func getQueueURL(ctx context.Context, svc *awssqs.Client, queue string) (*awssqs.GetQueueUrlOutput, error) {
	urlResult, err := svc.GetQueueUrl(ctx, &awssqs.GetQueueUrlInput{
		QueueName: &queue,
	})
	if err != nil {
		return nil, err
	}

	return urlResult, nil
}

func sqsClient(ctx context.Context) (*awssqs.Client, error) {
	region := config.WithRegion(shared.GetConfig("AWS_REGION"))
	cfg, err := config.LoadDefaultConfig(ctx, region)
	profile := shared.GetConfig("AWS_PROFILE")
	if profile != "" {
		cfg, err = config.LoadDefaultConfig(ctx, region, config.WithSharedConfigProfile(profile))
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting aws config: %s", err.Error()))
	}

	/*
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Config: cfg,
		}))
	*/
	//github.com/aws/aws-sdk-go-v2/aws
	//github.com/aws/aws-sdk-go/aws

	return awssqs.NewFromConfig(cfg), nil
}
