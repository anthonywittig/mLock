package datastore

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"mlock/shared"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func GetDB(ctx context.Context) (*sql.DB, error) {
	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting context data: %s", err.Error())
	}

	if cd.DB != nil {
		log.Print("already have DB, no need to create another")
		return cd.DB, nil
	}

	db, err := sql.Open("pgx", getConnectionString())
	if err != nil {
		return nil, fmt.Errorf("error opening DB: %s", err.Error())
	}

	cd.DB = db

	return db, nil
}

func getConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		shared.GetConfig("DB_HOST"),
		5432,
		shared.GetConfig("DB_USER"),
		shared.GetConfig("DB_PASSWORD"),
		shared.GetConfig("DB_NAME"),
	)
}

func GetCurrentDatabase(db *sql.DB) string {
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

func GetDatabases(db *sql.DB) []string {
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
