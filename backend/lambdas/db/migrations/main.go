package main

import (
	"context"
	"database/sql"
	"fmt"
	"mlock/shared"
	"mlock/shared/dynamo/user"
	"mlock/shared/postgres"

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

	db, err := postgres.GetDB(ctx)
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

	if err := user.MigrateUsers(ctx); err != nil {
		return Response{}, fmt.Errorf("error migrating dynamo users: %s", err.Error())
	}
	/*
		//
		// code for moving users from postgres to dynamo
		//
		users, err := pguser.GetAll(ctx)
		if err != nil {
			return Response{}, fmt.Errorf("error getting users: %s", err.Error())
		}
		errors := []error{}
		for _, u := range users {
			cd, err := shared.GetContextData(ctx)
			if err != nil {
				return Response{}, fmt.Errorf("error getting context data: %s", err.Error())
			}
			cd.User = &shared.User{Email: u.Email, CreatedBy: u.CreatedBy}
			if err := user.Put(ctx, u.Email); err != nil {
				errors = append(errors, err)
			}
		}
		if len(errors) > 0 {
			return Response{}, fmt.Errorf("error(s) inserting user(s): %+v", errors)
		}
		//
		// end code for moving users from postgres to dynamo
		//
	*/

	dbName := postgres.GetCurrentDatabase(ctx)

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
