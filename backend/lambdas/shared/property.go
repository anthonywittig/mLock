package shared

import "github.com/google/uuid"

type Property struct {
	ControllerID string    `json:"controllerId"`
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	UpdatedBy    string    `json:"updatedBy"`
}
