package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/property"
	"mlock/lambdas/shared/dynamo/unit"
	"mlock/lambdas/shared/dynamo/user"
	mshared "mlock/shared"
	"time"

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
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	log.Printf("starting migrations\n")
	if err := mshared.LoadConfig(); err != nil {
		return Response{}, fmt.Errorf("error loading config: %s", err.Error())
	}

	ctx = shared.CreateContextData(ctx)

	log.Printf("migrating user...\n")
	if err := user.Migrate(ctx); err != nil {
		return Response{}, fmt.Errorf("error migrating dynamo users: %s", err.Error())
	}
	log.Printf("migrated user\n")
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

	if err := property.Migrate(ctx); err != nil {
		return Response{}, fmt.Errorf("error migrating dynamo properties: %s", err.Error())
	}
	/*
		//
		// code for moving properties from postgres to dynamo
		//
		items, err := pgProperty.GetAll(ctx)
		if err != nil {
			return Response{}, fmt.Errorf("error getting items: %s", err.Error())
		}
		errors := []error{}
		for _, i := range items {
			cd, err := shared.GetContextData(ctx)
			if err != nil {
				return Response{}, fmt.Errorf("error getting context data: %s", err.Error())
			}
			cd.User = &shared.User{Email: i.CreatedBy}
			if _, err := property.PutID(ctx, i.Name, i.ID); err != nil {
				errors = append(errors, err)
			}
		}
		if len(errors) > 0 {
			return Response{}, fmt.Errorf("error(s) inserting items(s): %+v", errors)
		}
		//
		// end code for moving properties from postgres to dynamo
		//
	*/

	if err := unit.Migrate(ctx); err != nil {
		return Response{}, fmt.Errorf("error migrating dynamo properties: %s", err.Error())
	}
	//
	// code for moving properties from postgres to dynamo
	//
	/*
		properties, err := property.List(ctx)
		if err != nil {
			return Response{}, fmt.Errorf("error getting properties: %s", err.Error())
		}

		items, err := pgUnit.GetAll(ctx)
		if err != nil {
			return Response{}, fmt.Errorf("error getting items: %s", err.Error())
		}
		errors := []error{}
		for _, i := range items {
			cd, err := shared.GetContextData(ctx)
			if err != nil {
				return Response{}, fmt.Errorf("error getting context data: %s", err.Error())
			}
			cd.User = &shared.User{Email: i.UpdatedBy}
			found := false
			for _, p := range properties {
				if p.ID == i.PropertyID.String() {
					found = true
					//if _, err := unit.PutCal(ctx, i.Name, p.Name, i.CalendarURL); err != nil {
					if _, err := unit.Put(ctx, "", shared.Unit{
						Name:         i.Name,
						PropertyName: p.Name,
						CalendarURL:  i.CalendarURL,
					}); err != nil {
						errors = append(errors, err)
					}
				}
			}
			if !found {
				errors = append(errors, fmt.Errorf("can't find property for id: %s", i.PropertyID))
			}
		}
		if len(errors) > 0 {
			return Response{}, fmt.Errorf("error(s) inserting items(s): %+v", errors)
		}
		//
		// end code for moving properties from postgres to dynamo
		//
	*/

	return Response{
		Messages: []string{"ok"},
	}, nil
}

func migrateUp(db *sql.DB) error {
	goose.SetTableName("goose_db_version")
	return goose.Up(db, "goosemigrations")
}
