package lockengine

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"
	"strings"
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

type EmailService interface {
	SendEamil(ctx context.Context, subject string, body string) error
}

type PropertyRepository interface {
	GetCached(ctx context.Context, id uuid.UUID) (shared.Property, bool, error)
}

type LockEngine struct {
	DeviceController   DeviceController
	DeviceRepository   DeviceRepository
	EmailService       EmailService
	PropertyRepository PropertyRepository
}

type lockState struct {
	Exists          bool
	RequestToAdd    []*shared.DeviceManagedLockCode
	RequestToRemove []*shared.DeviceManagedLockCode
}

func NewLockEngine(dc DeviceController, dr DeviceRepository, es EmailService, pr PropertyRepository) *LockEngine {
	return &LockEngine{
		DeviceController:   dc,
		DeviceRepository:   dr,
		EmailService:       es,
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
		lockStates := l.getLockStates(now, d)

		needToSave, err := l.calculateAndSendLockCommands(ctx, d, lockStates)
		if err != nil {
			return fmt.Errorf("error calculating and sending lock commands: %s", err.Error())
		}

		if len(needToSave) > 0 {
			if err := l.DeviceRepository.AppendToAuditLog(ctx, d, needToSave); err != nil {
				return fmt.Errorf("error appending to audit log: %s", err.Error())
			}

			if _, err := l.DeviceRepository.Put(ctx, d); err != nil {
				return fmt.Errorf("error putting device: %s", err.Error())
			}

			if err := l.sendEmailForAuditLogs(ctx, d, needToSave); err != nil {
				return fmt.Errorf("error sending email: %s", err.Error())
			}
		}
	}

	return nil
}

func (l *LockEngine) calculateAndSendLockCommands(ctx context.Context, device shared.Device, lockStates map[string]*lockState) ([]*shared.DeviceManagedLockCode, error) {
	needToSave := []*shared.DeviceManagedLockCode{}

	for code, ls := range lockStates {
		prop, ok, err := l.PropertyRepository.GetCached(ctx, device.PropertyID)
		if err != nil {
			return nil, fmt.Errorf("error getting property: %s", err.Error())
		}
		if !ok {
			return nil, fmt.Errorf("error finding property: %s", device.PropertyID)
		}

		if ls.Exists {
			if len(ls.RequestToAdd) > 0 {
				for _, mlc := range ls.RequestToRemove {
					if mlc.Status != shared.DeviceManagedLockCodeStatus5Complete {
						mlc.Status = shared.DeviceManagedLockCodeStatus5Complete
						mlc.Note = "Leaving lock code as it's in use."
						needToSave = append(needToSave, mlc)
					}
				}
				for _, mlc := range ls.RequestToAdd {
					if mlc.Status != shared.DeviceManagedLockCodeStatus3Enabled {
						mlc.Status = shared.DeviceManagedLockCodeStatus3Enabled
						mlc.Note = "Lock code present."
						needToSave = append(needToSave, mlc)
					}
				}
			} else if len(ls.RequestToRemove) > 0 {
				if err := l.DeviceController.RemoveLockCode(ctx, prop, device, code); err != nil {
					return nil, fmt.Errorf("error removing lock code: %s", err.Error())
				}

				for _, mlc := range ls.RequestToRemove {
					mlc.Status = shared.DeviceManagedLockCodeStatus4Removing
					mlc.Note = "Attempting to remove lock code."
					needToSave = append(needToSave, mlc)
				}
			}
		} else { // !ls.Exists
			if len(ls.RequestToAdd) > 0 {
				if err := l.DeviceController.AddLockCode(ctx, prop, device, code); err != nil {
					return nil, fmt.Errorf("error removing lock code: %s", err.Error())
				}

				for _, mlc := range ls.RequestToAdd {
					mlc.Status = shared.DeviceManagedLockCodeStatus2Adding
					mlc.Note = "Attempting to add lock code."
					needToSave = append(needToSave, mlc)
				}
			}

			for _, mlc := range ls.RequestToRemove {
				if mlc.Status != shared.DeviceManagedLockCodeStatus5Complete {
					mlc.Status = shared.DeviceManagedLockCodeStatus5Complete

					mlc.Note = "Code was removed."
					if len(ls.RequestToAdd) > 0 {
						mlc.Note = "Code is currently in use; nothing more to do."
					}

					needToSave = append(needToSave, mlc)
				}
			}
		}
	}

	return needToSave, nil
}

func (l *LockEngine) getLockStates(now time.Time, d shared.Device) map[string]*lockState {
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
			ls.RequestToAdd = append(ls.RequestToAdd, mlc)
		} else {
			ls.RequestToRemove = append(ls.RequestToRemove, mlc)
		}
	}

	return lockStates
}

func (l *LockEngine) sendEmailForAuditLogs(ctx context.Context, d shared.Device, needToSave []*shared.DeviceManagedLockCode) error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("<h1>New Audit Logs For Device: %s</h1>", d.RawDevice.Name))
	sb.WriteString("<ul>")
	for _, m := range needToSave {
		sb.WriteString(fmt.Sprintf("<li>Code: %s, Status: %s</li>", m.Code, m.Status))
	}
	sb.WriteString("</ul>")

	if err := l.EmailService.SendEamil(ctx, "MursetLock - Added Audit Log Entries", sb.String()); err != nil {
		return fmt.Errorf("error sending email: %s", err.Error())
	}

	return nil
}
