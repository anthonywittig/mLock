package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mlock/shared"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type CreateBody struct {
	GoogleToken string
}

type CreateResponse struct {
	Message string
}

type DeleteResponse struct {
	Message string
}

func main() {
	shared.StartAPILambda(HandleRequest, []string{})
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	switch req.HTTPMethod {
	case "DELETE":
		return delete(ctx, req)
	case "OPTIONS":
		return shared.NewAPIResponse(http.StatusOK, nil)
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
	tokenData, err := shared.GetUserFromToken(ctx, body.GoogleToken)
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

	resp, err := shared.NewAPIResponse(http.StatusOK, CreateResponse{Message: "ok"})
	if err != nil {
		return nil, err
	}

	if err := resp.SetAuthCookie(body.GoogleToken); err != nil {
		return nil, err
	}

	return resp, nil
}

func delete(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	// For now, we just need to delete the cookie.
	resp, err := shared.NewAPIResponse(http.StatusOK, DeleteResponse{Message: "deleted session"})
	if err != nil {
		return nil, err
	}

	if err := resp.DeleteAuthCookie(); err != nil {
		return nil, err
	}

	return resp, nil
}
