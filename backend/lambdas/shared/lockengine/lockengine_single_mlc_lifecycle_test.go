package lockengine_test

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

type singleMLCLifecycleTest struct {
	ctx context.Context
	d   shared.Device
	dc  *mock_lockengine.MockDeviceController
	dr  *mock_lockengine.MockDeviceRepository
	le  *lockengine.LockEngine
	mlc *shared.DeviceManagedLockCode
	p   shared.Property
	pr  *mock_lockengine.MockPropertyRepository
	t   *testing.T
}

func newSingleMLCLifecycleTest(t *testing.T) *singleMLCLifecycleTest {
	return &singleMLCLifecycleTest{
		ctx: context.Background(),
		t:   t,
	}
}

func Test_ScheduledBeforeNotExists(t *testing.T) {
	// MLC is scheduled but hasn't started yet.
	// Code does not exist.

	s := newSingleMLCLifecycleTest(t)

	s.mlc = &shared.DeviceManagedLockCode{
		Code:    "1234",
		EndAt:   time.Now().Add(4 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus1Scheduled,
		StartAt: time.Now().Add(3 * time.Hour),
	}

	s.generateLockengineSingleMLCLifecycleTest(false)

	err := s.le.UpdateLocks(s.ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus1Scheduled, s.mlc.Status)
}

func Test_ScheduledBeforeExists(t *testing.T) {
	// MLC is scheduled but hasn't started yet.
	// Code does exist, treat it as an unmanaged code and do nothing.

	s := newSingleMLCLifecycleTest(t)

	s.mlc = &shared.DeviceManagedLockCode{
		Code:    "1234",
		EndAt:   time.Now().Add(4 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus1Scheduled,
		StartAt: time.Now().Add(3 * time.Hour),
	}

	s.generateLockengineSingleMLCLifecycleTest(true)

	err := s.le.UpdateLocks(s.ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus1Scheduled, s.mlc.Status)
}

func Test_ScheduledStartedNotExists(t *testing.T) {
	// MLC is scheduled and has started.
	// Code does not exist.

	s := newSingleMLCLifecycleTest(t)

	s.mlc = &shared.DeviceManagedLockCode{
		Code:    "1234",
		EndAt:   time.Now().Add(4 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus1Scheduled,
		StartAt: time.Now().Add(-1 * time.Minute),
	}

	s.generateLockengineSingleMLCLifecycleTest(false)

	s.pr.EXPECT().GetCached(s.ctx, s.p.ID).Return(s.p, true, nil)
	s.dc.EXPECT().AddLockCode(s.ctx, s.p, s.d, s.mlc.Code).Return(nil)
	s.dr.EXPECT().AppendToAuditLog(s.ctx, s.d, []*shared.DeviceManagedLockCode{s.mlc}).Return(nil)
	s.dr.EXPECT().Put(s.ctx, s.d).Return(shared.Device{}, nil)

	err := s.le.UpdateLocks(s.ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus2Adding, s.mlc.Status)
}

func Test_ScheduledStartedExists(t *testing.T) {
	// MLC is scheduled and has started.
	// Code does exist.

	s := newSingleMLCLifecycleTest(t)

	s.mlc = &shared.DeviceManagedLockCode{
		Code:    "1234",
		EndAt:   time.Now().Add(4 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus1Scheduled,
		StartAt: time.Now().Add(-1 * time.Minute),
	}

	s.generateLockengineSingleMLCLifecycleTest(true)

	s.pr.EXPECT().GetCached(s.ctx, s.p.ID).Return(s.p, true, nil)
	s.dr.EXPECT().AppendToAuditLog(s.ctx, s.d, []*shared.DeviceManagedLockCode{s.mlc}).Return(nil)
	s.dr.EXPECT().Put(s.ctx, s.d).Return(shared.Device{}, nil)

	err := s.le.UpdateLocks(s.ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus3Enabled, s.mlc.Status)
}

func Test_AddingStartedNotExists(t *testing.T) {
	// MLC is adding, has started, but doesn't exist.

	s := newSingleMLCLifecycleTest(t)

	s.mlc = &shared.DeviceManagedLockCode{
		Code:    "1234",
		EndAt:   time.Now().Add(4 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus2Adding,
		StartAt: time.Now().Add(-1 * time.Minute),
	}

	s.generateLockengineSingleMLCLifecycleTest(false)

	s.pr.EXPECT().GetCached(s.ctx, s.p.ID).Return(s.p, true, nil)
	s.dc.EXPECT().AddLockCode(s.ctx, s.p, s.d, s.mlc.Code).Return(nil)
	s.dr.EXPECT().AppendToAuditLog(s.ctx, s.d, []*shared.DeviceManagedLockCode{s.mlc}).Return(nil)
	s.dr.EXPECT().Put(s.ctx, s.d).Return(shared.Device{}, nil)

	err := s.le.UpdateLocks(s.ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus2Adding, s.mlc.Status)
}

func Test_AddingStartedExists(t *testing.T) {
	// MLC is adding, has started, and does exist.

	s := newSingleMLCLifecycleTest(t)

	s.mlc = &shared.DeviceManagedLockCode{
		Code:    "1234",
		EndAt:   time.Now().Add(4 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus2Adding,
		StartAt: time.Now().Add(-1 * time.Minute),
	}

	s.generateLockengineSingleMLCLifecycleTest(true)

	s.pr.EXPECT().GetCached(s.ctx, s.p.ID).Return(s.p, true, nil)
	s.dr.EXPECT().AppendToAuditLog(s.ctx, s.d, []*shared.DeviceManagedLockCode{s.mlc}).Return(nil)
	s.dr.EXPECT().Put(s.ctx, s.d).Return(shared.Device{}, nil)

	err := s.le.UpdateLocks(s.ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus3Enabled, s.mlc.Status)
}

func Test_EnabledStartedExists(t *testing.T) {
	// MLC is enabled, has started, and does exist.

	s := newSingleMLCLifecycleTest(t)

	s.mlc = &shared.DeviceManagedLockCode{
		Code:    "1234",
		EndAt:   time.Now().Add(4 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus3Enabled,
		StartAt: time.Now().Add(-1 * time.Minute),
	}

	s.generateLockengineSingleMLCLifecycleTest(true)

	s.pr.EXPECT().GetCached(s.ctx, s.p.ID).Return(s.p, true, nil)

	err := s.le.UpdateLocks(s.ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus3Enabled, s.mlc.Status)
}

func Test_EnabledStartedNotExists(t *testing.T) {
	// MLC is enabled, has started, and does not exist.

	s := newSingleMLCLifecycleTest(t)

	s.mlc = &shared.DeviceManagedLockCode{
		Code:    "1234",
		EndAt:   time.Now().Add(4 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus3Enabled,
		StartAt: time.Now().Add(-1 * time.Minute),
	}

	s.generateLockengineSingleMLCLifecycleTest(false)

	s.pr.EXPECT().GetCached(s.ctx, s.p.ID).Return(s.p, true, nil)
	s.dc.EXPECT().AddLockCode(s.ctx, s.p, s.d, s.mlc.Code).Return(nil)
	s.dr.EXPECT().AppendToAuditLog(s.ctx, s.d, []*shared.DeviceManagedLockCode{s.mlc}).Return(nil)
	s.dr.EXPECT().Put(s.ctx, s.d).Return(shared.Device{}, nil)

	err := s.le.UpdateLocks(s.ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus2Adding, s.mlc.Status)
}

func Test_EnabledEndedExists(t *testing.T) {
	// MLC is enabled, has ended, and does exist.

	s := newSingleMLCLifecycleTest(t)

	s.mlc = &shared.DeviceManagedLockCode{
		Code:    "1234",
		EndAt:   time.Now().Add(-1 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus3Enabled,
		StartAt: time.Now().Add(-2 * time.Hour),
	}

	s.generateLockengineSingleMLCLifecycleTest(true)

	s.pr.EXPECT().GetCached(s.ctx, s.p.ID).Return(s.p, true, nil)
	s.dc.EXPECT().RemoveLockCode(s.ctx, s.p, s.d, s.mlc.Code).Return(nil)
	s.dr.EXPECT().AppendToAuditLog(s.ctx, s.d, []*shared.DeviceManagedLockCode{s.mlc}).Return(nil)
	s.dr.EXPECT().Put(s.ctx, s.d).Return(shared.Device{}, nil)

	err := s.le.UpdateLocks(s.ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus4Removing, s.mlc.Status)
}

func Test_RemovingEndedExists(t *testing.T) {
	// MLC is removing, has ended, and does exist.

	s := newSingleMLCLifecycleTest(t)

	s.mlc = &shared.DeviceManagedLockCode{
		Code:    "1234",
		EndAt:   time.Now().Add(-1 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus4Removing,
		StartAt: time.Now().Add(-2 * time.Hour),
	}

	s.generateLockengineSingleMLCLifecycleTest(true)

	s.pr.EXPECT().GetCached(s.ctx, s.p.ID).Return(s.p, true, nil)
	s.dc.EXPECT().RemoveLockCode(s.ctx, s.p, s.d, s.mlc.Code).Return(nil)
	s.dr.EXPECT().AppendToAuditLog(s.ctx, s.d, []*shared.DeviceManagedLockCode{s.mlc}).Return(nil)
	s.dr.EXPECT().Put(s.ctx, s.d).Return(shared.Device{}, nil)

	err := s.le.UpdateLocks(s.ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus4Removing, s.mlc.Status)
}

func Test_RemovingEndedNotExists(t *testing.T) {
	// MLC is removing, has ended, and does not exist.

	s := newSingleMLCLifecycleTest(t)

	s.mlc = &shared.DeviceManagedLockCode{
		Code:    "1234",
		EndAt:   time.Now().Add(-1 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus4Removing,
		StartAt: time.Now().Add(-2 * time.Hour),
	}

	s.generateLockengineSingleMLCLifecycleTest(false)

	s.pr.EXPECT().GetCached(s.ctx, s.p.ID).Return(s.p, true, nil)
	s.dr.EXPECT().AppendToAuditLog(s.ctx, s.d, []*shared.DeviceManagedLockCode{s.mlc}).Return(nil)
	s.dr.EXPECT().Put(s.ctx, s.d).Return(shared.Device{}, nil)

	err := s.le.UpdateLocks(s.ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus5Complete, s.mlc.Status)
}

func Test_CompletedEndedNotExists(t *testing.T) {
	// MLC is completed, has ended, and does not exist.

	s := newSingleMLCLifecycleTest(t)

	s.mlc = &shared.DeviceManagedLockCode{
		Code:    "1234",
		EndAt:   time.Now().Add(-1 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus5Complete,
		StartAt: time.Now().Add(-2 * time.Hour),
	}

	s.generateLockengineSingleMLCLifecycleTest(false)

	s.pr.EXPECT().GetCached(s.ctx, s.p.ID).Return(s.p, true, nil)

	err := s.le.UpdateLocks(s.ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus5Complete, s.mlc.Status)
}

func Test_CompletedEndedLongAgoNotExists(t *testing.T) {
	// MLC is completed, has ended (long ago), and does not exist.

	s := newSingleMLCLifecycleTest(t)

	s.mlc = &shared.DeviceManagedLockCode{
		Code:    "1234",
		EndAt:   time.Now().Add(-1 * time.Hour * 24 * 8),
		Status:  shared.DeviceManagedLockCodeStatus5Complete,
		StartAt: time.Now().Add(-1 * time.Hour * 24 * 9),
	}

	s.generateLockengineSingleMLCLifecycleTest(false)

	s.pr.EXPECT().GetCached(s.ctx, s.p.ID).Return(s.p, true, nil)

	s.dr.EXPECT().AppendToAuditLog(
		s.ctx,
		gomock.Any(),
		gomock.Any(),
	).Do(func(ctx context.Context, d shared.Device, managedLockCodes []*shared.DeviceManagedLockCode) {
		assert.Equal(t, s.d.ID, d.ID)
		assert.Equal(t, 0, len(d.ManagedLockCodes))

		assert.Equal(t, 1, len(managedLockCodes))
		mlc := managedLockCodes[0]
		assert.Equal(t, "1234", mlc.Code)
		assert.Equal(t, shared.DeviceManagedLockCodeStatus5Complete, mlc.Status)
	}).Return(nil)

	s.dr.EXPECT().Put(
		s.ctx,
		gomock.Any(),
	).Do(func(ctx context.Context, d shared.Device) {
		assert.Equal(t, s.d.ID, d.ID)
		assert.Equal(t, 0, len(d.ManagedLockCodes))
	}).Return(shared.Device{}, nil)

	err := s.le.UpdateLocks(s.ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus5Complete, s.mlc.Status)
}

func Test_CompletedEndedExists(t *testing.T) {
	// MLC is completed, has ended, and does exist.

	s := newSingleMLCLifecycleTest(t)

	s.mlc = &shared.DeviceManagedLockCode{
		Code:    "1234",
		EndAt:   time.Now().Add(-1 * time.Hour),
		Status:  shared.DeviceManagedLockCodeStatus5Complete,
		StartAt: time.Now().Add(-2 * time.Hour),
	}

	s.generateLockengineSingleMLCLifecycleTest(true)

	s.pr.EXPECT().GetCached(s.ctx, s.p.ID).Return(s.p, true, nil)
	s.dc.EXPECT().RemoveLockCode(s.ctx, s.p, s.d, s.mlc.Code).Return(nil)
	s.dr.EXPECT().AppendToAuditLog(s.ctx, s.d, []*shared.DeviceManagedLockCode{s.mlc}).Return(nil)
	s.dr.EXPECT().Put(s.ctx, s.d).Return(shared.Device{}, nil)

	err := s.le.UpdateLocks(s.ctx)
	assert.Nil(t, err)

	assert.Equal(t, shared.DeviceManagedLockCodeStatus4Removing, s.mlc.Status)
}

func (s *singleMLCLifecycleTest) generateLockengineSingleMLCLifecycleTest(codeExists bool) {
	s.p = shared.Property{
		ControllerID: "9876",
		ID:           uuid.New(),
	}

	s.d = shared.Device{
		ID:               uuid.New(),
		ManagedLockCodes: []*shared.DeviceManagedLockCode{s.mlc},
		PropertyID:       s.p.ID,
	}

	if codeExists {
		s.d.RawDevice.LockCodes = []shared.RawDeviceLockCode{
			{
				Code: s.mlc.Code,
			},
		}
	}

	s.le, s.dc, s.dr, s.pr = newLockEngine(s.t)

	s.dr.EXPECT().List(s.ctx).Return(
		[]shared.Device{s.d},
		nil,
	)
}
