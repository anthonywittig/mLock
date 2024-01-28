package main

import (
	"context"
	"fmt"
	"mlock/lambdas/helpers"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/climatecontrol"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ListResponse struct {
	Entities []shared.ClimateControl `json:"entities"`
}

func main() {
	helpers.StartAPILambda(HandleRequest, []string{helpers.MiddlewareAuth})
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
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

	return shared.NewAPIResponse(http.StatusOK, ListResponse{Entities: entities})
}
