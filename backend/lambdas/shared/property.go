package shared

import "github.com/google/uuid"

type Property struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	UpdatedBy string    `json:"updatedBy"`
}
