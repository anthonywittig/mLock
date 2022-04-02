package shared

import (
	"sort"
	"time"

	"github.com/google/uuid"
)

type Device struct {
	Battery struct {
		LastUpdatedAt *time.Time `json:"lastUpdatedAt"`
		Level         string     `json:"level"` // Could probably do a numeric type, but this simplifies some things (e.g. "NAN").
	} `json:"battery"`
	ControllerID      string                   `json:"controllerId"`
	History           []DeviceHistory          `json:"history"`
	ID                uuid.UUID                `json:"id"`
	LastRefreshedAt   time.Time                `json:"lastRefreshedAt"`
	LastWentOfflineAt *time.Time               `json:"lastWentOfflineAt"`
	LastWentOnlineAt  *time.Time               `json:"lastWentOnlineAt"`
	ManagedLockCodes  []*DeviceManagedLockCode `json:"managedLockCodes"`
	RawDevice         RawDevice                `json:"rawDevice"`
	UnitID            *uuid.UUID               `json:"unitId"`
}

type DeviceHistory struct {
	Description string    `json:"description"`
	RawDevice   RawDevice `json:"rawDevice"`
	RecordedAt  time.Time `json:"recordedAt"`
}

type RawDevice struct {
	Battery      RawDeviceBattery    `json:"battery"`
	Category     string              `json:"category"`
	DeviceTypeID string              `json:"deviceTypeId"`
	ID           string              `json:"id"`
	LockCodes    []RawDeviceLockCode `json:"lockCodes"`
	Name         string              `json:"name"`
	Status       string              `json:"status"`
}

type RawDeviceBattery struct {
	BatteryPowered bool `json:"batteryPowered"`
	Level          int  `json:"level"`
}
type RawDeviceLockCode struct {
	Code string `json:"code"`
	Mode string `json:"mode"`
	Name string `json:"name"`
	Slot int    `json:"slot"` // TODO: change json to `-` if we're not using it.
}

const (
	DeviceCodeModeEnabled = "enabled"
	DeviceStatusOffline   = "OFFLINE"
	DeviceStatusOnline    = "ONLINE"
)

func (d *Device) GenerateUnmanagedLockCodes() []RawDeviceLockCode {
	umlcs := []RawDeviceLockCode{}

	for _, c := range d.RawDevice.LockCodes {
		found := false
		for _, mlc := range d.ManagedLockCodes {
			if c.Code == mlc.Code {
				found = true
				break
			}
		}
		if !found {
			umlcs = append(umlcs, c)
		}
	}

	sort.Slice(
		umlcs,
		func(a, b int) bool {
			return umlcs[a].Code < umlcs[b].Code
		},
	)

	return umlcs
}

func (d *Device) GetManagedLockCode(id uuid.UUID) *DeviceManagedLockCode {
	for _, mlc := range d.ManagedLockCodes {
		if mlc.ID == id {
			return mlc
		}
	}
	return nil
}
