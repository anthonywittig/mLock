package shared

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	MiddlewareAuth = "middlewareAuth"
)

func StartAPILambda(
	handler func(context.Context, events.APIGatewayProxyRequest) (*APIResponse, error),
	middlewares []string,
) {
	lambda.Start(handlerWrapper(handler, middlewares))
}

func handlerWrapper(
	handler func(context.Context, events.APIGatewayProxyRequest) (*APIResponse, error),
	middlewares []string,
) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		if err := LoadConfig(); err != nil {
			err = fmt.Errorf("error loading config: %s", err.Error())
			log.Print(err.Error())
			return events.APIGatewayProxyResponse{}, err
		}

		ctx = CreateContextData(ctx)

		// Super lame middleware, maybe we'll need something better one day.
		err := handleMiddlewares(ctx, req, middlewares)
		if err != nil {
			var apiErr *APIError
			if ok := errors.As(err, &apiErr); ok {
				resp, err := NewAPIResponse(apiErr.StatusCode, apiErr)
				if err != nil {
					err = fmt.Errorf("error creating api response: %s", err.Error())
					log.Print(err.Error())
					return events.APIGatewayProxyResponse{}, err
				}
				return resp.Proxy, nil
			}

			err = fmt.Errorf("error handling middleware: %s", err.Error())
			log.Print(err.Error())
			return events.APIGatewayProxyResponse{}, err
		}

		resp, err := handler(ctx, req)
		if err != nil {
			log.Printf("error in lambda: %s\n", err.Error())
		}
		return resp.Proxy, err
	}
}

func handleMiddlewares(ctx context.Context, req events.APIGatewayProxyRequest, middlewares []string) error {
	if err := AddAuthToContext(ctx, req); err != nil {
		return fmt.Errorf("error adding auth: %s", err.Error())
	}

	for _, middleware := range middlewares {
		switch middleware {
		case MiddlewareAuth:
			user, err := GetAuthUser(ctx)
			if err != nil {
				return fmt.Errorf("error getting auth user: %s", err.Error())
			}
			if user == nil {
				return &APIError{
					StatusCode: http.StatusUnauthorized,
					Message:    "unauthorized",
				}
			}
		default:
			log.Printf("error in lambda, unhandled middleware: %s\n", middleware)
			return errors.New("unhandled middleware")
		}
	}

	return nil
}
