package shared

import "time"

type Reservation struct {
	ID                string    `json:"id"`
	Start             time.Time `json:"start"`
	End               time.Time `json:"end"`
	DoorCode          string    `json:"doorCode"`
	TransactionNumber string    `json:"transactionNumber"`
}
