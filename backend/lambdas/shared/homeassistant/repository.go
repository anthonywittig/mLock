package homeassistant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"mlock/lambdas/shared"
	mshared "mlock/shared"
)

type statesListResponseEntity struct {
	EntityID string `json:"entity_id"`
}

type Repository struct {
	authToken string
	baseURL   string
}

func NewRepository() (*Repository, error) {
	authToken, err := mshared.GetConfig("HOME_ASSISTANT_AUTH_TOKEN")
	if err != nil {
		return nil, fmt.Errorf("error getting HOME_ASSISTANT_AUTH_TOKEN: %s", err.Error())
	}

	baseURL, err := mshared.GetConfig("HOME_ASSISTANT_BASE_URL")
	if err != nil {
		return nil, fmt.Errorf("error getting HOME_ASSISTANT_BASE_URL: %s", err.Error())
	}

	return &Repository{
		authToken: authToken,
		baseURL:   baseURL,
	}, nil
}

func (r *Repository) GetClimateControl(ctx context.Context, id string) (shared.RawClimateControl, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/api/states/%s", r.baseURL, id),
		nil,
	)
	if err != nil {
		return shared.RawClimateControl{}, fmt.Errorf("error creating request: %s", err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+r.authToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return shared.RawClimateControl{}, fmt.Errorf("error doing request: %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return shared.RawClimateControl{}, fmt.Errorf("error reading body: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		return shared.RawClimateControl{}, fmt.Errorf("non-200 status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	var body shared.RawClimateControl
	if err := json.Unmarshal(respBody, &body); err != nil {
		return shared.RawClimateControl{}, fmt.Errorf(
			"error unmarshalling body: %s; %s",
			err.Error(),
			respBody,
		)
	}

	return body, nil
}

func (r *Repository) ListClimateControls(ctx context.Context) ([]shared.RawClimateControl, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/api/states", r.baseURL),
		nil,
	)
	if err != nil {
		return []shared.RawClimateControl{}, fmt.Errorf("error creating request: %s", err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+r.authToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return []shared.RawClimateControl{}, fmt.Errorf("error doing request: %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []shared.RawClimateControl{}, fmt.Errorf("error reading body: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		return []shared.RawClimateControl{}, fmt.Errorf("non-200 status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	var body []statesListResponseEntity
	if err := json.Unmarshal(respBody, &body); err != nil {
		return []shared.RawClimateControl{}, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	climateControls := []shared.RawClimateControl{}
	for _, entity := range body {
		if entity.EntityID[:7] != "climate" {
			continue
		}
		// Could parallelize this, but maybe we'll be a nicer client if we don't.
		cc, err := r.GetClimateControl(ctx, entity.EntityID)
		if err != nil {
			return []shared.RawClimateControl{}, fmt.Errorf(
				"error getting climate control for entityID: %s; %s",
				entity.EntityID,
				err.Error(),
			)
		}
		climateControls = append(climateControls, cc)
	}

	return climateControls, nil
}

func (r *Repository) SetToDesiredState(ctx context.Context, climateControl shared.ClimateControl) error {
	if climateControl.DesiredState.HVACMode != climateControl.ActualState.HVACMode {
		if err := r.SetHVACMode(ctx, climateControl, climateControl.DesiredState.HVACMode); err != nil {
			return fmt.Errorf("error setting HVAC mode: %s", err.Error())
		}
	}

	if climateControl.DesiredState.Temperature != climateControl.ActualState.Temperature {
		if err := r.SetTemperature(ctx, climateControl, climateControl.DesiredState.Temperature); err != nil {
			return fmt.Errorf("error setting temperature: %s", err.Error())
		}
	}

	return nil
}

func (r *Repository) SetHVACMode(ctx context.Context, climateControl shared.ClimateControl, hvacMode string) error {
	jsonData, err := json.Marshal(struct {
		EntityID string `json:"entity_id"`
		HVACMode string `json:"hvac_mode"`
	}{
		EntityID: climateControl.RawClimateControl.EntityID,
		HVACMode: hvacMode,
	})
	if err != nil {
		return fmt.Errorf("error marshaling data: %s", err.Error())
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/services/climate/set_hvac_mode", r.baseURL),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+r.authToken)
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

	if resp.StatusCode != 200 {
		return fmt.Errorf("non-200 status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (r *Repository) SetTemperature(ctx context.Context, climateControl shared.ClimateControl, temperature int) error {
	jsonData, err := json.Marshal(struct {
		EntityID    string `json:"entity_id"`
		Temperature int    `json:"temperature"`
	}{
		EntityID:    climateControl.RawClimateControl.EntityID,
		Temperature: temperature,
	})
	if err != nil {
		return fmt.Errorf("error marshaling data: %s", err.Error())
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/services/climate/set_temperature", r.baseURL),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+r.authToken)
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

	if resp.StatusCode != 200 {
		return fmt.Errorf("non-200 status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
