package ical

import (
	"context"
	"fmt"
	"io/ioutil"
	"mlock/shared"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	timeFormat = "20060102T150405Z"
)

func Get(ctx context.Context, url string) ([]shared.Reservation, error) {
	data, err := getData(ctx, url)
	if err != nil {
		return []shared.Reservation{}, fmt.Errorf("error getting data: %s", err.Error())
	}

	ress, err := parseReservations(data)
	if err != nil {
		return []shared.Reservation{}, fmt.Errorf("error parsing reservations: %s", err.Error())
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

	return ress, nil
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
		return []shared.Reservation{}, fmt.Errorf("no reservations to parse")
	}

	if len(lines) < 4 {
		return []shared.Reservation{}, fmt.Errorf("expect at least 4 lines to parse")
	}

	if _, err := getMatch("^BEGIN:VCALENDA(R)$", lines[0]); err != nil {
		return []shared.Reservation{}, fmt.Errorf("expected begin vcalendar: %s", lines[0])
	}

	if _, err := getMatch("^VERSION:(.*)$", lines[1]); err != nil {
		return []shared.Reservation{}, fmt.Errorf("expected version: %s", lines[1])
	}

	if _, err := getMatch("^PRODID:(.*)$", lines[2]); err != nil {
		return []shared.Reservation{}, fmt.Errorf("expected PRODID: %s", lines[2])
	}

	reservations := []shared.Reservation{}
	i := 3
	for {
		if lines[i] == "END:VCALENDAR" {
			// Could just return here, but we'll do a little extra validation.
			break
		}

		if _, err := getMatch("^BEGIN:(VEVENT)$", lines[i]); err != nil {
			return []shared.Reservation{}, fmt.Errorf("expected begin vevent: %s", lines[i])
		}
		i++

		res := shared.Reservation{}

		//m, err := getMatch(".*U.*I.*D.*:(.*)", lines[i])
		//m, err := getMatch("^UID:(.*)$", "UID:6799168@LiveRez.com")
		m, err := getMatch("^UID:(.*)$", lines[i])
		if err != nil {
			return []shared.Reservation{}, fmt.Errorf("expected UID: %s (%s)", lines[i], err.Error())
		}
		i++
		res.ID = m

		if _, err := getMatch("^DTSTAMP:(.*)$", lines[i]); err != nil {
			return []shared.Reservation{}, fmt.Errorf("expected DTSTAMP: %s", lines[i])
		}
		i++

		m, err = getMatch("^DTSTART:(.*)$", lines[i])
		if err != nil {
			return []shared.Reservation{}, fmt.Errorf("expected DTSTART: %s", lines[i])
		}
		i++
		start, err := time.Parse(timeFormat, m)
		res.Start = start

		m, err = getMatch("^DTEND:(.*)$", lines[i])
		if err != nil {
			return []shared.Reservation{}, fmt.Errorf("expected DTEND: %s", lines[i])
		}
		i++
		end, err := time.Parse(timeFormat, m)
		res.End = end

		m, err = getMatch("^SUMMARY:(.*)$", lines[i])
		if err != nil {
			return []shared.Reservation{}, fmt.Errorf("expected SUMMARY: %s", lines[i])
		}
		i++
		res.Summary = m
		res.TransactionNumber = res.Summary // The summary is the transaction number.

		if _, err := getMatch("^DESCRIPTION:(.*)$", lines[i]); err != nil {
			return []shared.Reservation{}, fmt.Errorf("expected DESCRIPTION: %s", lines[i])
		}
		i++

		m, err = getMatch("^STATUS:(.*)$", lines[i])
		if err != nil {
			return []shared.Reservation{}, fmt.Errorf("expected STATUS: %s", lines[i])
		}
		i++
		res.Status = m

		if _, err := getMatch("^END:(VEVENT)$", lines[i]); err != nil {
			return []shared.Reservation{}, fmt.Errorf("expected END:VEVENT: %s", lines[i])
		}
		i++

		reservations = append(reservations, res)
	}

	if lines[i] != "END:VCALENDAR" {
		return []shared.Reservation{}, fmt.Errorf("expected END:VCALENDAR: %s", lines[i])
	}
	i++

	if lines[i] != "" {
		return []shared.Reservation{}, fmt.Errorf("expected blank line: %s", lines[i])
	}
	i++

	if i != len(lines) {
		return []shared.Reservation{}, fmt.Errorf("line lenght doesn't _line_ up")
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
