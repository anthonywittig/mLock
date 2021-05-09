package shared

import "github.com/google/uuid"

type Unit struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	PropertyID  uuid.UUID `json:"propertyId"`
	CalendarURL string    `json:"calendarUrl"`
	UpdatedBy   string    `json:"updatedBy"`
}
