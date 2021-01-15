package shared

import "time"

type Reservation struct {
	ID      string
	Start   time.Time
	End     time.Time
	Summary string
	Status  string
}
