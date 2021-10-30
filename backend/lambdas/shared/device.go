package shared

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	Battery struct {
		LastUpdatedAt *time.Time `json:"lastUpdatedAt"`
		Level         string     `json:"level"` // Could probably do a numeric type, but this simplifies some things (e.g. "NAN").
	} `json:"battery"`
	History           []DeviceHistory `json:"history"`
	ID                uuid.UUID       `json:"id"`
	LastRefreshedAt   time.Time       `json:"lastRefreshedAt"`
	LastWentOfflineAt *time.Time      `json:"lastWentOfflineAt"`
	PropertyID        uuid.UUID       `json:"propertyId"`
	RawDevice         RawDevice       `json:"rawDevice"`
	UnitID            *uuid.UUID      `json:"unitId"`
}

type DeviceHistory struct {
	Description string    `json:"description"`
	RawDevice   RawDevice `json:"rawDevice"`
	RecordedAt  time.Time `json:"recordedAt"`
}

type RawDevice struct {
	Battery   RawDeviceBattery    `json:"battery"`
	Category  string              `json:"category"`
	ID        string              `json:"id"`
	LockCodes []RawDeviceLockCode `json:"lockCodes"`
	Name      string              `json:"name"`
	Status    string              `json:"status"`
}

type RawDeviceBattery struct {
	BatteryPowered bool `json:"batteryPowered"`
	Level          int  `json:"level"`
}
type RawDeviceLockCode struct {
	Code string `json:"code"`
	Mode string `json:"mode"`
	Name string `json:"name"`
	Slot int    `json:"slot"`
}

const (
	DeviceStatusOffline = "OFFLINE"
	DeviceStatusOnline  = "ONLINE"
)
