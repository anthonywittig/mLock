package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/shared"
	"mlock/shared/token"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "POST":
		return create(ctx, req)
	default:
		return shared.APIResponse(http.StatusNotImplemented, fmt.Errorf("not implemented - %s", req.HTTPMethod))
	}
}

func create(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var body CreateBody
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	// Validate the token.
	_, err := token.GetUserFromToken(nil, body.GoogleToken)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %s", err.Error())
	}

	// In the future we could create our own token, but for now we'll just piggy back on Google's.
	return shared.APIResponse(http.StatusOK, CreateResponse{Token: body.GoogleToken})
}
