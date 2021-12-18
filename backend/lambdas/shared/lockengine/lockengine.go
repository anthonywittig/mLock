package lockengine

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"
	"time"

	"github.com/google/uuid"
)

type DeviceController interface {
	AddLockCode(ctx context.Context, prop shared.Property, device shared.Device, code string) error
	RemoveLockCode(ctx context.Context, prop shared.Property, device shared.Device, code string) error
}

type DeviceRepository interface {
	AppendToAuditLog(ctx context.Context, device shared.Device, managedLockCodes []*shared.DeviceManagedLockCode) error
	List(ctx context.Context) ([]shared.Device, error)
	Put(ctx context.Context, item shared.Device) (shared.Device, error)
}

type PropertyRepository interface {
	GetCached(ctx context.Context, id uuid.UUID) (shared.Property, bool, error)
}

type LockEngine struct {
	DeviceController   DeviceController
	DeviceRepository   DeviceRepository
	PropertyRepository PropertyRepository
}

type lockState struct {
	Exists          bool
	RequestToAdd    []*shared.DeviceManagedLockCode
	RequestToRemove []*shared.DeviceManagedLockCode
}

func NewLockEngine(dc DeviceController, dr DeviceRepository, pr PropertyRepository) *LockEngine {
	return &LockEngine{
		DeviceController:   dc,
		DeviceRepository:   dr,
		PropertyRepository: pr,
	}
}

func (l *LockEngine) UpdateLocks(ctx context.Context) error {
	ds, err := l.DeviceRepository.List(ctx)
	if err != nil {
		return fmt.Errorf("error getting devices: %s", err.Error())
	}

	now := time.Now()

	for _, d := range ds {
		needToSave := []*shared.DeviceManagedLockCode{}
		lockStates, err := l.getLockStates(now, d)
		if err != nil {
			return fmt.Errorf("error getting lock states for device: %s, error: %s", d.ID, err.Error())
		}

		for code, ls := range lockStates {
			prop, ok, err := l.PropertyRepository.GetCached(ctx, d.PropertyID)
			if err != nil {
				return fmt.Errorf("error getting property: %s", err.Error())
			}
			if !ok {
				return fmt.Errorf("error finding property: %s", d.PropertyID)
			}

			if ls.Exists {
				if len(ls.RequestToAdd) > 0 {
					for _, mlc := range ls.RequestToRemove {
						if mlc.Status == shared.DeviceManagedLockCodeStatus3Enabled || mlc.Status == shared.DeviceManagedLockCodeStatus4Removing {
							mlc.Status = shared.DeviceManagedLockCodeStatus5Complete
							mlc.Note = "Leaving lock code as it's in use."
							needToSave = append(needToSave, mlc)
						} else if mlc.Status != shared.DeviceManagedLockCodeStatus5Complete {
							return fmt.Errorf("unexpected status for exitsting remove with add, device: %s, id: %s, state: %s", d.ID, mlc.ID, mlc.Status)
						}
					}
					for _, mlc := range ls.RequestToAdd {
						if mlc.Status == shared.DeviceManagedLockCodeStatus1Scheduled || mlc.Status == shared.DeviceManagedLockCodeStatus2Adding {
							mlc.Status = shared.DeviceManagedLockCodeStatus3Enabled
							mlc.Note = "Lock code present."
							needToSave = append(needToSave, mlc)
						}
					}
				} else if len(ls.RequestToRemove) > 0 {
					if err := l.DeviceController.RemoveLockCode(ctx, prop, d, code); err != nil {
						return fmt.Errorf("error removing lock code: %s", err.Error())
					}

					for _, mlc := range ls.RequestToRemove {
						if mlc.Status == shared.DeviceManagedLockCodeStatus3Enabled || mlc.Status == shared.DeviceManagedLockCodeStatus5Complete {
							mlc.Status = shared.DeviceManagedLockCodeStatus4Removing
							mlc.Note = "Attempting to remove lock code."
							needToSave = append(needToSave, mlc)
						} else if mlc.Status == shared.DeviceManagedLockCodeStatus4Removing {
							// Assume this is a retry, do nothing.
						} else {
							return fmt.Errorf("unexpected status for existing remove without add, device: %s, id: %s, state: %s", d.ID, mlc.ID, mlc.Status)
						}
					}
				}
			} else { // !ls.Exists
				if len(ls.RequestToAdd) > 0 {
					if err := l.DeviceController.AddLockCode(ctx, prop, d, code); err != nil {
						return fmt.Errorf("error removing lock code: %s", err.Error())
					}

					for _, mlc := range ls.RequestToAdd {
						if mlc.Status == shared.DeviceManagedLockCodeStatus1Scheduled || mlc.Status == shared.DeviceManagedLockCodeStatus3Enabled {
							mlc.Status = shared.DeviceManagedLockCodeStatus2Adding
							mlc.Note = "Attempting to add lock code."
							needToSave = append(needToSave, mlc)
						} else if mlc.Status == shared.DeviceManagedLockCodeStatus2Adding {
							// Assume this is a retry, do nothing.
						} else {
							return fmt.Errorf("unexpected status for existing remove without add, device: %s, id: %s, state: %s", d.ID, mlc.ID, mlc.Status)
						}
					}

					if len(ls.RequestToRemove) > 1 {
						return fmt.Errorf("need to implement!")
					}
				} else {
					return fmt.Errorf("need to implement!")
				}

				/*
					for _, mlc := range ls.RequestToRemove {
						// tell them good job
					}

					if len(ls.RequestToAdd) > 0 {
						// add it

						for _, mlc := range ls.RequestToAdd {
						}
					}
				*/
			}
		}

		if len(needToSave) > 0 {
			if err := l.DeviceRepository.AppendToAuditLog(ctx, d, needToSave); err != nil {
				return fmt.Errorf("error appending to audit log: %s", err.Error())
			}

			if _, err := l.DeviceRepository.Put(ctx, d); err != nil {
				return fmt.Errorf("error putting device: %s", err.Error())
			}
		}
	}

	return nil
}

func (l *LockEngine) getLockStates(now time.Time, d shared.Device) (map[string]*lockState, error) {
	lockStates := map[string]*lockState{}

	for _, mlc := range d.ManagedLockCodes {
		if !mlc.HasStarted(now) {
			continue
		}

		ls, ok := lockStates[mlc.Code]
		if !ok {
			ls = &lockState{}
			lockStates[mlc.Code] = ls
			for _, lc := range d.RawDevice.LockCodes {
				if lc.Code == mlc.Code {
					ls.Exists = true
					break
				}
			}
		}

		if mlc.CodeShouldBePresent(now) {
			/*
				if mlc.Status != shared.DeviceManagedLockCodeStatus2Adding && mlc.Status != shared.DeviceManagedLockCodeStatus3Enabled {
					return nil, fmt.Errorf("mlc has unexpected status when requesting to add, id: %s, state: %s", mlc.ID, mlc.Status)
				}
			*/
			ls.RequestToAdd = append(ls.RequestToAdd, mlc)
		} else {
			/*
				if mlc.Status != shared.DeviceManagedLockCodeStatus4Removing && mlc.Status != shared.DeviceManagedLockCodeStatus5Complete {
					return nil, fmt.Errorf("mlc has unexpected status when requesting to remove, id: %s, state: %s", mlc.ID, mlc.Status)
				}
			*/
			ls.RequestToRemove = append(ls.RequestToRemove, mlc)
		}
	}

	return lockStates, nil
}
