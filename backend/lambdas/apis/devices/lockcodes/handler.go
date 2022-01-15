package lockcodes

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/device"
	"net/http"
	"regexp"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type CreateRequest struct {
	Code    string    `json:"code"`
	EndAt   time.Time `json:"endAt"`
	StartAt time.Time `json:"startAt"`
}

type CreateResponse struct {
	Entity shared.Device `json:"entity"`
}

type UpdateBody struct {
	EndAt time.Time `json:"endAt"`
}

type UpdateResponse struct {
	Entity shared.Device `json:"entity"`
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*shared.APIResponse, error) {
	r, err := regexp.Compile(`^/devices/([0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12})/lock-codes/([0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12})?`)
	if err != nil {
		return nil, fmt.Errorf("error generating regex: %s", err.Error())
	}

	match := r.FindStringSubmatch(req.Path)

	if len(match) < 2 {
		return nil, fmt.Errorf("regex didn't match path")
	}

	deviceID, err := uuid.Parse(match[1])
	if err != nil {
		return nil, fmt.Errorf("error parsing device id: %s", err.Error())
	}

	d, ok, err := device.NewRepository().Get(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("error getting entity: %s", err.Error())
	}
	if !ok {
		return nil, fmt.Errorf("unable to find entity: %s", deviceID)
	}

	switch req.HTTPMethod {
	case "POST":
		return create(ctx, req, d)
	case "PUT":
		if len(match) != 3 {
			return nil, fmt.Errorf("regex didn't match path for PUT")
		}

		mlcID, err := uuid.Parse(match[2])
		if err != nil {
			return nil, fmt.Errorf("error parsing device id: %s", err.Error())
		}

		return update(ctx, req, d, mlcID)
	default:
		return shared.NewAPIResponse(http.StatusNotImplemented, "not implemented")
	}
}

func create(ctx context.Context, req events.APIGatewayProxyRequest, d shared.Device) (*shared.APIResponse, error) {
	// TODO: create audit log.

	var body CreateRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get context data: %s", err.Error())
	}

	currentUser := cd.User
	if currentUser == nil {
		return nil, fmt.Errorf("no current user")
	}

	mlc := &shared.DeviceManagedLockCode{
		Code:    body.Code,
		EndAt:   body.EndAt,
		ID:      uuid.New(),
		Note:    fmt.Sprintf("Added by %s.", currentUser.Email),
		Status:  shared.DeviceManagedLockCodeStatus1Scheduled,
		StartAt: body.StartAt,
	}

	if mlc.StartAt.After(mlc.EndAt) {
		return nil, fmt.Errorf("can't start after it ends")
	}

	d.ManagedLockCodes = append(d.ManagedLockCodes, mlc)

	if err := device.NewRepository().AppendToAuditLog(ctx, d, []*shared.DeviceManagedLockCode{mlc}); err != nil {
		return nil, fmt.Errorf("error appending to audit log: %s", err.Error())
	}

	d, err = device.NewRepository().Put(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("error updating device: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, CreateResponse{Entity: d})
}

func update(ctx context.Context, req events.APIGatewayProxyRequest, d shared.Device, mlcID uuid.UUID) (*shared.APIResponse, error) {
	var body UpdateBody
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	mlc := d.GetManagedLockCode(mlcID)
	if mlc == nil {
		return nil, fmt.Errorf("unable to find managed lock code")
	}

	if mlc.ReservationID != "" {
		return nil, fmt.Errorf("can't edit a reservation lock code")
	}

	mlc.EndAt = body.EndAt
	if mlc.StartAt.After(mlc.EndAt) {
		return nil, fmt.Errorf("can't start after it ends")
	}

	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get context data: %s", err.Error())
	}

	currentUser := cd.User
	if currentUser == nil {
		return nil, fmt.Errorf("no current user")
	}

	mlc.Note = fmt.Sprintf("Edited by %s.", currentUser.Email)

	if device.NewRepository().AppendToAuditLog(ctx, d, []*shared.DeviceManagedLockCode{mlc}); err != nil {
		return nil, fmt.Errorf("error appending to audit log: %s", err.Error())
	}

	d, err = device.NewRepository().Put(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("error updating device: %s", err.Error())
	}

	return shared.NewAPIResponse(http.StatusOK, UpdateResponse{Entity: d})
}
