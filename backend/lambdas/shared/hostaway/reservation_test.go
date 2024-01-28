package hostaway_test

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/hostaway"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/joho/godotenv"

	"github.com/google/uuid"
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
	mux.HandleFunc("/v1/accessTokens", func(w http.ResponseWriter, r *http.Request) {
		mockJSON, _ := json.Marshal(authData{
			AccessToken: "accessTokenValue",
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(mockJSON)
	})
	/*
		mux.HandleFunc("/v1/listings", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			query := r.URL.Query()
			offset := query.Get("offset")
			if offset == "0" {
				mockJSON, _ := json.Marshal(listingsPage{
					Result: []listing{
						{
							ID:                  25,
							InternalListingName: "01A Bunkhouse",
						},
					},
					Status: "success",
				})
				w.Write(mockJSON)
			} else {
				mockJSON, _ := json.Marshal(listingsPage{
					Result: []listing{},
					Status: "success",
				})
				w.Write(mockJSON)

			}
		})
	*/
	mux.HandleFunc("/v1/listings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			mockJSON, _ := json.Marshal(listingsPage{
				Status: "success",
				Result: []listing{},
			})
			w.Write(mockJSON)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/v1/reservations", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		query := r.URL.Query()
		offset := query.Get("offset")
		if offset == "0" {
			mockJSON, _ := json.Marshal(reservationsPage{
				Result: []reservation{
					{
						ArrivalDate:           "2021-11-22",
						ChannelID:             1,
						CheckInTime:           17, // Hour later than normal.
						CheckOutTime:          10, // Hour earlier than normal.
						DepartureDate:         "2027-11-26",
						HostawayReservationID: "21107569",
						ListingMapID:          25,
						Status:                "new",
					},
				},
				Status: "success",
			})
			w.Write(mockJSON)
		} else {
			mockJSON, _ := json.Marshal(listingsPage{
				Result: []listing{},
				Status: "success",
			})
			w.Write(mockJSON)

		}
	})
	mux.HandleFunc("/v1/reservations/21107569", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			mockJSON, _ := json.Marshal(reservationUpdateResponse{
				Status:  "",
				Message: "",
			})
			w.Write(mockJSON)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	// Setup a mock HTTP server with the custom mux
	server := httptest.NewServer(mux)
	defer server.Close()

	tz, err := time.LoadLocation("America/Denver")
	assert.Nil(t, err)

	r := hostaway.NewRepository(tz, server.URL)

	unitID := uuid.New()
	reservationsByUnit, err := r.GetForUnits(context.Background(), []shared.Unit{
		{
			ID:                unitID,
			Name:              "01A",
			RemotePropertyURL: "https://dashboard.hostaway.com/listing/25",
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(reservationsByUnit))

	reservations := reservationsByUnit[unitID]
	assert.Equal(t, 1, len(reservations))

	reservation := reservations[0]
	assert.Equal(t, "21107569", reservation.ID)
	assert.Equal(t, "2021-11-22T16:00:00-07:00", reservation.Start.Format(time.RFC3339))
	assert.Equal(t, "2027-11-26T11:00:00-07:00", reservation.End.Format(time.RFC3339))
	assert.Equal(t, "21107569", reservation.TransactionNumber)
}

func loadConfig() error {
	if err := godotenv.Load(".env.test"); err != nil {
		return fmt.Errorf("error loading .env file: %s", err.Error())
	}
	return nil
}
