package main

import (
	"context"
	"fmt"
	"log"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/device"
	"mlock/lambdas/shared/dynamo/unit"
	"mlock/lambdas/shared/ezlo"
	"mlock/lambdas/shared/ical"
	"mlock/lambdas/shared/ses"
	mshared "mlock/shared"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type MyEvent struct {
}

type Response struct {
	Messages []string `json:"Messages"`
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

	// TODO: This should really move somewhere else.
	updateDevices(ctx)

	// get all units
	units, err := unit.List(ctx)
	if err != nil {
		return Response{}, fmt.Errorf("error getting entities: %s", err.Error())
	}

	// pull calendar info for those with a link
	reservations, err := getReservations(ctx, units)
	if err != nil {
		return Response{}, fmt.Errorf("error getting reservations: %s", err.Error())
	}

	var sb strings.Builder
	for i, ress := range reservations {
		// Start/End dates are UTC but they're really naive and just the day. Check-in is at 4 pm and check-out is at 11 am.
		now := time.Now()
		notTooFarAway := now.Add(5 * time.Minute)
		if len(ress) != 0 {
			upcomingRess := []string{}
			for _, r := range ress {
				if r.Start.After(now) && r.Start.Before(notTooFarAway) {
					upcomingRess = append(upcomingRess, fmt.Sprintf("<li>tx:%s<ul>%s (start) (%f hours till)</ul><ul>%s (end) (%f hours till)</ul></li>", r.TransactionNumber, r.Start, time.Until(r.Start).Hours(), r.End, time.Until(r.End).Hours()))
				}
			}

			if len(upcomingRess) == 0 {
				continue
			}

			u := units[i]
			sb.WriteString(fmt.Sprintf("<h1>Unit %s</h1>", u.Name))
			sb.WriteString("<ul>")
			for _, s := range upcomingRess {
				sb.WriteString(s)
			}
			sb.WriteString("</ul>")
		}
	}

	message := sb.String()

	if message == "" {
		return Response{
			Messages: []string{"no reservations to return"},
		}, nil
	}

	if err := ses.SendEamil(ctx, "Upcoming Reservations", message); err != nil {
		return Response{}, fmt.Errorf("error sending email: %s", err.Error())
	}

	return Response{
		Messages: []string{message},
	}, nil
}

func getReservations(ctx context.Context, units []shared.Unit) ([][]shared.Reservation, error) {
	reservations := make([][]shared.Reservation, len(units))

	g, ctx := errgroup.WithContext(ctx)
	for i, unit := range units {
		i, unit := i, unit // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			if unit.CalendarURL != "" {
				ress, err := ical.Get(ctx, unit.CalendarURL)
				if err != nil {
					return fmt.Errorf("error getting calendar items: %s", err.Error())
				}
				reservations[i] = ress
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return [][]shared.Reservation{}, fmt.Errorf("error getting reservations: %s", err.Error())
	}

	return reservations, nil
}

func updateDevices(ctx context.Context) error {
	ds, err := ezlo.GetDevices(ctx)
	if err != nil {
		return fmt.Errorf("error getting devices: %s", err.Error())
	}

	eds, err := device.List(ctx)
	if err != nil {
		return fmt.Errorf("error getting devices: %s", err.Error())
	}

	transitioningToOfflineDevices := []shared.Device{}
	offlineDevices := []shared.Device{}

	for _, d := range ds {
		device := shared.Device{
			History: []shared.DeviceHistory{
				{
					Description: "Initial State",
					EZLODevice:  d,
					RecordedAt:  time.Now(),
				},
			},
			ID: uuid.New(),
		}

		for _, ed := range eds {
			what do we do about properties?, v2?

			if ed.PropertyID == property.ID && ed.HABThing.UID == t.UID {
				// We found a match.
				d = ed

				wasOffline := t.StatusInfo.Status == shared.DeviceStatusOffline
				isOffline := d.HABThing.StatusInfo.Status == shared.DeviceStatusOffline
				if isOffline {
					offlineDevices = append(offlineDevices, d)
					if !wasOffline {
						now := time.Now()
						d.LastWentOfflineAt = &now
						transitioningToOfflineDevices = append(transitioningToOfflineDevices, d)
					}
				}

				statusChanged := (t.StatusInfo.Status != d.HABThing.StatusInfo.Status) || (t.StatusInfo.StatusDetail != d.HABThing.StatusInfo.StatusDetail)
				if statusChanged {
					d.History = append(d.History, shared.DeviceHistory{
						Description: "Status Changed",
						HABThing:    t,
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
		d.HABThing = t
		d.LastRefreshedAt = time.Now()

		if _, err := device.Put(ctx, d); err != nil {
			return fmt.Errorf("error putting device: %s", err.Error())
		}
	}

	if err := sendOfflineDeviceEmail(ctx, transitioningToOfflineDevices, offlineDevices); err != nil {
		return fmt.Errorf("error sending offline device email: %s", err.Error())
	}

	return nil
}
