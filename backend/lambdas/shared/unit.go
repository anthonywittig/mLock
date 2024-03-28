package shared

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Unit struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	PropertyID        uuid.UUID `json:"propertyId"`
	RemotePropertyURL string    `json:"remotePropertyUrl"`
	UpdatedBy         string    `json:"updatedBy"`
}

type UnitOccupancyStatus struct {
	Date time.Time `json:"date"`
	Noon struct {
		Occupied         bool                    `json:"occupied"`
		ManagedLockCodes []DeviceManagedLockCode `json:"managedLockCodes"`
	} `json:"noon"`
	FourPM struct {
		Occupied         bool                    `json:"occupied"`
		ManagedLockCodes []DeviceManagedLockCode `json:"managedLockCodes"`
	} `json:"fourPm"`
}

func (u *Unit) GetRemotePropertyID() int {
	// RemotePropertyURL is of the form:
	// https://dashboard.hostaway.com/listing/211374
	if u.RemotePropertyURL == "" {
		return -1
	}
	split := strings.Split(u.RemotePropertyURL, "/")
	id := split[len(split)-1]
	intID, err := strconv.Atoi(id)
	if err != nil {
		return -1
	}
	return intID
}

func (u *Unit) OccupancyStatusForDay(devices []Device, day time.Time) (UnitOccupancyStatus, error) {
	if day.Hour() != 0 || day.Minute() != 0 || day.Second() != 0 || day.Nanosecond() != 0 {
		return UnitOccupancyStatus{}, fmt.Errorf("day must have zero time fields")
	}

	noon := day.Add(12 * time.Hour)
	fourPM := day.Add(16 * time.Hour)

	unitOccupiedStatus := UnitOccupancyStatus{
		Date: day,
	}

	for _, d := range devices {
		if d.UnitID != nil && *d.UnitID == u.ID {
			for _, mlc := range d.ManagedLockCodes {
				if mlc.Reservation.ID != "" {
					// Warning: we're basing these off of the lock start/stop times which should have a buffer around when we expect the unit to be occupied. If we adjust that buffer we might break this logic.
					if (mlc.StartAt.Before(noon) || mlc.StartAt == noon) && mlc.EndAt.After(noon) {
						unitOccupiedStatus.Noon.Occupied = true
						unitOccupiedStatus.Noon.ManagedLockCodes = append(unitOccupiedStatus.Noon.ManagedLockCodes, *mlc)
					}
					if (mlc.StartAt.Before(fourPM) || mlc.StartAt == fourPM) && mlc.EndAt.After(fourPM) {
						unitOccupiedStatus.FourPM.Occupied = true
						unitOccupiedStatus.FourPM.ManagedLockCodes = append(unitOccupiedStatus.FourPM.ManagedLockCodes, *mlc)
					}
				}
			}
		}
	}

	return unitOccupiedStatus, nil
}
