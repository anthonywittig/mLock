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

	if err := mshared.LoadConfig(); err != nil {
		return Response{}, fmt.Errorf("error loading config: %s", err.Error())
	}

	queueName, err := mshared.GetConfig("AWS_SQS_QUEUE_PREFIX")
	if err != nil {
		return Response{}, fmt.Errorf("error getting queue prefix: %s", err.Error())
	}

	queueName = queueName + "-in.fifo"

	log.Printf("adding message to \"%s\"\n", queueName)

	s, err := sqs.GetClient(ctx)
	if err != nil {
		return Response{}, fmt.Errorf("error getting sqs client: %s", err.Error())
	}

	if err := s.SendMessage(ctx, queueName, mshared.HabCommandListThings(fmt.Sprintf("hello there @ %s - requesting a list", time.Now().String()))); err != nil {
		return Response{}, fmt.Errorf("error sending message: %s", err.Error())
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
