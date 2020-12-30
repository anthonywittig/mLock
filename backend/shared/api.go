package shared

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type APIResponse struct {
	Proxy *events.APIGatewayProxyResponse
}

func NewAPIResponse(status int, body interface{}) (*APIResponse, error) {
	resp := &events.APIGatewayProxyResponse{Headers: map[string]string{
		"Content-Type":                "application/json",
		"access-control-allow-origin": "*",
	}}
	resp.StatusCode = status
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshalling body: %s", err.Error())
	}
	resp.Body = string(jsonBody)
	return &APIResponse{Proxy: resp}, nil
}
