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

	// Should probably create a repository for this, but we're just listing them for now.
	controllers, err := ezlo.GetControllers(ctx)
	if err != nil {
		return Response{}, fmt.Errorf("error getting controllers: %s", err.Error())
	}

	// Old list, might be handy when transitioning to the API.
	// "84900483",
	// "84901957",
	// "84902073",
	// "84902125",
	// "84902278",
	// "84902829",
	// "84903155",
	// "84903375",
	// "84904595",
	// "84905147",
	// "84907280",
	// "84907380",
	// "84908450",
	// "84909917",
	// "84911082",
	// "90000849",
	// "90001483",
	// "90001613",
	// "90001793",
	// "90001871",
	// "90002061",
	// "90002779",
	// "90003204",
	// "90003500",
	// "90003526",
	// "90003983",
	// "90010778",
	// "90010799",
	// "90010815",
	// "90011629",
	// "90012570",
	// "92001809",
	for _, c := range controllers {
		ctxUpdateDevices, cancel := context.WithTimeout(ctx, 40*time.Second)
		defer cancel()

		if err := updateDevicesFromController(ctxUpdateDevices, emailService, c.PKDevice, deviceController); err != nil {
			if strings.Contains(err.Error(), "cloud.error.controller_not_connected") {
				// We get a ton of these when we're swapping out controllers. This should be temporary (but we know how that goes)...
				continue
			}
			if err2 := emailService.SendEamil(
				ctx,
				"MursetLock - Error updating devices from controller.",
				fmt.Sprintf("Controller ID: %s; error: %s", c.PKDevice, err.Error()),
			); err2 != nil {
				return Response{}, fmt.Errorf("error sending error email for updating devices for controller: %s, error: %s", c.PKDevice, err.Error())
			}
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

	return Response{
		Message: "ok",
	}, nil
}

func updateDevicesFromController(
	ctx context.Context,
	emailService *ses.EmailService,
	controllerID string,
	deviceController *ezlo.DeviceController,
) error {
	rds, err := deviceController.GetDevices(ctx, controllerID)
	if err != nil {
		return fmt.Errorf("error getting devices from controller: %s", err.Error())
	}

	deviceRepository := device.NewRepository()

	eds, err := deviceRepository.List(ctx)
	if err != nil {
		return fmt.Errorf("error getting devices from repository: %s", err.Error())
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
			if ed.ControllerID == controllerID && ed.RawDevice.ID == rd.ID {
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

		d.ControllerID = controllerID
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
