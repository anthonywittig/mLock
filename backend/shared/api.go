package shared

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func APIResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
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
	return resp, nil
}
