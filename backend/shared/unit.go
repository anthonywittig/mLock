package shared

import "github.com/google/uuid"

type UnitX struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	PropertyID  uuid.UUID `json:"propertyId"`
	CalendarURL string    `json:"calendarUrl"`
	UpdatedBy   string    `json:"updatedBy"`
}

type Unit2 struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	PropertyName string `json:"propertyName"`
	CalendarURL  string `json:"calendarUrl"`
	UpdatedBy    string `json:"updatedBy"`
}
