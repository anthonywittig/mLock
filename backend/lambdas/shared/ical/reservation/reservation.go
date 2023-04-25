package reservation

import (
	"context"
	"fmt"
	"io/ioutil"
	"mlock/lambdas/shared"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type Repository struct {
	timeZone *time.Location
}

const (
	timeFormat = "20060102T150405Z"
)

func NewRepository(timeZone *time.Location) *Repository {
	return &Repository{
		timeZone: timeZone,
	}
}

func (r *Repository) Get(ctx context.Context, url string) ([]shared.Reservation, error) {
	data, err := getData(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("error getting data: %s", err.Error())
	}

	ress, err := parseReservations(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing reservations: %s", err.Error())
	}

	sort.Slice(ress, func(i, j int) bool {
		r1 := ress[i]
		r2 := ress[j]

		if r1.Start.Before(r2.Start) {
			return true
		}

		if r1.Start.Equal(r2.Start) {
			return r1.End.Before(r2.End)
		}

		return false
	})

	ress, err = r.updateTimestampsForCheckInOut(ress)
	if err != nil {
		return nil, fmt.Errorf("error updating reservation check in/out times: %s", err.Error())
	}

	return ress, nil
}

func (r *Repository) GetForUnits(ctx context.Context, units []shared.Unit) (map[uuid.UUID][]shared.Reservation, error) {
	reservations := make([][]shared.Reservation, len(units))

	g, ctx := errgroup.WithContext(ctx)
	for i, unit := range units {
		i, unit := i, unit // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			if unit.CalendarURL != "" {
				ress, err := r.Get(ctx, unit.CalendarURL)
				if err != nil {
					return fmt.Errorf("error getting calendar items: %s", err.Error())
				}
				reservations[i] = ress
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("error getting reservations: %s", err.Error())
	}

	// Might be able to do this inside the goroutine but I'm too lazy to figure out the possible errors with concurrency (or what other data structures to use).
	byUnitID := map[uuid.UUID][]shared.Reservation{}
	for i, unit := range units {
		byUnitID[unit.ID] = reservations[i]
	}

	return byUnitID, nil
}

func getData(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %s", err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error during call: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading body: %s", err.Error())
	}

	return string(body), nil
}

func parseReservations(data string) ([]shared.Reservation, error) {
	lines := strings.Split(data, "\r\n")

	if len(lines) == 0 {
		return nil, fmt.Errorf("no reservations to parse")
	}

	if len(lines) < 4 {
		return nil, fmt.Errorf("expect at least 4 lines to parse")
	}

	if _, err := getMatch("^BEGIN:VCALENDA(R)$", lines[0]); err != nil {
		return nil, fmt.Errorf("expected begin vcalendar: %s", lines[0])
	}

	if _, err := getMatch("^VERSION:(.*)$", lines[1]); err != nil {
		return nil, fmt.Errorf("expected version: %s", lines[1])
	}

	if _, err := getMatch("^PRODID:(.*)$", lines[2]); err != nil {
		return nil, fmt.Errorf("expected PRODID: %s", lines[2])
	}

	reservations := []shared.Reservation{}
	i := 3
	for {
		if lines[i] == "END:VCALENDAR" {
			// Could just return here, but we'll do a little extra validation.
			break
		}

		if _, err := getMatch("^BEGIN:(VEVENT)$", lines[i]); err != nil {
			return nil, fmt.Errorf("expected begin vevent: %s", lines[i])
		}
		i++

		res := shared.Reservation{}

		m, err := getMatch("^UID:(.*)$", lines[i])
		if err != nil {
			return nil, fmt.Errorf("expected UID: %s (%s)", lines[i], err.Error())
		}
		i++
		res.ID = m

		if _, err := getMatch("^DTSTAMP:(.*)$", lines[i]); err != nil {
			return nil, fmt.Errorf("expected DTSTAMP: %s", lines[i])
		}
		i++

		m, err = getMatch("^DTSTART:(.*)$", lines[i])
		if err != nil {
			return nil, fmt.Errorf("expected DTSTART: %s", lines[i])
		}
		i++
		start, err := time.Parse(timeFormat, m)
		res.Start = start

		m, err = getMatch("^DTEND:(.*)$", lines[i])
		if err != nil {
			return nil, fmt.Errorf("expected DTEND: %s", lines[i])
		}
		i++
		end, err := time.Parse(timeFormat, m)
		res.End = end

		m, err = getMatch("^SUMMARY:(.*)$", lines[i])
		if err != nil {
			return nil, fmt.Errorf("expected SUMMARY: %s", lines[i])
		}
		i++
		res.Summary = m

		if len(res.Summary) < 4 {
			return nil, fmt.Errorf("unexpectedly short summary (%s) for reservation ID: %s", res.Summary, res.ID)
		}
		res.TransactionNumber = res.Summary // The summary is the transaction number.

		if _, err := getMatch("^DESCRIPTION:(.*)$", lines[i]); err != nil {
			return nil, fmt.Errorf("expected DESCRIPTION: %s", lines[i])
		}
		i++

		m, err = getMatch("^STATUS:(.*)$", lines[i])
		if err != nil {
			return nil, fmt.Errorf("expected STATUS: %s", lines[i])
		}
		i++
		res.Status = m

		if _, err := getMatch("^END:(VEVENT)$", lines[i]); err != nil {
			return nil, fmt.Errorf("expected END:VEVENT: %s", lines[i])
		}
		i++

		reservations = append(reservations, res)
	}

	if lines[i] != "END:VCALENDAR" {
		return nil, fmt.Errorf("expected END:VCALENDAR: %s", lines[i])
	}
	i++

	if lines[i] != "" {
		return nil, fmt.Errorf("expected blank line: %s", lines[i])
	}
	i++

	if i != len(lines) {
		return nil, fmt.Errorf("line length doesn't _line_ up")
	}

	return reservations, nil
}

