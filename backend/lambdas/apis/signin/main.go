package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mlock/shared"
	"mlock/shared/token"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type CreateBody struct {
	GoogleToken string
}

type CreateResponse struct {
	Token string
}

func main() {
	if err := shared.LoadConfig(); err != nil {
		fmt.Printf("ERROR loading config: %s\n", err.Error())
		return
	}
	shared.StartAPILambda(HandleRequest)
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	switch req.HTTPMethod {
	case "POST":
		return create(ctx, req)
	default:
		return shared.NewAPIResponse(http.StatusNotImplemented, fmt.Errorf("not implemented - %s", req.HTTPMethod))
	}
}

func create(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	var body CreateBody
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	// Validate the token.
	tokenData, err := token.GetUserFromToken(nil, body.GoogleToken)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %s", err.Error())
	}

	if tokenData.Error != nil {
		var apiErr *shared.APIError
		if ok := errors.As(tokenData.Error, &apiErr); ok {
			return shared.NewAPIResponse(apiErr.StatusCode, apiErr)
		}

		return nil, tokenData.Error
	}

	return shared.NewAPIResponse(http.StatusOK, CreateResponse{Token: body.GoogleToken})
}
