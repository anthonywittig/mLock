package shared

import (
	"time"

	"github.com/google/uuid"
)

type ClimateControl struct {
	History           []ClimateControlHistory `json:"history"`
	ID                uuid.UUID               `json:"id"`
	LastRefreshedAt   time.Time               `json:"lastRefreshedAt"`
	RawClimateControl RawClimateControl       `json:"rawClimateControl"`
}

type ClimateControlHistory struct {
	Description       string            `json:"description"`
	RawClimateControl RawClimateControl `json:"rawClimateControl"`
	RecordedAt        time.Time         `json:"recordedAt"`
}

type RawClimateControl struct {
	EntityID   string `json:"entity_id"`
	State      string `json:"state"`
	Attributes struct {
		HVACModes          []string `json:"hvac_modes"`
		MinTemp            int      `json:"min_temp"`
		MaxTemp            int      `json:"max_temp"`
		TargetTempStep     int      `json:"target_temp_step"`
		PresetModes        []string `json:"preset_modes"`
		CurrentTemperature int      `json:"current_temperature"`
		Temperature        int      `json:"temperature"`
		HVACAction         string   `json:"hvac_action"`
		PresetMode         string   `json:"preset_mode"`
		FriendlyName       string   `json:"friendly_name"`
		SupportedFeatures  int      `json:"supported_features"`
	} `json:"attributes"`
	LastChanged string `json:"last_changed"`
	LastUpdated string `json:"last_updated"`
	Context     struct {
		ID string `json:"id"`
		// ParentID string `json:"parent_id"` // Not sure what the real type of this is.
		// UserID string `json:"user_id"`  // Not sure what the real type of this is.
	} `json:"context"`
}
