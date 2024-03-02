package shared

import (
	"github.com/google/uuid"
)

type Miscellaneous struct {
	ID                             uuid.UUID              `json:"id"`
	ClimateControlOccupiedSettings ClimateControlSettings `json:"climateControlOccupiedSettings"`
	ClimateControlVacantSettings   ClimateControlSettings `json:"climateControlVacantSettings"`
}

type ClimateControlSettings struct {
	HVACMode    string `json:"hvacMode"`
	Temperature int    `json:"temperature"`
}
