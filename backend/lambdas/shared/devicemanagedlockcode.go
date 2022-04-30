package shared

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type DeviceManagedLockCode struct {
	Code              string                           `json:"code"`
	EndAt             time.Time                        `json:"endAt"`
	ID                uuid.UUID                        `json:"id"`
	Note              string                           `json:"note"`
	Reservation       DeviceManagedLockCodeReservation `json:"reservation"`
	Status            DeviceManagedLockCodeStatus      `json:"status"`
	StartAt           time.Time                        `json:"startAt"`
	StartedAddingAt   *time.Time                       `json:"startedAddingAt"`
	WasEnabledAt      *time.Time                       `json:"wasEnabledAt"`
	StartedRemovingAt *time.Time                       `json:"startedRemovingAt"`
	WasCompletedAt    *time.Time                       `json:"wasCompletedAt"`
}

type DeviceManagedLockCodeReservation struct {
	ID   string `json:"id"`
	Sync bool   `json:"sync"`
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

func (m *DeviceManagedLockCode) SetStatus(status DeviceManagedLockCodeStatus) error {
	if m.Status == status {
		return nil
	}
	m.Status = status

	now := time.Now()
	if status == DeviceManagedLockCodeStatus1Scheduled {
		m.StartedAddingAt = nil
		m.WasEnabledAt = nil
		m.StartedRemovingAt = nil
		m.WasCompletedAt = nil
	} else if status == DeviceManagedLockCodeStatus2Adding {
		m.StartedAddingAt = &now
		m.WasEnabledAt = nil
		m.StartedRemovingAt = nil
		m.WasCompletedAt = nil
	} else if status == DeviceManagedLockCodeStatus3Enabled {
		m.WasEnabledAt = &now
		m.StartedRemovingAt = nil
		m.WasCompletedAt = nil
	} else if status == DeviceManagedLockCodeStatus4Removing {
		m.StartedRemovingAt = &now
		m.WasCompletedAt = nil
	} else if status == DeviceManagedLockCodeStatus5Complete {
		m.WasCompletedAt = &now
	} else {
		return fmt.Errorf("unhandled status %s", status)
	}
	return nil
}
