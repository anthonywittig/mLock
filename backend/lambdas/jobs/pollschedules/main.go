package main

import (
	"context"
	"fmt"
	"log"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/device"
	"mlock/lambdas/shared/dynamo/property"
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
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	ctx = shared.CreateContextData(ctx)

	log.Printf("starting poll\n")

	if err := mshared.LoadConfig(); err != nil {
		return Response{}, fmt.Errorf("error loading config: %s", err.Error())
	}

	emailService, err := ses.NewEmailService(ctx)
	if err != nil {
		return Response{}, fmt.Errorf("error getting email service: %s", err.Error())
	}

	deviceRepository := device.NewRepository()
	reservationRepository := reservation.NewRepository()
	propertyRepository := property.NewRepository()
	unitRepository := unit.NewRepository()

	// Get the latest data from the controller and save it to the devices.
	ps, err := propertyRepository.List(ctx)
	if err != nil {
		return Response{}, fmt.Errorf("error getting properties: %s", err.Error())
	}
	for _, p := range ps {
		if err := updateDevicesFromController(ctx, emailService, p); err != nil {
			return Response{}, fmt.Errorf("error updating devices for property: %s, error: %s", p.Name, err.Error())
		}
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

	// Process and save any device changes to the controller.
	if err := lockengine.NewLockEngine(
		ezlo.NewLockCodeRepository(),
		deviceRepository,
		emailService,
		propertyRepository,
	).UpdateLocks(ctx); err != nil {
		return Response{}, fmt.Errorf("error updating lock codes: %s", err.Error())
	}

	return Response{
		Message: "ok",
	}, nil
}

func updateDevicesFromController(ctx context.Context, emailService *ses.EmailService, property shared.Property) error {
	rds, err := ezlo.GetDevices(ctx, property)
	if err != nil {
		return fmt.Errorf("error getting devices: %s", err.Error())
	}

	deviceRepository := device.NewRepository()

	eds, err := deviceRepository.List(ctx)
	if err != nil {
		return fmt.Errorf("error getting devices: %s", err.Error())
	}

	transitioningToOfflineDevices := []shared.Device{}
	offlineDevices := []shared.Device{}

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
			if ed.PropertyID == property.ID && ed.RawDevice.ID == rd.ID {
				// We found a match.
				d = ed

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
			}
		}

		d.PropertyID = property.ID
		d.RawDevice = rd
		d.LastRefreshedAt = time.Now()

		if _, err := deviceRepository.Put(ctx, d); err != nil {
			return fmt.Errorf("error putting device: %s", err.Error())
		}
	}

	if err := sendOfflineDeviceEmail(ctx, emailService, transitioningToOfflineDevices, offlineDevices); err != nil {
		return fmt.Errorf("error sending offline device email: %s", err.Error())
	}

	return nil
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

	if err := emailService.SendEamil(ctx, "MursetLock - Devices That Recently Went Offline", sb.String()); err != nil {
		return fmt.Errorf("error sending email: %s", err.Error())
	}

	return nil
}
