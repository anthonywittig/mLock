package shared

import (
	"mlock/lambdas/shared/ezlo"
	"time"

	"github.com/google/uuid"
)

type Device struct {
	Battery struct {
		LastUpdatedAt *time.Time `json:"lastUpdatedAt"`
		Level         string     `json:"level"` // Could probably do a numeric type, but this simplifies some things (e.g. "NAN").
	} `json:"battery"`
	EZLODevice        ezlo.Device     `json:"ezloDevice"`
	HABThing          HABThing        `json:"habThing"`
	History           []DeviceHistory `json:"history"`
	ID                uuid.UUID       `json:"id"`
	LastRefreshedAt   time.Time       `json:"lastRefreshedAt"`
	LastWentOfflineAt *time.Time      `json:"lastWentOfflineAt"`
	PropertyID        uuid.UUID       `json:"propertyId"`
	UnitID            *uuid.UUID      `json:"unitId"`
}

type DeviceHistory struct {
	Description string      `json:"description"`
	EZLODevice  ezlo.Device `json:"ezloDevice"`
	HABThing    HABThing    `json:"habThing"`
	RecordedAt  time.Time   `json:"recordedAt"`
}

func (d *Device) UpdateBatteryLevel(itemByName map[string]HABItem) bool {
	channel := d.HABThing.GetBatteryChannel()
	if channel == nil {
		return false
	}

	for _, link := range channel.LinkedItems {
		// We will probably only have a single linked item.

		item, ok := itemByName[link]
		if !ok {
			// We don't expect this to happen and could return an error...
			// return fmt.Errorf("battery channel had a linked item but we didn't find it: %s", link)
			continue
		}

		n := time.Now()
		d.Battery.LastUpdatedAt = &n
		d.Battery.Level = item.State
		return true
	}

	return false
}
