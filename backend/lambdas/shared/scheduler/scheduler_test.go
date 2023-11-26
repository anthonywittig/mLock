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

	s, dr, now, rr, ur := newScheduler(t, time.Now())

	ctx := context.Background()
	unit := shared.Unit{
		ID: uuid.New(),
	}
	device := shared.Device{
		ID:     uuid.New(),
		UnitID: &unit.ID,
	}
	pastReservation := shared.Reservation{
		ID:                "pastReservation",
		Start:             now.Add(-24 * time.Hour * 180),
		End:               now.Add(-24 * time.Hour * 179),
		DoorCode:          "3344",
		TransactionNumber: "11223344",
	}
	reservation := shared.Reservation{
		ID:                "currentReservation",
		Start:             now,
		End:               now.Add(1 * time.Hour),
		DoorCode:          "9876", // Different from the `TransactionNumber`.
		TransactionNumber: "12345678",
	}
	futureReservation := shared.Reservation{
		ID:                "futureReservation",
		Start:             now.Add(24 * time.Hour * 180),
		End:               now.Add(24 * time.Hour * 181),
		DoorCode:          "3456",
		TransactionNumber: "90123456",
	}

	ur.EXPECT().List(ctx).Return(
		[]shared.Unit{unit},
		nil,
	)

	rr.EXPECT().GetForUnits(ctx, []shared.Unit{unit}).Return(
		map[uuid.UUID][]shared.Reservation{
			unit.ID: {pastReservation, reservation, futureReservation},
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

		assert.Equal(t, 2, len(d.ManagedLockCodes))

		mlc := d.ManagedLockCodes[0]
		assert.Equal(t, "9876", mlc.Code)
		assert.Equal(t, reservation.ID, mlc.Reservation.ID)
		assert.Equal(t, true, mlc.Reservation.Sync)
		assert.Equal(t, reservation.Start.Add(-2*time.Hour), mlc.StartAt)
		assert.Equal(t, reservation.End.Add(30*time.Minute), mlc.EndAt)
		assert.Equal(t, shared.DeviceManagedLockCodeStatus1Scheduled, mlc.Status)

		mlc = d.ManagedLockCodes[1]
		assert.Equal(t, "3456", mlc.Code)
		assert.Equal(t, futureReservation.ID, mlc.Reservation.ID)
		assert.Equal(t, true, mlc.Reservation.Sync)
		assert.Equal(t, futureReservation.Start.Add(-2*time.Hour), mlc.StartAt)
		assert.Equal(t, futureReservation.End.Add(30*time.Minute), mlc.EndAt)
		assert.Equal(t, shared.DeviceManagedLockCodeStatus1Scheduled, mlc.Status)
	}).Return(nil)

	dr.EXPECT().Put(
		ctx,
		gomock.Any(),
	).Do(func(ctx context.Context, d shared.Device) {
		assert.Equal(t, device.ID, d.ID)
		assert.Equal(t, 2, len(d.ManagedLockCodes))
	})

	err := s.ReconcileReservationsAndLockCodes(ctx)
	assert.Nil(t, err)
}

