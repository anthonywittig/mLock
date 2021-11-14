package lockcodes

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/lambdas/shared"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type CreateRequest struct {
	Code    string    `json:"code"`
	EndAt   time.Time `json:"endAt"`
	StartAt time.Time `json:"startAt"`
}

type CreateResponse struct {
	//Entity shared.Property `json:"entity"`
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	switch req.HTTPMethod {
	case "POST":
		return create(ctx, req)
	default:
		return shared.NewAPIResponse(http.StatusNotImplemented, "not implemented")
	}
}

func create(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	var body CreateRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	// Package up the lock code.

	// Get the device.

	// Check for conflicting entry (same code with overlapping time ranges).
	//device.HasConflictingManagedLockCode

	// Save.

	// If timestamp is close kick off a request to update the locks? We probably want to keep it single threaded for now.

	/*
		entity, err := property.Put(ctx, shared.Property{
			ID: uuid.New(),
			Name: body.Name,
		})
		if err != nil {
			return nil, fmt.Errorf("error inserting entity: %s", err.Error())
		}
	*/

	//return shared.NewAPIResponse(http.StatusOK, CreateResponse{Entity: entity})
	return shared.NewAPIResponse(http.StatusOK, CreateResponse{})
}
