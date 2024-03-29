package shared

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type ContextData struct {
	DB   *sql.DB
	DY   *dynamodb.Client
	SES  *ses.Client
	SQS  *sqs.Client
	User *User
}

type contextKey int

const (
	contextDataKey contextKey = iota
)

func GetContextData(ctx context.Context) (*ContextData, error) {
	val := ctx.Value(contextDataKey)
	if val == nil {
		return nil, fmt.Errorf("context data is not set")
	}

	data, ok := val.(*ContextData)
	if !ok {
		return nil, fmt.Errorf("context data has wrong type")
	}

	return data, nil
}

func CreateContextData(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextDataKey, &ContextData{})
}
