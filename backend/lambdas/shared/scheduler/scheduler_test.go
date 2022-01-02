package scheduler_test

//go:generate mockgen -source=scheduler.go -destination mocks/mock_scheduler/scheduler.go

import (
	"context"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/scheduler"
	"mlock/lambdas/shared/scheduler/mocks/mock_scheduler"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_addMLC(t *testing.T) {
	// Add a MLC based on a reservation.

	s, dr, now, rr, ur := newScheduler(t)

	ctx := context.Background()
	unit := shared.Unit{
		ID:          uuid.New(),
		CalendarURL: "notBlank",
	}
	device := shared.Device{
		ID:     uuid.New(),
		UnitID: &unit.ID,
	}
	reservation := shared.Reservation{
		ID:                "someReservationID",
		Start:             now,
		End:               now.Add(1 * time.Hour),
		TransactionNumber: "12345678",
	}

	ur.EXPECT().List(ctx).Return(
		[]shared.Unit{unit},
		nil,
	)

	rr.EXPECT().GetForUnits(ctx, []shared.Unit{unit}).Return(
		map[uuid.UUID][]shared.Reservation{
			unit.ID: {reservation},
		},
		nil,
	)

	dr.EXPECT().List(ctx).Return(
		[]shared.Device{device},
		nil,
	)

	dr.EXPECT().AppendToAuditLog(
		ctx,
		gomock.Any(),
		gomock.Any(),
	).Do(func(ctx context.Context, d shared.Device, managedLockCodes []*shared.DeviceManagedLockCode) {
		assert.Equal(t, device.ID, d.ID)
		assert.Equal(t, d.ManagedLockCodes, managedLockCodes)

		assert.Equal(t, 1, len(d.ManagedLockCodes))
		mlc := d.ManagedLockCodes[0]
		assert.Equal(t, "5678", mlc.Code)
		assert.Equal(t, reservation.ID, mlc.ReservationID)
		assert.Equal(t, reservation.Start.Add(-30*time.Minute), mlc.StartAt)
		assert.Equal(t, reservation.End.Add(30*time.Minute), mlc.EndAt)
		assert.Equal(t, shared.DeviceManagedLockCodeStatus1Scheduled, mlc.Status)
	}).Return(nil)

	dr.EXPECT().Put(
		ctx,
		gomock.Any(),
	).Do(func(ctx context.Context, d shared.Device) {
		assert.Equal(t, device.ID, d.ID)
		assert.Equal(t, 1, len(d.ManagedLockCodes))
	})

	err := s.ReconcileReservationsAndLockCodes(ctx)
	assert.Nil(t, err)
}

func Test_noEditMLC(t *testing.T) {
	// Don't edit a MLC based on a reservation.

	s, dr, now, rr, ur := newScheduler(t)

	ctx := context.Background()
	unit := shared.Unit{
		ID:          uuid.New(),
		CalendarURL: "notBlank",
	}
	reservation := shared.Reservation{
		ID:                "someReservationID",
		Start:             now,
		End:               now.Add(1 * time.Hour),
		TransactionNumber: "12345678",
	}
	managedLockCode := &shared.DeviceManagedLockCode{
		ID:            uuid.New(),
		ReservationID: reservation.ID,
		Code:          "5678",
		StartAt:       reservation.Start.Add(-30 * time.Minute),
		EndAt:         reservation.End.Add(30 * time.Minute),
	}
	device := shared.Device{
		ID:               uuid.New(),
		ManagedLockCodes: []*shared.DeviceManagedLockCode{managedLockCode},
		UnitID:           &unit.ID,
	}

	ur.EXPECT().List(ctx).Return(
		[]shared.Unit{unit},
		nil,
	)

	rr.EXPECT().GetForUnits(ctx, []shared.Unit{unit}).Return(
		map[uuid.UUID][]shared.Reservation{
			unit.ID: {reservation},
		},
		nil,
	)

	dr.EXPECT().List(ctx).Return(
		[]shared.Device{device},
		nil,
	)

	// No dr.AppendToAuditLog or dr.Put because there are no modifications.

	err := s.ReconcileReservationsAndLockCodes(ctx)
	assert.Nil(t, err)
}

func Test_editMLC(t *testing.T) {
	// Edit a MLC based on a reservation.

	s, dr, now, rr, ur := newScheduler(t)

	ctx := context.Background()
	unit := shared.Unit{
		ID:          uuid.New(),
		CalendarURL: "notBlank",
	}
	reservation := shared.Reservation{
		ID:                "someReservationID",
		Start:             now,
		End:               now.Add(1 * time.Hour),
		TransactionNumber: "12345678",
	}
	managedLockCode := &shared.DeviceManagedLockCode{
		// The code, start, and end don't match (should probably test them individually).
		ID:            uuid.New(),
		ReservationID: reservation.ID,
		Code:          "1111",                                  // Should be "5678".
		StartAt:       reservation.Start.Add(-5 * time.Minute), // Should be -30.
		EndAt:         reservation.End.Add(5 * time.Minute),    // Should be +30.
	}
	device := shared.Device{
		ID:               uuid.New(),
		ManagedLockCodes: []*shared.DeviceManagedLockCode{managedLockCode},
		UnitID:           &unit.ID,
	}

	ur.EXPECT().List(ctx).Return(
		[]shared.Unit{unit},
		nil,
	)

	rr.EXPECT().GetForUnits(ctx, []shared.Unit{unit}).Return(
		map[uuid.UUID][]shared.Reservation{
			unit.ID: {reservation},
		},
		nil,
	)

	dr.EXPECT().List(ctx).Return(
		[]shared.Device{device},
		nil,
	)

	dr.EXPECT().AppendToAuditLog(
		ctx,
		gomock.Any(),
		gomock.Any(),
	).Do(func(ctx context.Context, d shared.Device, managedLockCodes []*shared.DeviceManagedLockCode) {
		assert.Equal(t, device.ID, d.ID)
		assert.Equal(t, d.ManagedLockCodes, managedLockCodes)

		assert.Equal(t, 1, len(d.ManagedLockCodes))
		mlc := d.ManagedLockCodes[0]
		assert.Equal(t, "5678", mlc.Code)
		assert.Equal(t, reservation.ID, mlc.ReservationID)
		assert.Equal(t, reservation.Start.Add(-30*time.Minute), mlc.StartAt)
		assert.Equal(t, reservation.End.Add(30*time.Minute), mlc.EndAt)
	}).Return(nil)

	dr.EXPECT().Put(
		ctx,
		gomock.Any(),
	).Do(func(ctx context.Context, d shared.Device) {
		assert.Equal(t, device.ID, d.ID)
		assert.Equal(t, 1, len(d.ManagedLockCodes))
	})

	err := s.ReconcileReservationsAndLockCodes(ctx)
	assert.Nil(t, err)
}

func Test_recentlyEndedReservation(t *testing.T) {
	// Do nothing for a reservation that has recently ended and doesn't have a corresponding MLC.

	s, dr, now, rr, ur := newScheduler(t)

	ctx := context.Background()
	unit := shared.Unit{
		ID:          uuid.New(),
		CalendarURL: "notBlank",
	}
	reservation := shared.Reservation{
		ID:                "someReservationID",
		Start:             now.Add(-1 * time.Hour),
		End:               now.Add(-31 * time.Minute), // To make it end a minute ago, we need to subtract the 30 that will get added...
		TransactionNumber: "12345678",
	}
	device := shared.Device{
		ID:               uuid.New(),
		ManagedLockCodes: []*shared.DeviceManagedLockCode{},
		UnitID:           &unit.ID,
	}

	ur.EXPECT().List(ctx).Return(
		[]shared.Unit{unit},
		nil,
	)

	rr.EXPECT().GetForUnits(ctx, []shared.Unit{unit}).Return(
		map[uuid.UUID][]shared.Reservation{
			unit.ID: {reservation},
		},
		nil,
	)

	dr.EXPECT().List(ctx).Return(
		[]shared.Device{device},
		nil,
	)

	// No dr.AppendToAuditLog or dr.Put because there are no modifications.

	err := s.ReconcileReservationsAndLockCodes(ctx)
	assert.Nil(t, err)
}

func Test_recentlyEndedMLC(t *testing.T) {
	// Do nothing for a MLC that has recently ended and doesn't have a corresponding reservation.

	s, dr, now, rr, ur := newScheduler(t)

	ctx := context.Background()
	unit := shared.Unit{
		ID:          uuid.New(),
		CalendarURL: "notBlank",
	}
	managedLockCode := &shared.DeviceManagedLockCode{
		ID:            uuid.New(),
		ReservationID: "someReservationIDThatDoesn'tExist",
		Code:          "1111",
		StartAt:       now.Add(-2 * time.Minute),
		EndAt:         now.Add(-1 * time.Minute),
	}
	device := shared.Device{
		ID:               uuid.New(),
		ManagedLockCodes: []*shared.DeviceManagedLockCode{managedLockCode},
		UnitID:           &unit.ID,
	}

	ur.EXPECT().List(ctx).Return(
		[]shared.Unit{unit},
		nil,
	)

	rr.EXPECT().GetForUnits(ctx, []shared.Unit{unit}).Return(
		map[uuid.UUID][]shared.Reservation{},
		nil,
	)

	dr.EXPECT().List(ctx).Return(
		[]shared.Device{device},
		nil,
	)

	// No dr.AppendToAuditLog or dr.Put because there are no modifications.

	err := s.ReconcileReservationsAndLockCodes(ctx)
	assert.Nil(t, err)
}

func Test_editMLCWithNoReservation(t *testing.T) {
	// Edit a MLC when it's reservation doesn't exist anymore.

	s, dr, now, rr, ur := newScheduler(t)

	ctx := context.Background()
	unit := shared.Unit{
		ID:          uuid.New(),
		CalendarURL: "notBlank",
	}
	managedLockCode := &shared.DeviceManagedLockCode{
		// The code, start, and end don't match (should probably test them individually).
		ID:            uuid.New(),
		ReservationID: "someReservationIDThatDoesn'tExist",
		Code:          "1111",
		StartAt:       now.Add(-5 * time.Hour),
		EndAt:         now.Add(5 * time.Hour),
	}
	device := shared.Device{
		ID:               uuid.New(),
		ManagedLockCodes: []*shared.DeviceManagedLockCode{managedLockCode},
		UnitID:           &unit.ID,
	}

	ur.EXPECT().List(ctx).Return(
		[]shared.Unit{unit},
		nil,
	)

	rr.EXPECT().GetForUnits(ctx, []shared.Unit{unit}).Return(
		map[uuid.UUID][]shared.Reservation{},
		nil,
	)

	dr.EXPECT().List(ctx).Return(
		[]shared.Device{device},
		nil,
	)

	dr.EXPECT().AppendToAuditLog(
		ctx,
		gomock.Any(),
		gomock.Any(),
	).Do(func(ctx context.Context, d shared.Device, managedLockCodes []*shared.DeviceManagedLockCode) {
		assert.Equal(t, device.ID, d.ID)
		assert.Equal(t, d.ManagedLockCodes, managedLockCodes)

		assert.Equal(t, 1, len(d.ManagedLockCodes))
		mlc := d.ManagedLockCodes[0]
		assert.Equal(t, "1111", mlc.Code)
		assert.Equal(t, managedLockCode.ReservationID, mlc.ReservationID)

		// The start time stays the same but the end time changes.
		assert.Equal(t, managedLockCode.StartAt, mlc.StartAt)
		assert.Equal(t, now.Add(1*time.Hour), mlc.EndAt)
	}).Return(nil)

	dr.EXPECT().Put(
		ctx,
		gomock.Any(),
	).Do(func(ctx context.Context, d shared.Device) {
		assert.Equal(t, device.ID, d.ID)
		assert.Equal(t, 1, len(d.ManagedLockCodes))
	})

	err := s.ReconcileReservationsAndLockCodes(ctx)
	assert.Nil(t, err)
}

func newScheduler(t *testing.T) (*scheduler.Scheduler, *mock_scheduler.MockDeviceRepository, time.Time, *mock_scheduler.MockReservationRepository, *mock_scheduler.MockUnitRepository) {
	ctrl := gomock.NewController(t)

	dr := mock_scheduler.NewMockDeviceRepository(ctrl)
	now := time.Now()
	rr := mock_scheduler.NewMockReservationRepository(ctrl)
	ur := mock_scheduler.NewMockUnitRepository(ctrl)

	s := scheduler.NewScheduler(dr, now, rr, ur)
	return s, dr, now, rr, ur
}
