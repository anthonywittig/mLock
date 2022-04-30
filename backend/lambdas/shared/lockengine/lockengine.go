package lockengine

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"
	"strings"
	"time"
)

type DeviceController interface {
	AddLockCode(ctx context.Context, device shared.Device, code string) error
	RemoveLockCode(ctx context.Context, device shared.Device, code string) error
}

type DeviceRepository interface {
	AppendToAuditLog(ctx context.Context, device shared.Device, managedLockCodes []*shared.DeviceManagedLockCode) error
	ListActive(ctx context.Context) ([]shared.Device, error)
	Put(ctx context.Context, item shared.Device) (shared.Device, error)
}

type EmailService interface {
	SendEamil(ctx context.Context, subject string, body string) error
}

type LockEngine struct {
	deviceController DeviceController
	deviceRepository DeviceRepository
	emailService     EmailService
	frontEndDomain   string
	timeZone         *time.Location
}

type lockState struct {
	Exists          bool
	RequestToAdd    []*shared.DeviceManagedLockCode
	RequestToRemove []*shared.DeviceManagedLockCode
}

func NewLockEngine(
	dc DeviceController,
	dr DeviceRepository,
	es EmailService,
	fed string,
	tz *time.Location,
) *LockEngine {
	return &LockEngine{
		deviceController: dc,
		deviceRepository: dr,
		emailService:     es,
		frontEndDomain:   fed,
		timeZone:         tz,
	}
}

func (l *LockEngine) UpdateLocks(ctx context.Context) error {
	ds, err := l.deviceRepository.ListActive(ctx)
	if err != nil {
		return fmt.Errorf("error getting devices: %s", err.Error())
	}

	now := time.Now()
	nearPast := now.Add(-1 * time.Hour * 24 * 7)
	minPastCount := 5

	for _, d := range ds {
		lockStates := l.getLockStates(now, d)

		needToSave, err := l.calculateAndSendLockCommands(ctx, d, lockStates)
		if err != nil {
			return fmt.Errorf("error calculating and sending lock commands: %s", err.Error())
		}

		nonDeletedMLCs := []*shared.DeviceManagedLockCode{}
		completedMLCsCount := 0

		for i, mlc := range d.ManagedLockCodes {
			if mlc.EndAt.Before(nearPast) && mlc.Status == shared.DeviceManagedLockCodeStatus5Complete {
				justUpdated := false
				for _, m := range needToSave {
					if m == mlc {
						justUpdated = true
						break
					}
				}
				if !justUpdated {
					completedMLCsCount = completedMLCsCount + 1
					if completedMLCsCount >= minPastCount {
						d.ManagedLockCodes[i].Note = "Deleting code as it completed a while ago."
						needToSave = append(needToSave, d.ManagedLockCodes[i])
						continue
					}
				}
			}
			nonDeletedMLCs = append(nonDeletedMLCs, d.ManagedLockCodes[i])
		}
		d.ManagedLockCodes = nonDeletedMLCs

		if len(needToSave) > 0 {
			if err := l.deviceRepository.AppendToAuditLog(ctx, d, needToSave); err != nil {
				return fmt.Errorf("error appending to audit log: %s", err.Error())
			}

			if _, err := l.deviceRepository.Put(ctx, d); err != nil {
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

	enhancedLogging := device.RawDevice.Name == ""
	if enhancedLogging {
		fmt.Printf("vvv calculateAndSendLockCommands for %s vvv\n", device.RawDevice.Name)
	}

	for code, ls := range lockStates {
		if ls.Exists {
			if len(ls.RequestToAdd) > 0 {
				for _, mlc := range ls.RequestToRemove {
					if mlc.Status != shared.DeviceManagedLockCodeStatus5Complete {
						if err := mlc.SetStatus(shared.DeviceManagedLockCodeStatus5Complete); err != nil {
							return []*shared.DeviceManagedLockCode{}, err
						}
						mlc.Note = "Leaving lock code as it's in use."
						needToSave = append(needToSave, mlc)
					}
				}
				for _, mlc := range ls.RequestToAdd {
					if mlc.Status != shared.DeviceManagedLockCodeStatus3Enabled {
						if err := mlc.SetStatus(shared.DeviceManagedLockCodeStatus3Enabled); err != nil {
							return []*shared.DeviceManagedLockCode{}, err
						}
						mlc.Note = "Lock code present."
						needToSave = append(needToSave, mlc)
					}
				}
			} else if len(ls.RequestToRemove) > 0 {
				note := "Attempting to remove lock code."
				err := l.deviceController.RemoveLockCode(ctx, device, code)
				if err != nil {
					// TODO: log metric?
					fmt.Printf("error removing lock code: %s", err.Error())
					note = "Error attempting to remove lock code."
				}

				for _, mlc := range ls.RequestToRemove {
					if err := mlc.SetStatus(shared.DeviceManagedLockCodeStatus4Removing); err != nil {
						return []*shared.DeviceManagedLockCode{}, err
					}
					mlc.Note = note
					needToSave = append(needToSave, mlc)
				}
			}
		} else { // !ls.Exists
			if len(ls.RequestToAdd) > 0 {
				note := "Attempting to add lock code."
				err := l.deviceController.AddLockCode(ctx, device, code)
				if err != nil {
					// TODO: log metric?
					fmt.Printf("error adding lock code: %s", err.Error())
					note = "Error attempting to add lock code."
				}

				for _, mlc := range ls.RequestToAdd {
					if err := mlc.SetStatus(shared.DeviceManagedLockCodeStatus2Adding); err != nil {
						return []*shared.DeviceManagedLockCode{}, err
					}
					mlc.Note = note
					needToSave = append(needToSave, mlc)
				}
			}

			for _, mlc := range ls.RequestToRemove {
				if mlc.Status != shared.DeviceManagedLockCodeStatus5Complete {
					if err := mlc.SetStatus(shared.DeviceManagedLockCodeStatus5Complete); err != nil {
						return []*shared.DeviceManagedLockCode{}, err
					}
					mlc.Note = "Code was removed."
					if len(ls.RequestToAdd) > 0 {
						mlc.Note = "Code is currently in use; nothing more to do."
					}

					needToSave = append(needToSave, mlc)
				}
			}
		}
	}

	if enhancedLogging {
		for _, n := range needToSave {
			fmt.Printf("---- calculateAndSendLockCommands; need to save: %+v\n", n)
		}
		fmt.Printf("^^^ calculateAndSendLockCommands for %s - need to save: %+v ^^^\n", device.RawDevice.Name, needToSave)
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

	link := fmt.Sprintf("%s/devices/%s", l.frontEndDomain, d.ID)
	sb.WriteString(fmt.Sprintf("New audit logs for device: <a href=\"%s\">%s</a>", link, d.RawDevice.Name))
	sb.WriteString("<ul>")
	for _, m := range needToSave {
		sb.WriteString(fmt.Sprintf("<li>Code: %s, Status: %s</li>", m.Code, m.Status))
	}
	sb.WriteString("</ul>")

	now := time.Now().In(l.timeZone)
	startOfWeek := now.AddDate(0, 0, -1*int(now.Weekday()))
	weekOf := startOfWeek.Format("week of 01/02/2006")

	subject := fmt.Sprintf("MursetLock - Added Audit Log Entries - %s - %s", d.RawDevice.Name, weekOf)
	if err := l.emailService.SendEamil(ctx, subject, sb.String()); err != nil {
		return fmt.Errorf("error sending email: %s", err.Error())
	}

	return nil
}
