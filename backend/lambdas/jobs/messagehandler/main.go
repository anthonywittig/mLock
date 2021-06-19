package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/device"
	"mlock/lambdas/shared/dynamo/property"
	mshared "mlock/shared"
	"mlock/shared/protos/messaging"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
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

	log.Printf("pretending that this event is for \"zzz Anthony's House\"")
	props, err := property.List(ctx)
	if err != nil {
		return Response{}, fmt.Errorf("error loading properties: %s", err.Error())
	}
	anthonysHouse := shared.Property{}
	for _, p := range props {
		if p.Name == "zzz Anthony's House" {
			anthonysHouse = p
		}
	}
	if anthonysHouse.ID == uuid.Nil {
		return Response{}, fmt.Errorf("error getting Anthony's house: %s", err.Error())
	}

	for _, m := range event.Records {
		mType, ok := m.MessageAttributes[mshared.SQSProtoMessageTypeKey]
		if !ok {
			// Send to deadletter queue?
			return Response{}, fmt.Errorf("can't get message type")
		}

		switch *mType.StringValue {
		case string((&messaging.OnPremHabCommandResponse{}).ProtoReflect().Descriptor().FullName()):
			log.Printf("message is an onprem hab command response\n")

			message, err := mshared.DecodeMessageOnPremHabCommandResponse(m.Body)
			if err != nil {
				return Response{}, fmt.Errorf("error decoding messages: %s", err.Error())
			}
			log.Printf("message description: %s\n", message.Description)
			log.Printf("message body: %s\n", string(message.Response))

			if err := handleOnPremHabResponse(ctx, anthonysHouse, message); err != nil {
				return Response{}, fmt.Errorf("error handling on-prem HAB response: %s", err.Error())
			}
		case string((&messaging.OnPremError{}).ProtoReflect().Descriptor().FullName()):
			log.Printf("message is an onprem error response\n")

			message := &messaging.OnPremError{}
			err := mshared.DecodeMessage(m.Body, message)
			if err != nil {
				return Response{}, fmt.Errorf("error decoding messages: %s", err.Error())
			}
			log.Printf("message's error message: %s\n", message.ErrorMessage)

			// TODO: we want to signal error, but we can't fail the request or it'll just make us reprocess the message later.
			return Response{}, nil
		default:
			return Response{}, fmt.Errorf("unhandled message type: %s", *mType.StringValue)
		}
	}

	return Response{
		Messages: []string{"yo!"},
	}, nil
}

func handleOnPremHabResponse(ctx context.Context, property shared.Property, in *messaging.OnPremHabCommandResponse) error {
	switch in.HabCommand.CommandType {
	case messaging.HabCommand_UNKNOWN:
		// Pretend things are a-ok so that we don't try to process this message again.
		logError(fmt.Errorf("command type is unknown - what have you done?"))
	case messaging.HabCommand_LIST_THINGS:
		log.Printf("got a list things command")
		if err := handleListThingsResponse(ctx, property, in); err != nil {
			// Just log and continue.
			logError(fmt.Errorf("error processing list things: %s", err.Error()))
		}

	}
	return nil
}

// logError is used when we want to log and error but we don't want to choke on the message (we want to mark it as processed and move on with our lives).
func logError(err error) {
	// TODO: do something that actually notifies us of an error.
	log.Printf("ERROR: %s", err.Error())
}

func handleListThingsResponse(ctx context.Context, property shared.Property, in *messaging.OnPremHabCommandResponse) error {
	ts := []shared.HABThing{}
	if err := json.Unmarshal(in.Response, &ts); err != nil {
		return fmt.Errorf("error parsing json: %s", err.Error())
	}

	eds, err := device.List(ctx)
	if err != nil {
		return fmt.Errorf("error parsing json: %s", err.Error())
	}

	for _, t := range ts {
		d := shared.Device{
			ID:              uuid.New(),
			PropertyID:      property.ID,
			HABThing:        t,
			LastRefreshedAt: time.Now(),
		}
		for _, ed := range eds {
			if ed.PropertyID == property.ID && ed.HABThing.UID == t.UID {
				// We found a match.
				d.ID = ed.ID
			}
		}

		if _, err := device.Put(ctx, d); err != nil {
			return fmt.Errorf("error putting device: %s", err.Error())
		}
	}

	return nil
}
