package main

import (
	"context"
	"fmt"
	"log"
	"mlock/lambdas/shared"
	mshared "mlock/shared"
	"mlock/shared/protos/messaging"
	"time"

	"github.com/aws/aws-lambda-go/events"
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

func HandleRequest(ctx context.Context, event events.SQSEvent) (Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	ctx = shared.CreateContextData(ctx)

	if err := mshared.LoadConfig(); err != nil {
		return Response{}, fmt.Errorf("error loading config: %s", err.Error())
	}

	log.Printf("handling message(s)")
	for _, m := range event.Records {
		mType, ok := m.MessageAttributes[mshared.SQSProtoMessageTypeKey]
		if !ok {
			// Send to deadletter queue?
			return Response{}, fmt.Errorf("can't get message type")
		}

		switch *mType.StringValue {
		case string((&messaging.OnPremResponse{}).ProtoReflect().Descriptor().FullName()):
			message, err := mshared.DecodeMessageOnPremResponse(m.Body)
			if err != nil {
				return Response{}, fmt.Errorf("error decoding messages: %s", err.Error())
			}
			log.Printf("message description: %s", message.Description)

			/*
				resp, err := hab.ProcessCommand(ctx, message)
				if err != nil {
					return Response{}, fmt.Errorf("error processing command: %s", err.Error())
				}
				log.Printf("response: %s", resp.Description)
			*/
		default:
			return Response{}, fmt.Errorf("unhandled message type: %s", *mType.StringValue)
		}
	}

	return Response{
		Messages: []string{"yo!"},
	}, nil
}
