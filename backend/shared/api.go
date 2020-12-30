package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type APIResponse struct {
	Proxy events.APIGatewayProxyResponse
}

const (
	SetCookieHeaderName = "Set-Cookie"
)

func NewAPIResponse(status int, body interface{}) (*APIResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{
		"Content-Type": "application/json",
		// pull in from config...
		"Access-Control-Allow-Origin":      GetConfig("FRONTEND_DOMAIN"),
		"Access-Control-Allow-Credentials": "true",
	}}
	resp.StatusCode = status
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshalling body: %s", err.Error())
	}
	resp.Body = string(jsonBody)
	return &APIResponse{Proxy: resp}, nil
}

func (a *APIResponse) AddCookie(name string, value string) error {
	if _, exists := a.Proxy.Headers[SetCookieHeaderName]; exists {
		return errors.New("need to implement multiple cookie support")
	}

	// TODO: we should check cookie values for illegal characters.

	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   60 * 60 * 24, // Seems like browsers handle session cookies differently, so we'll just set an expiration.
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	a.Proxy.Headers[SetCookieHeaderName] = cookie.String()

	return nil
}
