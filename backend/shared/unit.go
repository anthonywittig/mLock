package shared

type Unit struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	PropertyName string `json:"propertyName"`
	CalendarURL  string `json:"calendarUrl"`
	UpdatedBy    string `json:"updatedBy"`
}
