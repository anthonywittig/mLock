package hostaway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mlock/lambdas/shared"
	"net/http"
	"net/url"
	"strings"
	"time"

	mshared "mlock/shared"

	"github.com/google/uuid"
)

type authData struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type listing struct {
	DoorSecurityCode    string `json:"doorSecurityCode"`
	ID                  int    `json:"id"`
	InternalListingName string `json:"internalListingName"`
}

type listingsPage struct {
	Status string    `json:"status"`
	Result []listing `json:"result"`
	// Count  int       `json:"count"`
	// Limit  int       `json:"limit"`
	// Offset int       `json:"offset"`
}

type reservation struct {
	ArrivalDate           string `json:"arrivalDate"`
	DoorCode              string `json:"doorCode"`
	ChannelID             int    `json:"channelId"`
	CheckInTime           int    `json:"checkInTime"`
	CheckOutTime          int    `json:"checkOutTime"`
	DepartureDate         string `json:"departureDate"`
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
	// Count  int           `json:"count"`
	// Limit  int           `json:"limit"`
	// Offset int           `json:"offset"`
}

type Repository struct {
	hostawayURL string
	timeZone    *time.Location
}

func NewRepository(timeZone *time.Location, hostawayURL string) *Repository {
	if hostawayURL == "" {
		hostawayURL = "https://api.hostaway.com"
	}

	return &Repository{
		hostawayURL: hostawayURL,
		timeZone:    timeZone,
	}
}

func (r *Repository) Get(ctx context.Context, unit shared.Unit) ([]shared.Reservation, error) {
	reservations, err := r.GetForUnits(ctx, []shared.Unit{unit})
	if err != nil {
		return nil, fmt.Errorf("error getting reservations: %s", err.Error())
	}
	return reservations[unit.ID], nil
}

func (r *Repository) GetForUnits(ctx context.Context, units []shared.Unit) (map[uuid.UUID][]shared.Reservation, error) {
	accessToken, err := r.getAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting access token: %s\n", err.Error())
	}

	reservations, err := r.getReservations(ctx, accessToken, units)
	if err != nil {
		return nil, fmt.Errorf("error getting reservations: %s\n", err.Error())
	}

	var listings map[int]listing = map[int]listing{}

	reservationsByUnit := map[uuid.UUID][]shared.Reservation{}

	for _, reservation := range reservations {
		for _, unit := range units {
			if unit.GetRemotePropertyID() == reservation.ListingMapID {
				startDate, err := time.ParseInLocation("2006-01-02", reservation.ArrivalDate, r.timeZone)
				if err != nil {
					return map[uuid.UUID][]shared.Reservation{}, fmt.Errorf("error parsing start date: %s", err.Error())
				}
				// If they say they're going to be later than 4pm, assume 4pm.
				checkInHour := min(reservation.CheckInTime, 16)
				startDate = startDate.Add(time.Duration(checkInHour) * time.Hour)

				endDate, err := time.ParseInLocation("2006-01-02", reservation.DepartureDate, r.timeZone)
				if err != nil {
					return map[uuid.UUID][]shared.Reservation{}, fmt.Errorf("error parsing end date: %s", err.Error())
				}
				// If they say they're going to be earlier than 11am, assume 11am.
				checkOutHour := max(reservation.CheckOutTime, 11)
				endDate = endDate.Add(time.Duration(checkOutHour) * time.Hour)

				// This probably isn't the best place to do this, but if the `DoorCode` isn't set, set it.
				if reservation.DoorCode == "" {
					if len(listings) == 0 {
						listings, err = r.getListingsByID(ctx, accessToken)
						if err != nil {
							return map[uuid.UUID][]shared.Reservation{}, fmt.Errorf("error getting listings: %s", err.Error())
						}
					}
					listing, ok := listings[reservation.ListingMapID]
					if ok {
						reservation.DoorCode = listing.DoorSecurityCode
						fmt.Printf("setting door code to listing door code: %s\n", reservation.DoorCode)
					} else {
						reservation.DoorCode = reservation.HostawayReservationID[len(reservation.HostawayReservationID)-4:]
					}
					if err := r.setDoorCode(ctx, accessToken, reservation.HostawayReservationID, reservation.DoorCode); err != nil {
						return map[uuid.UUID][]shared.Reservation{}, fmt.Errorf("error setting door code: %s", err.Error())
					}
				}

				reservationsByUnit[unit.ID] = append(reservationsByUnit[unit.ID], shared.Reservation{
					ID:                reservation.HostawayReservationID,
					DoorCode:          reservation.DoorCode,
					TransactionNumber: reservation.HostawayReservationID,
					Start:             startDate,
					End:               endDate,
				})
				break
			}
		}
	}

	return reservationsByUnit, nil
}

