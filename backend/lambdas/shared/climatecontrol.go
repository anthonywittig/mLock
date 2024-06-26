package shared

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type ClimateControl struct {
	ActualState          ClimateControlActualState  `json:"actualState"`
	DesiredState         ClimateControlDesiredState `json:"desiredState"`
	History              []ClimateControlHistory    `json:"history"`
	ID                   uuid.UUID                  `json:"id"`
	LastRefreshedAt      time.Time                  `json:"lastRefreshedAt"`
	RawClimateControl    RawClimateControl          `json:"rawClimateControl"`
	SyncWithReservations bool                       `json:"syncWithReservations"`
}

type ClimateControlActualState struct {
	HVACMode    string `json:"hvacMode"`
	Temperature int    `json:"temperature"`
}

type ClimateControlDesiredState struct {
	AbandonAfter     time.Time  `json:"abandonAfter"`
	HVACMode         string     `json:"hvacMode"`
	Note             string     `json:"note"`
	SyncWithSettings bool       `json:"syncWithSettings"`
	Temperature      int        `json:"temperature"`
	WasSuccessfulAt  *time.Time `json:"wasSuccessfulAt"`
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
		Temperature        int      `json:"temperature"` // The target temperature.
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

func (c *ClimateControl) ActualStateMatchesDesiredState() bool {
	if c.ActualState.HVACMode == c.DesiredState.HVACMode && c.ActualState.HVACMode == "off" {
		// Some units don't let you update the temperature when the HVACMode is off.
		return true
	}
	return c.ActualState.HVACMode == c.DesiredState.HVACMode && c.ActualState.Temperature == c.DesiredState.Temperature
}

func (c *ClimateControl) GetFriendlyNamePrefix() string {
	return strings.Split(c.RawClimateControl.Attributes.FriendlyName, " ")[0]
}
