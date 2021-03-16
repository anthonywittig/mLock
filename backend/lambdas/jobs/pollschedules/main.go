package main

import (
	"context"
	"fmt"
	"log"
	"mlock/shared"
	"mlock/shared/dynamo/unit"
	"mlock/shared/ical"
	"mlock/shared/ses"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
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

	log.Printf("starting poll\n")

	if err := shared.LoadConfig(); err != nil {
		return Response{}, fmt.Errorf("error loading config: %s", err.Error())
	}

	ctx = shared.CreateContextData(ctx)

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

	messages := []string{}
	for i, ress := range reservations {
		if len(ress) != 0 {
			u := units[i]
			messages = append(messages, fmt.Sprintf("unit: %s, some reservation: %+v", u.Name, ress[0]))
		}
	}

	// TODO: send email when reservations would start (and end?) - also log somewhere?
	if err := ses.SendEamil(ctx, "test reservation", strings.Join(messages, "; ")); err != nil {
		return Response{}, fmt.Errorf("error sending email: %s", err.Error())
	}

	return Response{
		Messages: messages,
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
