package shared

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDevice_HasConflictingManagedLockCode_SimpleCases(t *testing.T) {
	device := Device{
		ManagedLockCodes: []DeviceManagedLockCode{
			{
				StartAt: getDate(t, "2021-01-02 00:00"),
				EndAt:   getDate(t, "2021-01-03 00:00"),
			},
		},
	}

	// Just before the existing range (outside hour buffer).
	hasConflict := device.HasConflictingManagedLockCode(DeviceManagedLockCode{
		StartAt: getDate(t, "2021-01-01 00:00"),
		EndAt:   getDate(t, "2021-01-01 22:00"),
	})
	assert.False(t, hasConflict)

	// Overlapping start range (within hour buffer).
	hasConflict = device.HasConflictingManagedLockCode(DeviceManagedLockCode{
		StartAt: getDate(t, "2021-01-01 00:00"),
		EndAt:   getDate(t, "2021-01-01 23:01"),
	})
	assert.True(t, hasConflict)

	// Just after the existing range (outside hour buffer).
	hasConflict = device.HasConflictingManagedLockCode(DeviceManagedLockCode{
		StartAt: getDate(t, "2021-01-03 01:00"),
		EndAt:   getDate(t, "2021-01-04 00:00"),
	})
	assert.False(t, hasConflict)

	// Overlapping end range (within hour buffer).
	hasConflict = device.HasConflictingManagedLockCode(DeviceManagedLockCode{
		StartAt: getDate(t, "2021-01-03 00:59"),
		EndAt:   getDate(t, "2021-01-04 00:00"),
	})
	assert.True(t, hasConflict)
}

func getDate(t *testing.T, date string) time.Time {
	dt, err := time.Parse("2006-01-02 15:04", date)
	assert.Nil(t, err)
	return dt
}