func getMatch(pattern string, input string) (string, error) {
	pat, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("error compiling regex: %s", err.Error())
	}

	matches := pat.FindAllStringSubmatch(input, -1)
	if len(matches) != 1 {
		return "", fmt.Errorf("error finding matches, expected length of 1 but was %d", len(matches))
	}

	submatch := matches[0]
	if len(submatch) != 2 {
		return "", fmt.Errorf("error expected submatch to have 2 items but was %d", len(submatch))
	}

	return submatch[1], nil
}

func (r *Repository) updateTimestampsForCheckInOut(ress []shared.Reservation) ([]shared.Reservation, error) {
	out := []shared.Reservation{}

	for _, res := range ress {
		t, err := r.getCheckinTimestamp(res.Start)
		if err != nil {
			return nil, fmt.Errorf("error getting check-in timestamp %s", err.Error())
		}
		res.Start = t

		t, err = r.getCheckoutTimestamp(res.End)
		if err != nil {
			return nil, fmt.Errorf("error getting check-out timestamp %s", err.Error())
		}
		res.End = t

		out = append(out, res)
	}

	return out, nil
}

func (r *Repository) getCheckinTimestamp(in time.Time) (time.Time, error) {
	utc := in.UTC()
	if utc.Hour() != 0 || utc.Minute() != 0 || utc.Second() != 0 || utc.Nanosecond() != 0 {
		return time.Time{}, fmt.Errorf("timestamp isn't midnight UTC: %s", in)
	}

	// Change to 4 pm (this should be configurable).
	fourPM := 16

	return time.Date(utc.Year(), utc.Month(), utc.Day(), fourPM, utc.Minute(), utc.Second(), utc.Nanosecond(), r.timeZone), nil
}

func (r *Repository) getCheckoutTimestamp(in time.Time) (time.Time, error) {
	utc := in.UTC()
	if utc.Hour() != 0 || utc.Minute() != 0 || utc.Second() != 0 || utc.Nanosecond() != 0 {
		return time.Time{}, fmt.Errorf("timestamp isn't midnight UTC: %s", in)
	}

	// Change to 11 am (this really should be configurable).
	elevenAM := 11

	return time.Date(utc.Year(), utc.Month(), utc.Day(), elevenAM, utc.Minute(), utc.Second(), utc.Nanosecond(), r.timeZone), nil

}
