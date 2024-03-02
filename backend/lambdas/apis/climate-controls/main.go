package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/lambdas/helpers"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/climatecontrol"
	"mlock/lambdas/shared/dynamo/miscellaneous"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ListResponse struct {
	Entities                       []shared.ClimateControl       `json:"entities"`
	ClimateControlOccupiedSettings shared.ClimateControlSettings `json:"climateControlOccupiedSettings"`
	ClimateControlVacantSettings   shared.ClimateControlSettings `json:"climateControlVacantSettings"`
}

type SettingsUpdateRequest struct {
	ClimateControlOccupiedSettings shared.ClimateControlSettings `json:"climateControlOccupiedSettings"`
	ClimateControlVacantSettings   shared.ClimateControlSettings `json:"climateControlVacantSettings"`
}

func main() {
	helpers.StartAPILambda(HandleRequest, []string{helpers.MiddlewareAuth})
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	match, err := regexp.MatchString(`^/climate-controls/settings`, req.Path)
	if err != nil {
		return shared.NewAPIResponse(http.StatusBadRequest, ErrorResponse{Error: "unable to parse request"})
	}
	if match {
		return settingsHandleRequest(ctx, req)
	}

	switch req.HTTPMethod {
	case "GET":
		return list(ctx, req)
	default:
		return shared.NewAPIResponse(http.StatusNotImplemented, "not implemented")
	}
}

func list(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	entities, err := climatecontrol.NewRepository().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting entities: %s", err.Error())
	}

	miscellaneous, ok, err := miscellaneous.NewRepository().Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting miscellaneous: %s", err.Error())
	}
	if !ok {
		return shared.NewAPIResponse(http.StatusNotFound, "miscellaneous not found")
	}

	return shared.NewAPIResponse(
		http.StatusOK, ListResponse{
			Entities:                       entities,
			ClimateControlOccupiedSettings: miscellaneous.ClimateControlOccupiedSettings,
			ClimateControlVacantSettings:   miscellaneous.ClimateControlVacantSettings,
		})
}

func settingsHandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	switch req.HTTPMethod {
	case "PUT":
		return updateSettings(ctx, req)
	default:
		return shared.NewAPIResponse(http.StatusNotImplemented, "not implemented")
	}
}

func updateSettings(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	var body SettingsUpdateRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	miscellaneousRepository := miscellaneous.NewRepository()

	miscellaneous, ok, err := miscellaneousRepository.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting miscellaneous: %s", err.Error())
	}
	if !ok {
		return shared.NewAPIResponse(http.StatusNotFound, "miscellaneous not found")
	}

	miscellaneous.ClimateControlOccupiedSettings = body.ClimateControlOccupiedSettings
	miscellaneous.ClimateControlVacantSettings = body.ClimateControlVacantSettings

	if _, err := miscellaneousRepository.Put(ctx, miscellaneous); err != nil {
		return nil, fmt.Errorf("error putting miscellaneous: %s", err.Error())
	}

	// We might want to kick off something to re-evaluate the current climate control settings.

	return list(ctx, events.APIGatewayProxyRequest{})
}
