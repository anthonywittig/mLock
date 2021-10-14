package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/device"
	"mlock/lambdas/shared/dynamo/property"
	"mlock/lambdas/shared/ses"
	"mlock/lambdas/shared/workflows/devicebatterylevel"
	mshared "mlock/shared"
	"mlock/shared/protos/messaging"
	"strings"
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

// TODO: move this to the property object.
var queueToPropertyName = map[string]string{
	"test-out.fifo": "zzz anthony's house",
	"rpi1-out.fifo": "zion's camp",
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
		prop, err := getPropertyForEventARN(ctx, m.EventSourceARN)
		if err != nil {
			return Response{}, fmt.Errorf("error getting property for \"%s\": %s", m.EventSourceARN, err.Error())
		}

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

			if err := handleOnPremHabResponse(ctx, prop, message); err != nil {
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
	case messaging.HabCommand_LIST_ITEMS:
		log.Printf("got a list items command")
		if err := devicebatterylevel.ProcessListItems(ctx, property, in); err != nil {
			// Just log and continue.
			logError(fmt.Errorf("error processing list items: %s", err.Error()))
		}
	default:
		// Just log and continue.
		logError(fmt.Errorf("unhandled command type: %s", in.HabCommand.CommandType.String()))
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
		return fmt.Errorf("error getting devices: %s", err.Error())
	}

	transitioningToOfflineDevices := []shared.Device{}
	offlineDevices := []shared.Device{}

	for _, t := range ts {
		d := shared.Device{
			History: []shared.DeviceHistory{
				{
					Description: "Initial State",
					HABThing:    t,
					RecordedAt:  time.Now(),
				},
			},
			ID: uuid.New(),
		}

		for _, ed := range eds {
			if ed.PropertyID == property.ID && ed.HABThing.UID == t.UID {
				// We found a match.
				d = ed

				wasOffline := t.StatusInfo.Status == shared.DeviceStatusOffline
				isOffline := d.HABThing.StatusInfo.Status == shared.DeviceStatusOffline
				if isOffline {
					offlineDevices = append(offlineDevices, d)
					if !wasOffline {
						now := time.Now()
						d.LastWentOfflineAt = &now
						transitioningToOfflineDevices = append(transitioningToOfflineDevices, d)
					}
				}

				statusChanged := (t.StatusInfo.Status != d.HABThing.StatusInfo.Status) || (t.StatusInfo.StatusDetail != d.HABThing.StatusInfo.StatusDetail)
				if statusChanged {
					d.History = append(d.History, shared.DeviceHistory{
						Description: "Status Changed",
						HABThing:    t,
						RecordedAt:  time.Now(),
					})
				}

				maxHistoryCount := 1
				historyStartIndex := len(d.History) - maxHistoryCount
				if historyStartIndex > 0 {
					d.History = d.History[historyStartIndex:]
				}
			}
		}

		d.PropertyID = property.ID
		d.HABThing = t
		d.LastRefreshedAt = time.Now()

		if _, err := device.Put(ctx, d); err != nil {
			return fmt.Errorf("error putting device: %s", err.Error())
		}
	}

	if err := sendOfflineDeviceEmail(ctx, transitioningToOfflineDevices, offlineDevices); err != nil {
		return fmt.Errorf("error sending offline device email: %s", err.Error())
	}

	return nil
}

func sendOfflineDeviceEmail(ctx context.Context, transitioningToOfflineDevices []shared.Device, offlineDevices []shared.Device) error {
	if len(transitioningToOfflineDevices) == 0 {
		return nil
	}

	var sb strings.Builder

	sb.WriteString("<h1>Devices That Recently Went Offline</h1>")
	sb.WriteString("<ul>")
	for _, d := range transitioningToOfflineDevices {
		sb.WriteString(fmt.Sprintf("<li>Device: %s</li>", d.HABThing.Label))
	}
	sb.WriteString("</ul>")

	sb.WriteString("<h1>Devices That Are Currently Offline</h1>")
	sb.WriteString("<ul>")
	for _, d := range offlineDevices {
		sb.WriteString(fmt.Sprintf("<li>Device: %s</li>", d.HABThing.Label))
	}
	sb.WriteString("</ul>")

	if err := ses.SendEamil(ctx, "MursetLock - Devices That Recently Went Offline", sb.String()); err != nil {
		return fmt.Errorf("error sending email: %s", err.Error())
	}

	return nil
}

func getPropertyForEventARN(ctx context.Context, eventARN string) (shared.Property, error) {
	propName := ""
	for k, v := range queueToPropertyName {
		if strings.Contains(eventARN, k) {
			propName = v
		}
	}
	if propName == "" {
		return shared.Property{}, fmt.Errorf("unable to get property name for \"%s\"", eventARN)
	}

	props, err := property.List(ctx)
	if err != nil {
		return shared.Property{}, fmt.Errorf("error loading properties: %s", err.Error())
	}

	prop := shared.Property{}
	for _, p := range props {
		if strings.Contains(strings.ToLower(p.Name), propName) {
			prop = p
		}
	}
	if prop.ID == uuid.Nil {
		return shared.Property{}, fmt.Errorf("error getting property for \"%s\" -> \"%s\": %s", eventARN, propName, err.Error())
	}

	return prop, nil
}

func getQueueForProperty(property shared.Property) (string, error) {
	propName := strings.ToLower(property.Name)
	for k, v := range queueToPropertyName {
		if propName == v {
			return k, nil
		}
	}
	return "", fmt.Errorf("unable to find property %s", propName)
}
