package app

import (
	"context"
	"fmt"
	"log"
	"mlock/onprem/hab"
	"mlock/onprem/sqs"
	"mlock/shared"
	"mlock/shared/protos/messaging"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type App struct {
	sqsClient *sqs.Client
}

func NewApp(sqsClient *sqs.Client) (*App, error) {
	return &App{
		sqsClient: sqsClient,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	shortDelay := 1 * time.Second
	longDelay := 1 * time.Minute

	ticker := time.NewTicker(shortDelay)

	for {
		select {
		case <-ctx.Done():
			log.Printf("context says it's time to quit, ending run loop")
			return nil
		case <-ticker.C:
			handledMessage, err := a.processMessages(ctx)
			if err != nil {
				return fmt.Errorf("error on process message: %s", err.Error())
			}

			nextDelay := longDelay
			if handledMessage {
				nextDelay = shortDelay
			}

			log.Printf("sleeping for %.0f seconds\n", nextDelay.Seconds())
			ticker.Reset(nextDelay)
		}
	}
}

func (a *App) processMessages(ctx context.Context) (bool, error) {
	messages, err := a.sqsClient.GetMessages(ctx)
	if err != nil {
		return false, fmt.Errorf("error getting messages: %s", err.Error())
	}

	for _, m := range messages {
		shouldAck, err := a.processMessage(ctx, m)
		if err != nil {
			return false, fmt.Errorf("error processing message: %s", err.Error())
		}
		if shouldAck {
			if err := a.sqsClient.AcknowledgeMessage(ctx, m); err != nil {
				return false, fmt.Errorf("error processing message: %s", err.Error())
			}
		} else {
			// Since someone has a reason for us to not ack this message, we should return saying we didn't handle the message.
			return false, nil
		}

	}

	return len(messages) > 0, nil
}

func (a *App) processMessage(ctx context.Context, m types.Message) (bool, error) {
	mType, ok := m.MessageAttributes[shared.SQSProtoMessageTypeKey]
	if !ok {
		if err := a.sendErrorMessage(ctx, "no message type"); err != nil {
			return false, fmt.Errorf("failed to send error message: %s", err.Error())
		}
		return true, nil
	}

	switch *mType.StringValue {
	case string((&messaging.HabCommand{}).ProtoReflect().Descriptor().FullName()):
		message, err := shared.DecodeMessageHabCommand(*m.Body)
		if err != nil {
			if err := a.sendErrorMessage(ctx, fmt.Sprintf("error decoding messages: %s", err.Error())); err != nil {
				return false, fmt.Errorf("failed to send error message: %s", err.Error())
			}
			return true, nil
		}

		resp, err := hab.ProcessCommand(ctx, message)
		if err != nil {
			if err := a.sendErrorMessage(ctx, fmt.Sprintf("error processing command: %s", err.Error())); err != nil {
				return false, fmt.Errorf("failed to send error message: %s", err.Error())
			}
			return true, nil
		}

		if err := a.sqsClient.SendMessage(ctx, resp); err != nil {
			if err := a.sendErrorMessage(ctx, fmt.Sprintf("error sending response: %s", err.Error())); err != nil {
				return false, fmt.Errorf("failed to send error message: %s", err.Error())
			}
			return true, nil
		}
	default:
		if err := a.sendErrorMessage(ctx, fmt.Sprintf("unhandled message type: %s", *mType.StringValue)); err != nil {
			return false, fmt.Errorf("failed to send error message: %s", err.Error())
		}
		return true, nil
	}

	return true, nil
}

func (a *App) sendErrorMessage(ctx context.Context, msg string) error {
	log.Printf("sending error message: %s\n", msg)
	if err := a.sqsClient.SendMessage(ctx, &messaging.OnPremError{
		ErrorMessage: msg,
	}); err != nil {
		return fmt.Errorf("error sending response: %s", err.Error())
	}
	return nil
}
