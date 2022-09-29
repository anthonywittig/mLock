package main

import (
	"context"
	"fmt"
	"log"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/device"
	"mlock/lambdas/shared/dynamo/unit"
	"mlock/lambdas/shared/ezlo"
	"mlock/lambdas/shared/ical/reservation"
	"mlock/lambdas/shared/lockengine"
	"mlock/lambdas/shared/scheduler"
	"mlock/lambdas/shared/ses"
	mshared "mlock/shared"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

type MyEvent struct {
}

type Response struct {
	Message string `json:"message"`
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, event MyEvent) (Response, error) {
	ctx = shared.CreateContextData(ctx)

	log.Printf("starting poll\n")

	if err := mshared.LoadConfig(); err != nil {
		return Response{}, fmt.Errorf("error loading config: %s", err.Error())
	}

	emailService, err := ses.NewEmailService(ctx)
	if err != nil {
		return Response{}, fmt.Errorf("error getting email service: %s", err.Error())
	}

	tzName, err := mshared.GetConfig("TIME_ZONE")
	if err != nil {
		return Response{}, fmt.Errorf("error getting time zone name: %s", err.Error())
	}

	tz, err := time.LoadLocation(tzName)
	if err != nil {
		return Response{}, fmt.Errorf("error getting time zone %s", err.Error())
	}

	connectionPool := ezlo.NewConnectionPool()
	defer connectionPool.Close()

	deviceController := ezlo.NewDeviceController(connectionPool)
	deviceRepository := device.NewRepository()
	reservationRepository := reservation.NewRepository(tz)
	unitRepository := unit.NewRepository()

	if err := updateDevicesFromController(
		ctx,
		emailService,
		deviceController,
		deviceRepository,
	); err != nil {
		return Response{}, fmt.Errorf("error updating devices from controller: %s", err.Error())
	}

	// Get the latest data from the reservations and save it to the devices.
	if err := scheduler.NewScheduler(
		deviceRepository,
		time.Now(),
		reservationRepository,
		unitRepository,
	).ReconcileReservationsAndLockCodes(ctx); err != nil {
		return Response{}, fmt.Errorf("error scheduling: %s", err.Error())
	}

	fed, err := mshared.GetConfig("FRONTEND_DOMAIN")
	if err != nil {
		return Response{}, fmt.Errorf("error getting front end domain: %s", err.Error())
	}

	// Process and save any device changes to the controller.
	if err := lockengine.NewLockEngine(
		deviceController,
		deviceRepository,
		emailService,
		fed,
		tz,
	).UpdateLocks(ctx); err != nil {
		return Response{}, fmt.Errorf("error updating lock codes: %s", err.Error())
	}

	/*
		This caused a few devices to get lost from their controllers, so we're disabling it for now.

		// Rediscover devices if needed.
		if err := rediscoverUnresponsiveDevices(
			ctx,
			deviceController,
			deviceRepository,
			emailService,
		); err != nil {
			return Response{}, fmt.Errorf("error rediscovering unresponsive devices: %s", err.Error())
		}
	*/

	return Response{
		Message: "ok",
	}, nil
}

func rediscoverUnresponsiveDevices(
	ctx context.Context,
	deviceController *ezlo.DeviceController,
	deviceRepository *device.Repository,
	emailService *ses.EmailService,
) error {
	devices, err := deviceRepository.List(ctx)
	if err != nil {
		return fmt.Errorf("error listing devices: %s", err.Error())
	}

	for _, d := range devices {
		if d.RawDevice.Status != shared.DeviceStatusOnline {
			continue
		}
		if d.LastRediscoveredAt != nil && time.Since(*d.LastRediscoveredAt) < 4*time.Hour {
			continue
		}

		var shouldRediscoverFor *shared.DeviceManagedLockCode = nil
		for _, mlc := range d.ManagedLockCodes {
			if mlc.Status != shared.DeviceManagedLockCodeStatus2Adding {
				continue
			}
			if time.Since(*mlc.StartedAddingAt) < 5*time.Minute {
				continue
			}
			shouldRediscoverFor = mlc
		}
		if shouldRediscoverFor == nil {
			continue
		}

		fmt.Printf("Rediscovering device %s\n", d.RawDevice.Name)
		if err := deviceController.RediscoverDevice(ctx, d); err != nil {
			return fmt.Errorf("error rediscovering device %s: %s", d.RawDevice.Name, err.Error())
		}

		now := time.Now()
		d.LastRediscoveredAt = &now
		d, err := deviceRepository.Put(ctx, d)
		if err != nil {
			return fmt.Errorf("error saving device %s: %s", d.RawDevice.Name, err.Error())
		}

		shouldRediscoverFor.Note = "Rediscovering device."
		if err := deviceRepository.AppendToAuditLog(ctx, d, []*shared.DeviceManagedLockCode{shouldRediscoverFor}); err != nil {
			return fmt.Errorf("error appending to audit log: %s", err.Error())
		}

		emailService.SendEmailToDevelopers(
			ctx,
			"Rediscovering device",
			fmt.Sprintf("Rediscovering device %s", d.RawDevice.Name),
		)
	}

	return nil
}

func updateDevicesFromController(
	ctx context.Context,
	emailService *ses.EmailService,
	deviceController *ezlo.DeviceController,
	deviceRepository *device.Repository,
) error {
	devices, err := deviceRepository.List(ctx)
	if err != nil {
		return fmt.Errorf("error getting devices from repository: %s", err.Error())
	}

	transitioningToOfflineDevices := []shared.Device{}
	offlineDevices := []shared.Device{}

	online, offline, err := ezlo.GetControllers(ctx)
	if err != nil {
		return fmt.Errorf("error getting controllers: %s", err.Error())
	}

	for _, c := range online {
		ctxUpdateDevices, cancel := context.WithTimeout(ctx, 40*time.Second)
		defer cancel()

		tTODevices, oDevices, err := updateOnlineDevicesFromController(
			ctxUpdateDevices,
			emailService,
			c.PKDevice,
			deviceController,
			deviceRepository,
			devices,
		)
		if err != nil {
			if strings.Contains(err.Error(), "cloud.error.controller_not_connected") {
				// We get a ton of these when we're swapping out controllers. This should be temporary (but we know how that goes)...
				continue
			}
			if err2 := emailService.SendEmailToDevelopers(
				ctx,
				"zcclock - Error updating devices from controller.",
				fmt.Sprintf("Controller ID: %s; error: %s", c.PKDevice, err.Error()),
			); err2 != nil {
				return fmt.Errorf("error sending error email for updating devices for controller: %s, error: %s", c.PKDevice, err.Error())
			}
		}
		transitioningToOfflineDevices = append(transitioningToOfflineDevices, tTODevices...)
		offlineDevices = append(offlineDevices, oDevices...)
	}

	for _, c := range offline {
		ctxUpdateDevices, cancel := context.WithTimeout(ctx, 40*time.Second)
		defer cancel()

		tTODevices, oDevices, err := updateOfflineDevicesFromController(
			ctxUpdateDevices,
			emailService,
			c.PKDevice,
			deviceRepository,
			devices,
		)
		if err != nil {
			if err2 := emailService.SendEmailToDevelopers(
				ctx,
				"zcclock - Error updating devices from controller.",
				fmt.Sprintf("Controller ID: %s; error: %s", c.PKDevice, err.Error()),
			); err2 != nil {
				return fmt.Errorf("error sending error email for updating devices for controller: %s, error: %s", c.PKDevice, err.Error())
			}
		}
		transitioningToOfflineDevices = append(transitioningToOfflineDevices, tTODevices...)
		offlineDevices = append(offlineDevices, oDevices...)
	}

	if err := sendOfflineDeviceEmail(ctx, emailService, transitioningToOfflineDevices, offlineDevices); err != nil {
		return fmt.Errorf("error sending offline device email: %s", err.Error())
	}

	return nil
}

func updateOfflineDevicesFromController(
	ctx context.Context,
	emailService *ses.EmailService,
	controllerID string,
	deviceRepository *device.Repository,
	devices []shared.Device,
) (
	[]shared.Device,
	[]shared.Device,
	error,
) {
	transitioningToOfflineDevices := []shared.Device{}
	offlineDevices := []shared.Device{}

	for _, ed := range devices {
		if ed.ControllerID != controllerID {
			continue
		}

		wasOffline := ed.RawDevice.Status == shared.DeviceStatusOffline

		offlineDevices = append(offlineDevices, ed)
		if !wasOffline {
			now := time.Now()
			ed.LastWentOfflineAt = &now
			transitioningToOfflineDevices = append(transitioningToOfflineDevices, ed)
		}

		ed.RawDevice.Status = shared.DeviceStatusOffline

		maxHistoryCount := 1
		historyStartIndex := len(ed.History) - maxHistoryCount
		if historyStartIndex > 0 {
			ed.History = ed.History[historyStartIndex:]
		}
		if ed.RawDevice.Status != shared.DeviceStatusOffline {
			ed.History = append(ed.History, shared.DeviceHistory{
				Description: "Status Changed",
				RawDevice:   ed.RawDevice,
				RecordedAt:  time.Now(),
			})
		}

		if _, err := deviceRepository.Put(ctx, ed); err != nil {
			return transitioningToOfflineDevices, offlineDevices, fmt.Errorf("error putting device: %s", err.Error())
		}
	}

	return transitioningToOfflineDevices, offlineDevices, nil
}

func updateOnlineDevicesFromController(
	ctx context.Context,
	emailService *ses.EmailService,
	controllerID string,
	deviceController *ezlo.DeviceController,
	deviceRepository *device.Repository,
	eds []shared.Device,
) (
	[]shared.Device,
	[]shared.Device,
	error,
) {
	transitioningToOfflineDevices := []shared.Device{}
	offlineDevices := []shared.Device{}

	rds, err := deviceController.GetDevices(ctx, controllerID)
	if err != nil {
		return transitioningToOfflineDevices, offlineDevices, fmt.Errorf("error getting devices from controller: %s", err.Error())
	}

	for _, rd := range rds {
		d := shared.Device{
			History: []shared.DeviceHistory{
				{
					Description: "Initial State",
					RecordedAt:  time.Now(),
					RawDevice:   rd,
				},
			},
			ID: uuid.New(),
		}

		for _, ed := range eds {
			if ed.ControllerID == controllerID && ed.RawDevice.ID == rd.ID {
				// We found a match.
				var tTOD []shared.Device
				var oDs []shared.Device
				d, tTOD, oDs = updateDeviceWithRawData(ed, rd)
				transitioningToOfflineDevices = append(transitioningToOfflineDevices, tTOD...)
				offlineDevices = append(offlineDevices, oDs...)
			}
		}

		d.ControllerID = controllerID
		d.RawDevice = rd
		d.LastRefreshedAt = time.Now()

		if _, err := deviceRepository.Put(ctx, d); err != nil {
			return transitioningToOfflineDevices, offlineDevices, fmt.Errorf("error putting device: %s", err.Error())
		}
	}

	// Look for devices that no longer exist on the controller.
	for _, ed := range eds {
		if ed.ControllerID != controllerID {
			continue
		}
		found := false
		for _, rd := range rds {
			if ed.RawDevice.ID == rd.ID {
				found = true
				break
			}
		}
		if found {
			continue
		}
		if ed.RawDevice.Status == shared.DeviceStatusOffline {
			continue
		}

		rd := ed.RawDevice
		rd.Status = shared.DeviceStatusOffline // Fake an offline status.
		d, tTOD, oDs := updateDeviceWithRawData(ed, rd)
		d.RawDevice = rd
		transitioningToOfflineDevices = append(transitioningToOfflineDevices, tTOD...)
		offlineDevices = append(offlineDevices, oDs...)

		if _, err := deviceRepository.Put(ctx, d); err != nil {
			return transitioningToOfflineDevices, offlineDevices, fmt.Errorf("error putting device: %s", err.Error())
		}
	}

	return transitioningToOfflineDevices, offlineDevices, nil
}

func updateDeviceWithRawData(d shared.Device, rd shared.RawDevice) (
	shared.Device,
	[]shared.Device,
	[]shared.Device,
) {
	transitioningToOfflineDevices := []shared.Device{}
	offlineDevices := []shared.Device{}

	wasOffline := d.RawDevice.Status == shared.DeviceStatusOffline
	isOffline := rd.Status == shared.DeviceStatusOffline

	if wasOffline && !isOffline {
		now := time.Now()
		d.LastWentOnlineAt = &now
	}

	if isOffline {
		offlineDevices = append(offlineDevices, d)
		if !wasOffline {
			now := time.Now()
			d.LastWentOfflineAt = &now
			transitioningToOfflineDevices = append(transitioningToOfflineDevices, d)
		}
	}

	statusChanged := d.RawDevice.Status != rd.Status
	if statusChanged {
		d.History = append(d.History, shared.DeviceHistory{
			Description: "Status Changed",
			RawDevice:   rd,
			RecordedAt:  time.Now(),
		})
	}

	maxHistoryCount := 1
	historyStartIndex := len(d.History) - maxHistoryCount
	if historyStartIndex > 0 {
		d.History = d.History[historyStartIndex:]
	}

	return d, transitioningToOfflineDevices, offlineDevices
}

func sendOfflineDeviceEmail(ctx context.Context, emailService *ses.EmailService, transitioningToOfflineDevices []shared.Device, offlineDevices []shared.Device) error {
	if len(transitioningToOfflineDevices) == 0 {
		return nil
	}

	var sb strings.Builder

	sb.WriteString("<h1>Devices That Recently Went Offline</h1>")
	sb.WriteString("<ul>")
	for _, d := range transitioningToOfflineDevices {
		sb.WriteString(fmt.Sprintf("<li>Device: %s</li>", d.RawDevice.Name))
	}
	sb.WriteString("</ul>")

	sb.WriteString("<h1>Devices That Are Currently Offline</h1>")
	sb.WriteString("<ul>")
	for _, d := range offlineDevices {
		sb.WriteString(fmt.Sprintf("<li>Device: %s</li>", d.RawDevice.Name))
	}
	sb.WriteString("</ul>")

	if err := emailService.SendEmailToAdmins(ctx, "zcclock - Devices That Recently Went Offline", sb.String()); err != nil {
		return fmt.Errorf("error sending email: %s", err.Error())
	}

	return nil
}
