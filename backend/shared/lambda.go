package shared

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func StartAPILambda(handler func(context.Context, events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)) {
	lambda.Start(handlerWrapper(handler))
}

func handlerWrapper(
	handler func(context.Context, events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error),
) func(context.Context, events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	return func(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			log.Printf("error in lambda: %s\n", err.Error())
		}
		return resp, err
	}
}
