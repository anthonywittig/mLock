package shared

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func GetDB(ctx context.Context) (*sql.DB, error) {
	cd, err := GetContextData(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting context data: %s", err.Error())
	}

	if cd.DB != nil {
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
		GetConfig("DB_HOST"),
		5432,
		GetConfig("DB_USER"),
		GetConfig("DB_PASSWORD"),
		GetConfig("DB_NAME"),
	)
}

func GetCurrentDatabase(ctx context.Context) string {
	db, err := GetDB(ctx)
	if err != nil {
		return fmt.Sprintf("error getting DB: %s", err.Error())
	}

	row := db.QueryRowContext(ctx, "SELECT current_database()")
	if row == nil {
		return "error querying DB"
	}

	var name string
	if err := row.Scan(&name); err != nil {
		return fmt.Sprintf("err: %s", err.Error())
	}

	return name
}

func GetDatabases(ctx context.Context) []string {
	db, err := GetDB(ctx)
	if err != nil {
		return []string{fmt.Sprintf("error getting DB: %s", err.Error())}
	}

	rows, err := db.QueryContext(ctx, "SELECT datname FROM pg_database")
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
