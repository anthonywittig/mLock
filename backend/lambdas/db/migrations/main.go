package main

import (
	"context"
	"database/sql"
	"fmt"
	"mlock/shared"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pressly/goose"
)

type MyEvent struct {
}

type Response struct {
	Messages []string `json:"Messages"`
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, event MyEvent) (Response, error) {
	if err := shared.LoadConfig(); err != nil {
		return Response{}, fmt.Errorf("error loading config: %s", err.Error())
	}

	ctx = shared.CreateContextData(ctx)

	db, err := shared.GetDB(ctx)
	if err != nil {
		return Response{}, fmt.Errorf("error opening DB: %s", err.Error())
	}
	defer db.Close()

	if err = db.PingContext(ctx); err != nil {
		return Response{}, fmt.Errorf("error pinging DB: %s", err.Error())
	}

	if err := migrateUp(db); err != nil {
		return Response{}, fmt.Errorf("error migrating DB: %s", err.Error())
	}

	dbName := shared.GetCurrentDatabase(ctx)

	return Response{
		Messages: []string{
			"current db",
			dbName,
		},
	}, nil
}

func migrateUp(db *sql.DB) error {
	goose.SetTableName("goose_db_version")
	return goose.Up(db, "goosemigrations")
}
