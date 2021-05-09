package shared

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mlock/shared"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type APIResponse struct {
	Proxy events.APIGatewayProxyResponse
}

const (
	CookieHeaderName    = "cookie"
	AuthCookie          = "auth-v1"
	SetCookieHeaderName = "Set-Cookie"
)

func NewAPIResponse(status int, body interface{}) (*APIResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{
		"Content-Type": "application/json",
		// Not all of the CORS headers need to be in every request, but to make things easy we'll include them all.
		"Access-Control-Allow-Methods":     "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH", // Can't use `*` with credentials.
		"Access-Control-Allow-Origin":      shared.GetConfig("FRONTEND_DOMAIN"),
		"Access-Control-Allow-Credentials": "true",
	}}
	resp.StatusCode = status

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshalling body: %s", err.Error())
		}
		resp.Body = string(jsonBody)
	}

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
		Path:     "/",
		MaxAge:   60 * 60 * 24, // Seems like browsers handle session cookies differently, so we'll just set an expiration.
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	a.Proxy.Headers[SetCookieHeaderName] = cookie.String()

	return nil
}

func (a *APIResponse) DeleteAuthCookie() error {
	// Clearing isn't exactly the same as deleting but we need to handle a messed up cookie anyway.
	return a.AddCookie(AuthCookie, "")
}

func (a *APIResponse) SetAuthCookie(token string) error {
	return a.AddCookie(AuthCookie, token)
}

func AddAuthToContext(ctx context.Context, req events.APIGatewayProxyRequest, userService UserService) error {
	cookies := req.Headers[CookieHeaderName]
	if cookies == "" {
		return nil
	}

	authCookieValue := ""
	for _, cookie := range strings.Split(cookies, ";") {
		cookie = strings.TrimSpace(cookie)
		cookieParts := strings.SplitN(cookie, "=", 2)
		if len(cookieParts) != 2 {
			return fmt.Errorf("unexpected cookie format: %s", cookie)
		}

		if cookieParts[0] == AuthCookie {
			authCookieValue = cookieParts[1]
			break
		}
	}
	if authCookieValue == "" {
		return nil
	}

	tokenData, err := GetUserFromToken(ctx, authCookieValue, userService)
	if err != nil {
		return fmt.Errorf("error getting user from token: %s", err.Error())
	}
	if tokenData.Error != nil || !tokenData.TokenValid || !tokenData.UserValid {
		// Could probably just check if the user is not nil.
		return nil
	}

	cd, err := GetContextData(ctx)
	if err != nil {
		return fmt.Errorf("error getting context data: %s", err.Error())
	}

	cd.User = &tokenData.User

	return nil
}

func GetAuthUser(ctx context.Context) (*User, error) {
	cd, err := GetContextData(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting context data: %s", err.Error())
	}

	return cd.User, nil
}
