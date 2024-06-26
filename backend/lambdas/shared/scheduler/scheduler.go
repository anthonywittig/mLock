package scheduler

import (
	"context"
	"fmt"
	"log"
	"mlock/lambdas/shared"
	"time"

	"github.com/google/uuid"
)

type Scheduler struct {
	dr  DeviceRepository
	now time.Time
	rr  ReservationRepository
	ur  UnitRepository
}

type DeviceRepository interface {
	AppendToAuditLog(ctx context.Context, device shared.Device, managedLockCodes []*shared.DeviceManagedLockCode) error
	List(ctx context.Context) ([]shared.Device, error)
	Put(ctx context.Context, item shared.Device) (shared.Device, error)
}

type ReservationRepository interface {
	GetForUnits(ctx context.Context, units []shared.Unit) (map[uuid.UUID][]shared.Reservation, error)
}

type UnitRepository interface {
	List(ctx context.Context) ([]shared.Unit, error)
}

func NewScheduler(dr DeviceRepository, now time.Time, rr ReservationRepository, ur UnitRepository) *Scheduler {
	return &Scheduler{
		dr:  dr,
		now: now,
		rr:  rr,
		ur:  ur,
	}
}

func (s *Scheduler) ReconcileReservationsAndLockCodes(ctx context.Context) error {
	units, err := s.ur.List(ctx)
	if err != nil {
		return fmt.Errorf("error getting units: %s", err.Error())
	}

	// Could be more selective and only use the units that have a device associated with them, but hopefully that's a minor optimization that doesn't matter.
	reservationsByUnit, err := s.rr.GetForUnits(ctx, units)
	if err != nil {
		return fmt.Errorf("error getting reservations: %s", err.Error())
	}

	devices, err := s.dr.List(ctx)
	if err != nil {
		return fmt.Errorf("error getting devices: %s", err.Error())
	}

	for _, d := range devices {
		if err := s.processDevice(ctx, d, reservationsByUnit); err != nil {
			return fmt.Errorf("error processing device: %s", err.Error())
		}
	}

	return nil
}

func (s *Scheduler) processDevice(ctx context.Context, device shared.Device, reservationsByUnit map[uuid.UUID][]shared.Reservation) error {
	if device.UnitID == nil {
		return nil
	}

	mlcByReservation := map[string]*shared.DeviceManagedLockCode{}
	for _, mlc := range device.ManagedLockCodes {
		if mlc.Reservation.ID != "" {
			mlcByReservation[mlc.Reservation.ID] = mlc
		}
	}

	relevantReservations, err := s.getRelevantReservations(reservationsByUnit[*device.UnitID])
	if err != nil {
		return fmt.Errorf("error getting relevant reservations: %s", err.Error())
	}

	// We want the lock codes to start and end with a buffer.
	for id, r := range relevantReservations {
		r.Start = r.Start.Add(shared.ReservationStartBufferInMinutes * time.Minute)
		r.End = r.End.Add(shared.ReservationEndBufferInMinutes * time.Minute)
		relevantReservations[id] = r
	}

	needToSave := []*shared.DeviceManagedLockCode{}

	for _, reservation := range relevantReservations {
		if reservation.DoorCode == "" {
			return fmt.Errorf("reservation %s has no door code", reservation.ID)
		}
		mlc, ok := mlcByReservation[reservation.ID]
		if !ok {
			if reservation.End.Before(s.now) {
				continue // Let's not get in a fight with some other part of the system that's removing these and just ignore it.
			}

			newMLC := &shared.DeviceManagedLockCode{
				Code:  reservation.DoorCode,
				EndAt: reservation.End,
				ID:    uuid.New(),
				Note:  fmt.Sprintf("Automatically created for reservation %s", reservation.TransactionNumber),
				Reservation: shared.DeviceManagedLockCodeReservation{
					ID:   reservation.ID,
					Sync: true,
				},
				Status:  shared.DeviceManagedLockCodeStatus1Scheduled,
				StartAt: reservation.Start,
			}

			device.ManagedLockCodes = append(device.ManagedLockCodes, newMLC)
			needToSave = append(needToSave, newMLC)
		} else if mlc.Reservation.Sync {
			changedFields := []string{}
			if mlc.Code != reservation.DoorCode {
				changedFields = append(changedFields, "code")
				mlc.Code = reservation.DoorCode
			}
			if !mlc.StartAt.Equal(reservation.Start) {
				changedFields = append(changedFields, "start")
				mlc.StartAt = reservation.Start
			}
			if !mlc.EndAt.Equal(reservation.End) {
				changedFields = append(changedFields, "end")
				mlc.EndAt = reservation.End
			}

			if len(changedFields) > 0 {
				mlc.Note = fmt.Sprintf("Updating to match reservation (fields: %v)", changedFields)
				needToSave = append(needToSave, mlc)
			}
		}
	}

	for _, mlc := range device.ManagedLockCodes {
		if mlc.Reservation.ID != "" && mlc.Reservation.Sync {
			if _, ok := relevantReservations[mlc.Reservation.ID]; !ok {
				if mlc.Status == shared.DeviceManagedLockCodeStatus1Scheduled {
					log.Printf("DEBUG: canceling reservation %s from device %s", mlc.Reservation.ID, device.RawDevice.Name)
					mlc.Note = "Reservation disappeared, assuming it was canceled; moving the start and end times to now"
					mlc.StartAt = s.now
					mlc.EndAt = s.now
					needToSave = append(needToSave, mlc)
				} else if mlc.Status == shared.DeviceManagedLockCodeStatus2Adding {
					log.Printf("DEBUG: canceling reservation %s from device %s", mlc.Reservation.ID, device.RawDevice.Name)
					mlc.Note = "Reservation disappeared, assuming it was canceled; moving the end time to now"
					mlc.EndAt = s.now
					needToSave = append(needToSave, mlc)
				} else if mlc.Status == shared.DeviceManagedLockCodeStatus3Enabled {
					// The reservation drops off the calendar at midnight (or there about) the day before it ends.
				} else if mlc.Status == shared.DeviceManagedLockCodeStatus4Removing || mlc.Status == shared.DeviceManagedLockCodeStatus5Complete {
					// Do nothing.
				} else {
					return fmt.Errorf("unhandled status %s", mlc.Status)
				}
			}
		}
	}

	if len(needToSave) > 0 {
		if err := s.dr.AppendToAuditLog(ctx, device, needToSave); err != nil {
			return fmt.Errorf("error appending to audit log: %s", err.Error())
		}

		_, err = s.dr.Put(ctx, device)
		if err != nil {
			return fmt.Errorf("error updating device: %s", err.Error())
		}
	}

	return nil
}

func (s *Scheduler) getRelevantReservations(reservations []shared.Reservation) (map[string]shared.Reservation, error) {
	relevantReservations := map[string]shared.Reservation{}
	for _, r := range reservations {
		if r.End.Before(s.now.Add(-1 * time.Hour)) {
			continue // If it ended an hour or more ago, ignore it.
		}

		// We don't control the IDs, so it's probably a good idea to double check that they're unique (at least within a unit's active reservations).
		if _, ok := relevantReservations[r.ID]; ok {
			return nil, fmt.Errorf("duplicate reservation found, reservation ID: %s", r.ID)
		}

		relevantReservations[r.ID] = r
	}

	return relevantReservations, nil
}
