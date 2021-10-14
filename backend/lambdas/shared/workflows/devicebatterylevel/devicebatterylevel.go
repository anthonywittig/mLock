package devicebatterylevel

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/device"
	mshared "mlock/shared"
	"mlock/shared/protos/messaging"
	"mlock/shared/sqs"
	"time"
)

/*
	A: kick off updates
		- requests all `items` as an SQS message (items have the battery levels) - handled by `B`
		- for each device that has a battery level channel but no linked item, kick off the item creation flow - handled by `C`
	B: process `items` - update all devices' battery levels
	C: process the `item` creation, kick off a `link` creation - handled by `D`
	D: process the `link` creation

*/

func KickOffUpdates(ctx context.Context, sqsClient *sqs.Client, queueNames []string) error {
	// Kick off request for all `items`, the response to this message will update the battery levels.
	for _, qn := range queueNames {

		// TODO: remove when we're ready to test everywhere.
		if qn != "test-in.fifo" {
			continue
		}

		log.Printf("adding list things message to \"%s\"\n", qn)
		if err := sqsClient.SendMessage(
			ctx,
			qn,
			mshared.HabCommandListItems(fmt.Sprintf("hello there @ %s - requesting an item list", time.Now().String())),
		); err != nil {
			return fmt.Errorf("error sending message: %s", err.Error())
		}
	}

	// Loop through all devices and kick off `item` creation where needed.
	ds, err := device.List(ctx)
	if err != nil {
		return fmt.Errorf("error getting devices: %s", err.Error())
	}

	for _, d := range ds {
		// TODO: remove hardcoded property ID when we're ready to release.
		if d.PropertyID.String() != "e0cb7d64-9e51-46db-ab55-ee7ff927baa3" {
			continue
		}
		// TODO...
	}

	return nil
}

func ProcessListItems(ctx context.Context, property shared.Property, in *messaging.OnPremHabCommandResponse) error {
	items := []shared.HABItem{}
	if err := json.Unmarshal(in.Response, &items); err != nil {
		return fmt.Errorf("error parsing json: %s", err.Error())
	}

	// Preprocess items to make them easier to work with later.
	itemByName := map[string]shared.HABItem{}
	for _, item := range items {
		if _, ok := itemByName[item.Name]; ok {
			return fmt.Errorf("we expect unique item names but found a duplicate for: %s", item.Name)
		}
		itemByName[item.Name] = item
	}

	ds, err := device.List(ctx)
	if err != nil {
		return fmt.Errorf("error getting devices: %s", err.Error())
	}

	for _, d := range ds {
		if d.PropertyID != property.ID {
			continue
		}

		if didUpdate := d.UpdateBatteryLevel(itemByName); didUpdate {
			if _, err := device.Put(ctx, d); err != nil {
				return fmt.Errorf("error putting device: %s", err.Error())
			}
		}
	}

	return nil
}

/*
func kickOffBatteryLinkIfNeeded(ctx context.Context, device shared.Device, property shared.Property) error {
	// TODO: remove hardcoded property ID when we're ready to release.
	if device.PropertyID.String() == "e0cb7d64-9e51-46db-ab55-ee7ff927baa3" {
		for _, c := range device.HABThing.Channels {
			if c.ID != "battery-level" {
				continue
			}
			if len(c.LinkedItems) != 0 {
				return nil
			}

			// kick off link creation

			s, err := sqs.GetClient(ctx)
			if err != nil {
				return fmt.Errorf("error getting sqs client: %s", err.Error())
			}

			qn, err := getQueueForProperty(property)
			if err != nil {
				return fmt.Errorf("error getting queue for property %s: %s", property.Name, err.Error())
			}

			command, err := mshared.HabCommandCreateBatteryLevelItem(
				fmt.Sprintf("hello there @ %s - requesting a battery level item", time.Now().String()),
				c.UID,
			)
			if err != nil {
				return fmt.Errorf("error generating command: %s", err.Error())
			}

			if err := s.SendMessage(ctx, qn, command); err != nil {
				return fmt.Errorf("error sending message: %s", err.Error())
			}
		}
	}

	return nil
}
*/
