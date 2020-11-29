package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
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
	ex, err := os.Executable()
	if err != nil {
		return Response{}, err
	}
	dir := path.Dir(ex)
	if err := godotenv.Load(dir + "/.env"); err != nil {
		return Response{}, errors.New("Error loading .env file")
	}

	db, err := sql.Open("pgx", getConnectionString())
	if err != nil {
		return Response{}, fmt.Errorf("error opening DB: %s", err.Error())
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return Response{}, fmt.Errorf("error pinging DB: %s", err.Error())
	}

	if err := migrateUp(db); err != nil {
		return Response{}, fmt.Errorf("error migrating DB: %s", err.Error())
	}

	dbName := getCurrentDatabase(db)

	return Response{
		Messages: []string{
			"current db",
			dbName,
		},
	}, nil
}

func getConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		getConfig("DB_HOST"),
		5432,
		getConfig("DB_USER"),
		getConfig("DB_PASSWORD"),
		getConfig("DB_NAME"),
	)
}

func getConfig(name string) string {
	val := os.Getenv(name)
	if val == "" {
		fmt.Printf("can't find config for \"%s\"\n", name)
	}
	return val
}

func getCurrentDatabase(db *sql.DB) string {
	row := db.QueryRow("SELECT current_database()")
	if row == nil {
		return "error querying DB"
	}

	var name string
	if err := row.Scan(&name); err != nil {
		return fmt.Sprintf("err: %s", err.Error())
	}

	return name
}

func getDatabases(db *sql.DB) []string {
	rows, err := db.Query("SELECT datname FROM pg_database")
	if err != nil {
		return []string{fmt.Sprintf("err: %s", err.Error())}
	}
	defer rows.Close()

	dbs := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			name = fmt.Sprintf("err: %s", err.Error())
		}
		dbs = append(dbs, name)
	}

	return dbs
}

func migrateUp(db *sql.DB) error {
	// maybe we should do more?
	goose.SetTableName("goose_db_version")
	return goose.Up(db, "goosemigrations")
}
