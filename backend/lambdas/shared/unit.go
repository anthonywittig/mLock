package shared

import (
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
	At   struct {
		Occupied         bool                    `json:"occupied"`
		ManagedLockCodes []DeviceManagedLockCode `json:"managedLockCodes"`
	} `json:"at"`
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

func (u *Unit) OccupancyStatusForDay(devices []Device, at time.Time) UnitOccupancyStatus {
	year, month, day := at.Date()
	date := time.Date(year, month, day, 0, 0, 0, 0, at.Location())

	noon := date.Add(12 * time.Hour)
	fourPM := date.Add(16 * time.Hour)

	unitOccupiedStatus := UnitOccupancyStatus{
		Date: date,
	}

	for _, d := range devices {
		if d.UnitID != nil && *d.UnitID == u.ID {
			for _, mlc := range d.ManagedLockCodes {
				if mlc.Reservation.ID != "" {
					reservationRealStart := mlc.StartAt.Add(-1 * time.Duration(ReservationStartBufferInMinutes) * time.Minute)
					// Warning: we're basing these off of the lock start/stop times which should have a buffer around when we expect the unit to be occupied. If we adjust that buffer we might break this logic.
					if (reservationRealStart.Before(at) || reservationRealStart == at) && mlc.EndAt.After(at) {
						unitOccupiedStatus.At.Occupied = true
						unitOccupiedStatus.At.ManagedLockCodes = append(unitOccupiedStatus.At.ManagedLockCodes, *mlc)
					}
					if (reservationRealStart.Before(noon) || reservationRealStart == noon) && mlc.EndAt.After(noon) {
						unitOccupiedStatus.Noon.Occupied = true
						unitOccupiedStatus.Noon.ManagedLockCodes = append(unitOccupiedStatus.Noon.ManagedLockCodes, *mlc)
					}
					if (reservationRealStart.Before(fourPM) || reservationRealStart == fourPM) && mlc.EndAt.After(fourPM) {
						unitOccupiedStatus.FourPM.Occupied = true
						unitOccupiedStatus.FourPM.ManagedLockCodes = append(unitOccupiedStatus.FourPM.ManagedLockCodes, *mlc)
					}
				}
			}
		}
	}

	return unitOccupiedStatus
}
