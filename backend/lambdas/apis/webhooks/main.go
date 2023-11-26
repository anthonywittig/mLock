package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/lambdas/helpers"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/sqs"

	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type hostawayRequest struct {
	Object *string `json:"object"`
	Event  *string `json:"event"`
}

type emptyResponse struct{}

func main() {
	helpers.StartAPILambda(HandleRequest, []string{helpers.MiddlewareAuth})
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	// When the webhook is first created, Hostaway sends a GET and POST request to verify the endpoint. I'm not sure if we need to support the GET.
	switch req.HTTPMethod {
	case "GET":
		// WARNING: the Hostaway webhook verification GET request does not include any auth, do not trust this.
		return shared.NewAPIResponse(http.StatusOK, emptyResponse{})
	case "POST":
		return post(ctx, req)
	default:
		return shared.NewAPIResponse(http.StatusNotImplemented, "not implemented")
	}
}

func post(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	var body hostawayRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	fmt.Printf("raw body: %+v\n", req.Body)
	if body.Object != nil {
		fmt.Printf("webhook body.object: %+v\n", *body.Object)
	}

	// When a webhook is created, Hostaway sends a POST request without an `object` or `event`. We also want to ignore any non-reservation requests.
	if body.Object == nil || *body.Object != "reservation" {
		return shared.NewAPIResponse(http.StatusOK, emptyResponse{})
	}

	sqsService, err := sqs.NewSQSService(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting sqs service: %s", err.Error())
	}

	if err := sqsService.SendBlankMessageToPollSchedulesQueue(ctx); err != nil {
		return nil, fmt.Errorf("error sending blank message to poll schedules queue: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, emptyResponse{})
}
