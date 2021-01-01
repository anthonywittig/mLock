package shared

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type handlerResponse struct {
	response events.APIGatewayProxyResponse
	err      error
}

type simpleBody struct {
	Message string
}

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
		ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()

		c := make(chan handlerResponse)
		go func() {
			c <- handleRequest(ctx, req, handler, middlewares)
		}()

		select {
		case <-ctx.Done():
			// Wait just a little to see if `handleRequest` will finish.
			waitForIt(c, 4*time.Second)
			resp, err := NewAPIResponse(http.StatusGatewayTimeout, simpleBody{Message: "forced timeout"})
			return resp.Proxy, err
		case hr := <-c:
			if hr.err != nil {
				log.Print(fmt.Printf("error handling request: %s", hr.err.Error()))
			}
			return hr.response, hr.err
		}
	}
}

func handleRequest(
	ctx context.Context,
	req events.APIGatewayProxyRequest,
	handler func(context.Context, events.APIGatewayProxyRequest) (*APIResponse, error),
	middlewares []string,
) handlerResponse {

	if err := LoadConfig(); err != nil {
		return handlerResponse{
			response: events.APIGatewayProxyResponse{},
			err:      fmt.Errorf("error loading config: %s", err.Error()),
		}
	}

	if req.HTTPMethod == "OPTIONS" {
		return optionsResponse()
	}

	// Must be done after `LoadConfig`.
	ctx = CreateContextData(ctx)

	// Super lame middleware, maybe we'll need something better one day.
	if err := handleMiddlewares(ctx, req, middlewares); err != nil {
		var apiErr *APIError
		if ok := errors.As(err, &apiErr); ok {
			resp, err := NewAPIResponse(apiErr.StatusCode, apiErr)
			if err != nil {
				return handlerResponse{
					response: events.APIGatewayProxyResponse{},
					err:      fmt.Errorf("error creating api response: %s", err.Error()),
				}
			}
			return handlerResponse{
				response: resp.Proxy,
				err:      nil,
			}
		}

		return handlerResponse{
			response: events.APIGatewayProxyResponse{},
			err:      fmt.Errorf("error handling middleware: %s", err.Error()),
		}
	}

	resp, err := handler(ctx, req)
	if err != nil {
		return handlerResponse{
			response: events.APIGatewayProxyResponse{},
			err:      fmt.Errorf("error executing handler: %s", err.Error()),
		}
	}

	return handlerResponse{
		response: resp.Proxy,
		err:      nil,
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
			return fmt.Errorf("error in lambda, unhandled middleware: %s", middleware)
		}
	}

	return nil
}

func waitForIt(c chan handlerResponse, d time.Duration) {
	ticker := time.NewTicker(d)

	select {
	case <-c:
		log.Print("handler response came soon after we canceled the context")
	case <-ticker.C:
		log.Print("handler response didn't come fast enough, returning without waiting")
	}
}

func optionsResponse() handlerResponse {
	resp, err := NewAPIResponse(http.StatusOK, nil)
	if err != nil {
		return handlerResponse{
			response: events.APIGatewayProxyResponse{},
			err:      fmt.Errorf("error creating options api response: %s", err.Error()),
		}
	}
	return handlerResponse{
		response: resp.Proxy,
		err:      nil,
	}
}
