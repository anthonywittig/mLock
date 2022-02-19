package lockengine_test

//go:generate mockgen -source=lockengine.go -destination mocks/mock_lockengine/lockengine.go

import (
	"context"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/lockengine"
	"mlock/lambdas/shared/lockengine/mocks/mock_lockengine"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// The organization of these tests is not good, consider breaking them out into separate files.
//
// Some tests to consider:
// * do nothing
//   * no managed lock codes
//   * managed lock code but no action to take
//     * codes already exist
//     * codes don't exist and don't need to be created
// * add lock codes
// * remove lock codes
// * add and remove lock codes
// * MLC is "active" but the code doesn't exist
//   * Should really add tests for all lock code states and verify when adding/removing.

func Test_AddLockCode(t *testing.T) {
	// We'll have a single device and managed lock code, and we'll add the lock code.

	ctx := context.Background()
	property := shared.Property{
		ControllerID: "9876",
		ID:           uuid.New(),
	}
	code := "5566"
	managedLockCode := &shared.DeviceManagedLockCode{
		Code:    code,
		EndAt:   time.Now().Add(1 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus1Scheduled,
		StartAt: time.Now().Add(-2 * time.Hour),
	}
	device := shared.Device{
		ID:               uuid.New(),
		ManagedLockCodes: []*shared.DeviceManagedLockCode{managedLockCode},
		PropertyID:       property.ID,
	}

	le, dc, dr, pr := newLockEngine(t)

	dc.EXPECT().AddLockCode(ctx, property, device, code).Return(nil)

	dr.EXPECT().List(ctx).Return(
		[]shared.Device{device},
		nil,
	)
	dr.EXPECT().Put(ctx, device).Return(shared.Device{}, nil)

	dr.EXPECT().AppendToAuditLog(ctx, device, []*shared.DeviceManagedLockCode{managedLockCode}).Return(nil)

	pr.EXPECT().GetCached(ctx, property.ID).Return(property, true, nil)

	err := le.UpdateLocks(ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus2Adding, managedLockCode.Status)
}

func Test_LeaveLockCode_MultipleMLC(t *testing.T) {
	// We'll have a single device and lock code, with multiple managed lock codes for the same code. One MLC will say to remove the code, the other will say to keep it.

	ctx := context.Background()
	property := shared.Property{
		ControllerID: "9876",
		ID:           uuid.New(),
	}
	code := "5566"
	activeManagedLockCode := &shared.DeviceManagedLockCode{
		Code:    code,
		EndAt:   time.Now().Add(1 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus3Enabled,
		StartAt: time.Now().Add(-2 * time.Hour),
	}
	expiredManagedLockCode := &shared.DeviceManagedLockCode{
		Code:    code,
		EndAt:   time.Now().Add(-3 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus4Removing,
		StartAt: time.Now().Add(-4 * time.Hour),
	}
	device := shared.Device{
		ID: uuid.New(),
		ManagedLockCodes: []*shared.DeviceManagedLockCode{
			activeManagedLockCode,
			expiredManagedLockCode,
		},
		PropertyID: property.ID,
		RawDevice: shared.RawDevice{
			LockCodes: []shared.RawDeviceLockCode{
				{
					Code: code,
				},
			},
		},
	}

	le, _, dr, pr := newLockEngine(t)

	dr.EXPECT().List(ctx).Return(
		[]shared.Device{device},
		nil,
	)

	dr.EXPECT().AppendToAuditLog(ctx, device, []*shared.DeviceManagedLockCode{expiredManagedLockCode}).Return(nil)
	dr.EXPECT().Put(ctx, device).Return(shared.Device{}, nil)

	pr.EXPECT().GetCached(ctx, property.ID).Return(property, true, nil)

	err := le.UpdateLocks(ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus3Enabled, activeManagedLockCode.Status)
	assert.Equal(t, shared.DeviceManagedLockCodeStatus5Complete, expiredManagedLockCode.Status)
}

func Test_LeaveLockCode_SingleMLC(t *testing.T) {
	// We'll have a single device, lock code, managed lock code, and we'll keep the lock code.

	ctx := context.Background()
	property := shared.Property{
		ControllerID: "9876",
		ID:           uuid.New(),
	}
	managedLockCode := &shared.DeviceManagedLockCode{
		Code:    "5566",
		EndAt:   time.Now().Add(1 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus3Enabled,
		StartAt: time.Now().Add(-2 * time.Hour),
	}
	device := shared.Device{
		ID:               uuid.New(),
		ManagedLockCodes: []*shared.DeviceManagedLockCode{managedLockCode},
		PropertyID:       property.ID,
		RawDevice: shared.RawDevice{
			LockCodes: []shared.RawDeviceLockCode{
				{
					Code: "5566",
				},
			},
		},
	}

	le, _, dr, pr := newLockEngine(t)

	dr.EXPECT().List(ctx).Return(
		[]shared.Device{device},
		nil,
	)

	pr.EXPECT().GetCached(ctx, property.ID).Return(property, true, nil)

	err := le.UpdateLocks(ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus3Enabled, managedLockCode.Status)
}

func Test_NoDevices(t *testing.T) {
	// With no devices, we shouldn't add or remove any lock codes.

	ctx := context.Background()
	le, _, dr, _ := newLockEngine(t)
	dr.EXPECT().List(ctx).Return([]shared.Device{}, nil)

	err := le.UpdateLocks(ctx)
	assert.Nil(t, err)

	// Since we didn't mock the DeviceController.* methods, this test will fail if any of them are called.
}

func Test_RemoveLockCode(t *testing.T) {
	// We'll have a single device, lock code, managed lock code, and remove the lock code.

	ctx := context.Background()
	property := shared.Property{
		ControllerID: "9876",
		ID:           uuid.New(),
	}
	code := "5566"
	managedLockCode := &shared.DeviceManagedLockCode{
		Code:    code,
		EndAt:   time.Now().Add(-1 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus3Enabled,
		StartAt: time.Now().Add(-2 * time.Hour),
	}
	device := shared.Device{
		ID:               uuid.New(),
		ManagedLockCodes: []*shared.DeviceManagedLockCode{managedLockCode},
		PropertyID:       property.ID,
		RawDevice: shared.RawDevice{
			LockCodes: []shared.RawDeviceLockCode{
				{
					Code: code,
				},
			},
		},
	}

	le, dc, dr, pr := newLockEngine(t)

	dc.EXPECT().RemoveLockCode(ctx, property, device, code).Return(nil)

	dr.EXPECT().List(ctx).Return(
		[]shared.Device{device},
		nil,
	)
	dr.EXPECT().Put(ctx, device).Return(shared.Device{}, nil)

	dr.EXPECT().AppendToAuditLog(ctx, device, []*shared.DeviceManagedLockCode{managedLockCode}).Return(nil)

	pr.EXPECT().GetCached(ctx, property.ID).Return(property, true, nil)

	err := le.UpdateLocks(ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus4Removing, managedLockCode.Status)
}

func Test_deleteOneButNotAllMLCs(t *testing.T) {
	// Because I like to mess up pointers in loops...

	ctx := context.Background()
	property := shared.Property{
		ControllerID: "9876",
		ID:           uuid.New(),
	}
	device := shared.Device{
		ID: uuid.New(),
		ManagedLockCodes: []*shared.DeviceManagedLockCode{
			// Old enough to delete, but since we keep the x most recent it'll hang out.
			{
				Code:    "0001",
				EndAt:   time.Now().Add(-1 * time.Hour * 24 * 8),
				Status:  shared.DeviceManagedLockCodeStatus5Complete,
				StartAt: time.Now().Add(-1 * time.Hour * 24 * 9),
			},
			// Old enough to delete, but since we keep the x most recent it'll hang out.
			{
				Code:    "0002",
				EndAt:   time.Now().Add(-1 * time.Hour * 24 * 9),
				Status:  shared.DeviceManagedLockCodeStatus5Complete,
				StartAt: time.Now().Add(-1 * time.Hour * 24 * 10),
			},
			// Old enough to delete, but since we keep the x most recent it'll hang out.
			{
				Code:    "0003",
				EndAt:   time.Now().Add(-1 * time.Hour * 24 * 10),
				Status:  shared.DeviceManagedLockCodeStatus5Complete,
				StartAt: time.Now().Add(-1 * time.Hour * 24 * 11),
			},
			// Old enough to delete, but since we keep the x most recent it'll hang out.
			{
				Code:    "0004",
				EndAt:   time.Now().Add(-1 * time.Hour * 24 * 11),
				Status:  shared.DeviceManagedLockCodeStatus5Complete,
				StartAt: time.Now().Add(-1 * time.Hour * 24 * 12),
			},
			// Expect to be deleted.
			{
				Code:    "0005",
				EndAt:   time.Now().Add(-1 * time.Hour * 24 * 12),
				Status:  shared.DeviceManagedLockCodeStatus5Complete,
				StartAt: time.Now().Add(-1 * time.Hour * 24 * 13),
			},
			// Expect to be deleted.
			{
				Code:    "0006",
				EndAt:   time.Now().Add(-1 * time.Hour * 24 * 13),
				Status:  shared.DeviceManagedLockCodeStatus5Complete,
				StartAt: time.Now().Add(-1 * time.Hour * 24 * 14),
			},
			// Don't delete because it's not a week old.
			{
				Code:    "9999",
				EndAt:   time.Now().Add(-1 * time.Hour * 24 * 1),
				Status:  shared.DeviceManagedLockCodeStatus5Complete,
				StartAt: time.Now().Add(-1 * time.Hour * 24 * 2),
			},
		},
		PropertyID: property.ID,
		RawDevice:  shared.RawDevice{},
	}

	le, _, dr, pr := newLockEngine(t)

	dr.EXPECT().List(ctx).Return(
		[]shared.Device{device},
		nil,
	)

	pr.EXPECT().GetCached(ctx, property.ID).Return(property, true, nil).AnyTimes()

	dr.EXPECT().AppendToAuditLog(
		ctx,
		gomock.Any(),
		gomock.Any(),
	).Do(func(ctx context.Context, d shared.Device, managedLockCodes []*shared.DeviceManagedLockCode) {
		assert.Equal(t, device.ID, d.ID)

		assert.Equal(t, 5, len(d.ManagedLockCodes))
		mlc := d.ManagedLockCodes[0]
		assert.Equal(t, "0001", mlc.Code)
		mlc = d.ManagedLockCodes[1]
		assert.Equal(t, "0002", mlc.Code)
		mlc = d.ManagedLockCodes[2]
		assert.Equal(t, "0003", mlc.Code)
		mlc = d.ManagedLockCodes[3]
		assert.Equal(t, "0004", mlc.Code)
		mlc = d.ManagedLockCodes[4]
		assert.Equal(t, "9999", mlc.Code)

		assert.Equal(t, 2, len(managedLockCodes))
		mlc = managedLockCodes[0]
		assert.Equal(t, "0005", mlc.Code)
		mlc = managedLockCodes[1]
		assert.Equal(t, "0006", mlc.Code)
	}).Return(nil)

	dr.EXPECT().Put(
		ctx,
		gomock.Any(),
	).Do(func(ctx context.Context, d shared.Device) {
		assert.Equal(t, d.ID, d.ID)
		assert.Equal(t, 5, len(d.ManagedLockCodes))
		mlc := d.ManagedLockCodes[0]
		assert.Equal(t, "0001", mlc.Code)
		mlc = d.ManagedLockCodes[1]
		assert.Equal(t, "0002", mlc.Code)
		mlc = d.ManagedLockCodes[2]
		assert.Equal(t, "0003", mlc.Code)
		mlc = d.ManagedLockCodes[3]
		assert.Equal(t, "0004", mlc.Code)
		mlc = d.ManagedLockCodes[4]
		assert.Equal(t, "9999", mlc.Code)
	}).Return(shared.Device{}, nil)

	err := le.UpdateLocks(ctx)
	assert.Nil(t, err)
}

func newLockEngine(t *testing.T) (*lockengine.LockEngine, *mock_lockengine.MockDeviceController, *mock_lockengine.MockDeviceRepository, *mock_lockengine.MockPropertyRepository) {
	ctrl := gomock.NewController(t)

	dc := mock_lockengine.NewMockDeviceController(ctrl)
	dr := mock_lockengine.NewMockDeviceRepository(ctrl)
	es := mock_lockengine.NewMockEmailService(ctrl)
	pr := mock_lockengine.NewMockPropertyRepository(ctrl)

	// This is probably temporary so we'll just completely mock it out all the way for now.
	es.EXPECT().SendEamil(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	le := lockengine.NewLockEngine(dc, dr, es, "", pr, time.UTC)
	return le, dc, dr, pr
}