func Test_noEditMLC(t *testing.T) {
	// Don't edit a MLC based on a reservation.

	s, dr, now, rr, ur := newScheduler(t, time.Now())

	ctx := context.Background()
	unit := shared.Unit{
		ID: uuid.New(),
	}
	reservation := shared.Reservation{
		ID:                "someReservationID",
		Start:             now,
		End:               now.Add(1 * time.Hour),
		DoorCode:          "5678",
		TransactionNumber: "12345678",
	}
	managedLockCode := &shared.DeviceManagedLockCode{
		ID: uuid.New(),
		Reservation: shared.DeviceManagedLockCodeReservation{
			ID:   reservation.ID,
			Sync: true,
		},
		Code:    "5678",
		StartAt: reservation.Start.Add(-2 * time.Hour),
		EndAt:   reservation.End.Add(30 * time.Minute),
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

func Test_noSyncMLC(t *testing.T) {
	// Don't sync a MLC based on a reservation.

	s, dr, now, rr, ur := newScheduler(t, time.Now())

	ctx := context.Background()
	unit := shared.Unit{
		ID: uuid.New(),
	}
	reservation := shared.Reservation{
		ID:                "someReservationID",
		Start:             now,
		End:               now.Add(1 * time.Hour),
		DoorCode:          "5678",
		TransactionNumber: "12345678",
	}
	managedLockCode := &shared.DeviceManagedLockCode{
		ID: uuid.New(),
		Reservation: shared.DeviceManagedLockCodeReservation{
			ID:   reservation.ID,
			Sync: false, // Because we're not syncing, we shouldn't see any changes.
		},
		Code:    "5678",
		StartAt: reservation.Start.Add(-1 * time.Hour * 24), // Starts a day earlier than the reservation says.
		EndAt:   reservation.End.Add(1 * time.Minute * 24),  // Ends a day later than the reservation says.
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

	s, dr, now, rr, ur := newScheduler(t, time.Now())

	ctx := context.Background()
	unit := shared.Unit{
		ID: uuid.New(),
	}
	reservation := shared.Reservation{
		ID:                "someReservationID",
		Start:             now,
		End:               now.Add(1 * time.Hour),
		DoorCode:          "5678",
		TransactionNumber: "12345678",
	}
	managedLockCode := &shared.DeviceManagedLockCode{
		// The code, start, and end don't match (should probably test them individually).
		ID: uuid.New(),
		Reservation: shared.DeviceManagedLockCodeReservation{
			ID:   reservation.ID,
			Sync: true,
		},
		Code:    "1111",                                  // Should be "5678".
		StartAt: reservation.Start.Add(-5 * time.Minute), // Should be -60?
		EndAt:   reservation.End.Add(5 * time.Minute),    // Should be +30.
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
		assert.Equal(t, reservation.ID, mlc.Reservation.ID)
		assert.Equal(t, true, mlc.Reservation.Sync)
		assert.Equal(t, reservation.Start.Add(-2*time.Hour), mlc.StartAt)
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

	s, dr, now, rr, ur := newScheduler(t, time.Now())

	ctx := context.Background()
	unit := shared.Unit{
		ID: uuid.New(),
	}
	reservation := shared.Reservation{
		ID:                "someReservationID",
		Start:             now.Add(-1 * time.Hour),
		End:               now.Add(-31 * time.Minute), // To make it end a minute ago, we need to subtract the 30 that will get added...
		DoorCode:          "5678",
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

	s, dr, now, rr, ur := newScheduler(t, time.Now())

	ctx := context.Background()
	unit := shared.Unit{
		ID: uuid.New(),
	}
	managedLockCode := &shared.DeviceManagedLockCode{
		ID: uuid.New(),
		Reservation: shared.DeviceManagedLockCodeReservation{
			ID:   "someReservationIDThatDoesNotExist",
			Sync: true,
		},
		Code:    "1111",
		StartAt: now.Add(-2 * time.Minute),
		EndAt:   now.Add(-1 * time.Minute),
		Status:  shared.DeviceManagedLockCodeStatus4Removing,
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

	type TestCase struct {
		OriginalStartAt time.Time
		OriginalEndAt   time.Time
		OriginalStatus  shared.DeviceManagedLockCodeStatus
	}

	now := time.Now()

	for _, testCase := range []TestCase{
		// Reservation has ended.
		{
			OriginalStartAt: now.Add(-10 * time.Hour),
			OriginalEndAt:   now.Add(-5 * time.Hour),
			OriginalStatus:  shared.DeviceManagedLockCodeStatus5Complete,
		},
		// Reservation is in progress.
		{
			OriginalStartAt: now.Add(-10 * time.Hour),
			OriginalEndAt:   now.Add(10 * time.Hour),
			OriginalStatus:  shared.DeviceManagedLockCodeStatus3Enabled,
		},
		// Reservation hasn't started.
		{
			OriginalStartAt: now.Add(10 * time.Hour),
			OriginalEndAt:   now.Add(20 * time.Hour),
			OriginalStatus:  shared.DeviceManagedLockCodeStatus1Scheduled,
		},
	} {
		s, dr, _, rr, ur := newScheduler(t, now)

		ctx := context.Background()
		unit := shared.Unit{
			ID: uuid.New(),
		}
		managedLockCode := &shared.DeviceManagedLockCode{
			ID: uuid.New(),
			Reservation: shared.DeviceManagedLockCodeReservation{
				ID:   "someReservationIDThatDoesNotExist",
				Sync: false,
			},
			Code:    "1111",
			StartAt: testCase.OriginalStartAt,
			EndAt:   testCase.OriginalEndAt,
			Status:  testCase.OriginalStatus,
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

		err := s.ReconcileReservationsAndLockCodes(ctx)
		assert.Nil(t, err)
	}
}

func newScheduler(t *testing.T, now time.Time) (*scheduler.Scheduler, *mock_scheduler.MockDeviceRepository, time.Time, *mock_scheduler.MockReservationRepository, *mock_scheduler.MockUnitRepository) {
	ctrl := gomock.NewController(t)

	dr := mock_scheduler.NewMockDeviceRepository(ctrl)
	rr := mock_scheduler.NewMockReservationRepository(ctrl)
	ur := mock_scheduler.NewMockUnitRepository(ctrl)

	s := scheduler.NewScheduler(dr, now, rr, ur)
	return s, dr, now, rr, ur
}
