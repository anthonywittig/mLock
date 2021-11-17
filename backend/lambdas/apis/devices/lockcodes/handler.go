package lockcodes

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/device"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type CreateRequest struct {
	DeviceID uuid.UUID `json:"deviceId"`
	Code     string    `json:"code"`
	EndAt    time.Time `json:"endAt"`
	StartAt  time.Time `json:"startAt"`
}

type CreateResponse struct {
	Entity shared.Device `json:"entity"`
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

	mlc := shared.DeviceManagedLockCode{
		Code:    body.Code,
		EndAt:   body.EndAt,
		ID:      uuid.New(),
		Status:  shared.DeviceManagedLockCodeStatusScheduled,
		StartAt: body.StartAt,
	}

	d, ok, err := device.Get(ctx, body.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("error getting entity: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("unable to find entity: %s", body.DeviceID)
	}

	if d.HasConflictingManagedLockCode(mlc) {
		return nil, fmt.Errorf("conflicting lock code already exists")
	}

	d.ManagedLockCodes = append(d.ManagedLockCodes, mlc)
	d, err = device.Put(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("error updating device: %s", err.Error())
	}

	// If timestamp is close kick off a request to update the locks? We probably want to keep it single threaded for now.

	return shared.NewAPIResponse(http.StatusOK, CreateResponse{Entity: d})
}
