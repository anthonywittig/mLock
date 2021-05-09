package shared

import "time"

type Reservation struct {
	ID                string    `json:"id"`
	Start             time.Time `json:"start"`
	End               time.Time `json:"end"`
	Summary           string    `json:"summary"`
	Status            string    `json:"status"`
	TransactionNumber string    `json:"transactionNumber"`
}
