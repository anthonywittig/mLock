package homeassistant_test

import (
	"context"
	"fmt"
	"mlock/lambdas/shared/homeassistant"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"

	"github.com/stretchr/testify/assert"
)

type authData struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type listing struct {
	ID                  int    `json:"id"`
	InternalListingName string `json:"internalListingName"`
}

type listingsPage struct {
	Status string    `json:"status"`
	Result []listing `json:"result"`
	Count  int       `json:"count"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}
type reservation struct {
	ArrivalDate           string `json:"arrivalDate"`
	ChannelID             int    `json:"channelId"`
	CheckInTime           int    `json:"checkInTime"`
	CheckOutTime          int    `json:"checkOutTime"`
	DepartureDate         string `json:"departureDate"`
	DoorCode              string `json:"doorCode"`
	HostawayReservationID string `json:"hostawayReservationId"`
	ListingMapID          int    `json:"listingMapId"`
	Status                string `json:"status"`
}

type reservationUpdateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type reservationsPage struct {
	Status string        `json:"status"`
	Result []reservation `json:"result"`
	Count  int           `json:"count"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}

func Test_aFewThings(t *testing.T) {
	assert.Nil(t, loadConfig())

	mux := http.NewServeMux()
	mux.HandleFunc("/api/states", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(getMockStatesData())
	})

	mux.HandleFunc("/api/states/", func(w http.ResponseWriter, r *http.Request) {
		pathSegments := strings.Split(r.URL.Path, "/")
		id := pathSegments[len(pathSegments)-1]

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(getMockEntityData(id))
	})

	// Setup a mock HTTP server with the custom mux
	server := httptest.NewServer(mux)
	defer server.Close()

	os.Setenv("HOME_ASSISTANT_BASE_URL", server.URL)

	r, err := homeassistant.NewRepository()
	assert.Nil(t, err)

	climateControls, err := r.ListClimateControls(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 20, len(climateControls))

	// Verify a few fields on one of the climate controls.
	climateControl5 := climateControls[4]
	/*
		{
			"entity_id": "climate.1e47e1fd",
			"state": "off",
			"attributes": {
				"hvac_modes": [
					"auto",
					"cool",
					"dry",
					"fan_only",
					"heat",
					"off"
				],
				"min_temp": 46,
				"max_temp": 86,
				"target_temp_step": 1,
				"fan_modes": [
					"auto",
					"low",
					"medium low",
					"medium",
					"medium high",
					"high"
				],
				"preset_modes": [
					"eco",
					"away",
					"boost",
					"none",
					"sleep"
				],
				"swing_modes": [
					"off",
					"vertical",
					"horizontal",
					"both"
				],
				"current_temperature": 54,
				"temperature": 70,
				"fan_mode": "medium",
				"preset_mode": "none",
				"swing_mode": "vertical",
				"friendly_name": "09A Gree",
				"supported_features": 57
			},
			"last_changed": "2024-01-26T20:40:23.045316+00:00",
			"last_updated": "2024-01-27T17:11:22.943704+00:00",
			"context": {
				"id": "01HN5YF3HZKMBQBN6RNS7A31RM",
				"parent_id": null,
				"user_id": null
			}
		}
	*/
	assert.Equal(t, "climate.1e47e1fd", climateControl5.EntityID)
	assert.Equal(t, "off", climateControl5.State)
	assert.Equal(t, 46, climateControl5.Attributes.MinTemp)
	assert.Equal(t, "09A Gree", climateControl5.Attributes.FriendlyName)
}

func loadConfig() error {
	if err := godotenv.Load(".env.test"); err != nil {
		return fmt.Errorf("error loading .env file: %s", err.Error())
	}

	return nil
}
