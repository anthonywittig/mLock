package shared

import (
	"time"

	"github.com/google/uuid"
)

type DeviceManagedLockCode struct {
	Code          string                      `json:"code"`
	EndAt         time.Time                   `json:"endAt"`
	ID            uuid.UUID                   `json:"id"`
	Note          string                      `json:"note"`
	ReservationID string                      `json:"reservationId"`
	Status        DeviceManagedLockCodeStatus `json:"status"`
	StartAt       time.Time                   `json:"startAt"`
}

type DeviceManagedLockCodeStatus string

const (
	// We add a 1, 2, ..., n to signify that the state should only progress.
	DeviceManagedLockCodeStatus1Scheduled DeviceManagedLockCodeStatus = "Scheduled"
	DeviceManagedLockCodeStatus2Adding    DeviceManagedLockCodeStatus = "Adding"
	DeviceManagedLockCodeStatus3Enabled   DeviceManagedLockCodeStatus = "Enabled"
	DeviceManagedLockCodeStatus4Removing  DeviceManagedLockCodeStatus = "Removing"
	DeviceManagedLockCodeStatus5Complete  DeviceManagedLockCodeStatus = "Complete"
)

func (m *DeviceManagedLockCode) HasEnded(now time.Time) bool {
	return now.After(m.EndAt)
}

func (m *DeviceManagedLockCode) HasStarted(now time.Time) bool {
	return now.After(m.StartAt)
}

func (m *DeviceManagedLockCode) CodeShouldBePresent(now time.Time) bool {
	if !m.HasStarted(now) {
		return false
	}

	return !m.HasEnded(now)
}
