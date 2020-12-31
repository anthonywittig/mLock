package shared

import (
	"context"
	"errors"
	"fmt"
	"log"

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

		resp, err := handler(ctx, req)
		if err != nil {
			log.Printf("error in lambda: %s\n", err.Error())
		}
		return resp.Proxy, err
	}
}

func handleMiddlewares(ctx context.Context, req events.APIGatewayProxyRequest, middlewares []string) (string, error) {
	// Add stuff to context here?

	for _, middleware := range middlewares {
		switch middleware {
		case MiddlewareAuth:
			// wait for it
		default:
			log.Printf("error in lambda, unhandled middleware: %s\n", middleware)
			return "", errors.New("unhandled middleware")
		}
	}

	return "", nil
}
