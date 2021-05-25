package main

import (
	"context"
	"fmt"
	"log"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/unit"
	"mlock/lambdas/shared/ical"
	"mlock/lambdas/shared/ses"
	"mlock/lambdas/shared/sqs"
	mshared "mlock/shared"
	"mlock/shared/protos/messaging"
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

	ctx = shared.CreateContextData(ctx)

	log.Printf("starting poll\n")

	queuePrefix := "test"
	log.Printf("adding message to \"%s\"\n", queuePrefix)
	if err := sqs.SendMessage(ctx, queuePrefix, &messaging.HabCommand{
		Description: fmt.Sprintf("hello there @ %s", time.Now().String()),
	}); err != nil {
		return Response{}, fmt.Errorf("error sending message: %s", err.Error())
	}

	if err := mshared.LoadConfig(); err != nil {
		return Response{}, fmt.Errorf("error loading config: %s", err.Error())
	}

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
		now := time.Now()
		notTooFarAway := now.Add(45 * time.Minute)
		if len(ress) != 0 {
			upcomingRess := []string{}
			for _, r := range ress {
				if r.Start.After(now) && r.Start.Before(notTooFarAway) {
					upcomingRess = append(upcomingRess, fmt.Sprintf("<li>tx:%s<ul>%s (start) (%f hours till)</ul><ul>%s (end) (%f hours till)</ul></li>", r.TransactionNumber, r.Start, r.Start.Sub(time.Now()).Hours(), r.End, r.End.Sub(time.Now()).Hours()))
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