func (r *Repository) getAccessToken(ctx context.Context) (authData, error) {
	accountId, err := mshared.GetConfig("HOSTAWAY_ACCOUNT_ID")
	if err != nil {
		return authData{}, fmt.Errorf("error getting accountId: %s", err.Error())
	}

	apiKey, err := mshared.GetConfig("HOSTAWAY_API_KEY")
	if err != nil {
		return authData{}, fmt.Errorf("error getting apiKey: %s", err.Error())
	}

	bodyData := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {accountId},
		"client_secret": {apiKey},
		"scope":         {"general"},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/v1/accessTokens", r.hostawayURL),
		strings.NewReader(bodyData.Encode()),
	)
	if err != nil {
		return authData{}, fmt.Errorf("error creating request: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return authData{}, fmt.Errorf("error doing request: %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return authData{}, fmt.Errorf("error reading body: %s", err.Error())
	}

	body := authData{}
	if err := json.Unmarshal(respBody, &body); err != nil {
		return authData{}, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	return body, nil
}

func (r *Repository) getListingsByID(ctx context.Context, authToken authData) (map[int]listing, error) {
	result := map[int]listing{}
	page := 0
	for {
		pageResult, err := getPage(listingsPage{}, r, ctx, authToken, "listings", page, []string{})
		if err != nil {
			return map[int]listing{}, fmt.Errorf("error getting listings page: %s", err.Error())
		}
		if pageResult.Status != "success" {
			return map[int]listing{}, fmt.Errorf("error getting listings page, non-success status: %s", pageResult.Status)
		}

		if len(pageResult.Result) == 0 {
			break
		}
		if page > 20 {
			return map[int]listing{}, fmt.Errorf("too many pages")
		}

		for _, listing := range pageResult.Result {
			result[listing.ID] = listing
		}
		page++
	}
	return result, nil
}

func (r *Repository) getReservations(ctx context.Context, authToken authData, units []shared.Unit) ([]reservation, error) {
	twoDaysAgo := time.Now().Add(-48 * time.Hour).Format("2006-01-02")
	extraParameters := []string{"sortOrder=arrivalDate"}
	if len(units) == 1 {
		extraParameters = append(extraParameters, fmt.Sprintf("listingId=%d", units[0].GetRemotePropertyID()))
	}
	result := []reservation{}
	page := 0
	for {
		pageResult, err := getPage(reservationsPage{}, r, ctx, authToken, "reservations", page, extraParameters)
		if err != nil {
			return []reservation{}, fmt.Errorf("error getting reservations page: %s", err.Error())
		}
		if pageResult.Status != "success" {
			return []reservation{}, fmt.Errorf("error getting reservations page, non-success status: %s", pageResult.Status)
		}

		if len(pageResult.Result) == 0 {
			return result, nil
		}
		if page > 20 {
			return []reservation{}, fmt.Errorf("too many pages")
		}
	ReservationLoop:
		for _, reservation := range pageResult.Result {
			for _, statusToIgnore := range []string{
				"cancelled",
				// Need to see if we should ignore these.
				// "declined",
				// "inquiry",
				// "inquiryNotPossible",
			} {
				if reservation.Status == statusToIgnore {
					continue ReservationLoop
				}
			}
			if reservation.Status != "new" && reservation.Status != "modified" {
				fmt.Printf("unhandled reservation status: %s for hostaway reservation ID: %s\n", reservation.Status, reservation.HostawayReservationID)
			}
			if reservation.DepartureDate < twoDaysAgo {
				continue
			}
			result = append(result, reservation)
		}
		page++
	}
}

func getPage[T any](
	emptyT T,
	r *Repository,
	ctx context.Context,
	authToken authData,
	resource string,
	page int,
	queryParams []string,
) (T, error) {
	limit := 100
	offset := page * limit

	extraQueryParams := ""
	for _, param := range queryParams {
		extraQueryParams += "&" + param
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/v1/%s?limit=%d&offset=%d%s", r.hostawayURL, resource, limit, offset, extraQueryParams),
		nil,
	)
	if err != nil {
		return emptyT, fmt.Errorf("error creating request: %s", err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+authToken.AccessToken)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return emptyT, fmt.Errorf("error doing request: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return emptyT, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return emptyT, fmt.Errorf("error reading body: %s", err.Error())
	}

	body := emptyT
	if err := json.Unmarshal(respBody, &body); err != nil {
		return emptyT, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	return body, nil
}

func (r *Repository) setDoorCode(ctx context.Context, authToken authData, reservationID string, doorCode string) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf("%s/v1/reservations/%s", r.hostawayURL, reservationID),
		strings.NewReader(fmt.Sprintf(`{"doorCode":"%s"}`, doorCode)),
	)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+authToken.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error doing request: %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading body: %s", err.Error())
	}

	var body reservationUpdateResponse
	if err := json.Unmarshal(respBody, &body); err != nil {
		return fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	if resp.StatusCode == 403 && strings.HasPrefix(body.Message, "Requested dates are not available") {
		// There were a lot of duplicate reservations when we first moved over to Hostaway...
		fmt.Printf("non-fatal error setting door code for reservation %s; message: %s\n", reservationID, body.Message)
		return nil
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("non-200 status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
