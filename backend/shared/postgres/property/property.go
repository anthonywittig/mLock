package property

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mlock/shared"
	"mlock/shared/postgres"
	"strings"

	"github.com/google/uuid"
)

func GetByID(ctx context.Context, id string) (shared.PropertyX, bool, error) {
	db, err := postgres.GetDB(ctx)
	if err != nil {
		return shared.PropertyX{}, false, fmt.Errorf("error getting DB: %s", err.Error())
	}

	// Verify id is a uuid.
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return shared.PropertyX{}, false, fmt.Errorf("error parsing ID: %s", err.Error())
	}

	row := db.QueryRowContext(ctx, "SELECT id, name, created_by FROM properties WHERE id = $1", parsedID)
	var idResult string
	var nameResult string
	var createdByResult string
	if err := row.Scan(&idResult, &nameResult, &createdByResult); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shared.PropertyX{}, false, nil // Not really an error.
		}

		return shared.PropertyX{}, false, fmt.Errorf("error scanning row: %s", err.Error())
	}

	return shared.PropertyX{ID: idResult, Name: nameResult, CreatedBy: createdByResult}, true, nil
}

func GetByName(ctx context.Context, name string) (shared.PropertyX, bool, error) {
	name = strings.TrimSpace(name)

	db, err := postgres.GetDB(ctx)
	if err != nil {
		return shared.PropertyX{}, false, fmt.Errorf("error getting DB: %s", err.Error())
	}

	row := db.QueryRowContext(ctx, "SELECT id, name, created_by FROM properties WHERE name = $1", name)
	var idResult string
	var nameResult string
	var createdByResult string
	if err := row.Scan(&idResult, &nameResult, &createdByResult); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shared.PropertyX{}, false, nil // Not really an error.
		}

		return shared.PropertyX{}, false, fmt.Errorf("error scanning row: %s", err.Error())
	}

	return shared.PropertyX{ID: idResult, Name: nameResult, CreatedBy: createdByResult}, true, nil
}

func GetAll(ctx context.Context) ([]shared.PropertyX, error) {
	db, err := postgres.GetDB(ctx)
	if err != nil {
		return []shared.PropertyX{}, fmt.Errorf("error getting DB: %s", err.Error())
	}

	rows, err := db.QueryContext(ctx, "SELECT id, name, created_by FROM properties ORDER BY name")
	if err != nil {
		return []shared.PropertyX{}, fmt.Errorf("error doing query: %s", err.Error())
	}
	defer rows.Close()

	entities := []shared.PropertyX{}
	for rows.Next() {
		var id string
		var name string
		var createdBy string
		if err := rows.Scan(&id, &name, &createdBy); err != nil {
			return []shared.PropertyX{}, fmt.Errorf("error scanning row: %s", err.Error())
		}
		entities = append(entities, shared.PropertyX{ID: id, Name: name, CreatedBy: createdBy})
	}

	return entities, nil
}

func Delete(ctx context.Context, id string) error {
	db, err := postgres.GetDB(ctx)
	if err != nil {
		return fmt.Errorf("error getting DB: %s", err.Error())
	}

	// No audit trail for deletes. :(

	// Verify id is a uuid.
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("error parsing ID: %s", err.Error())
	}

	_, err = db.ExecContext(
		ctx,
		`DELETE FROM properties WHERE id = $1`,
		parsedID,
	)
	if err != nil {
		return fmt.Errorf("error deleting from DB: %s", err.Error())
	}

	return nil
}

func Insert(ctx context.Context, name string) (shared.PropertyX, error) {
	name = strings.TrimSpace(name)

	db, err := postgres.GetDB(ctx)
	if err != nil {
		return shared.PropertyX{}, fmt.Errorf("error getting DB: %s", err.Error())
	}

	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return shared.PropertyX{}, fmt.Errorf("can't get context data: %s", err.Error())
	}

	currentUser := cd.User
	if currentUser == nil {
		return shared.PropertyX{}, fmt.Errorf("no current user")
	}

	_, err = db.ExecContext(
		ctx,
		`INSERT INTO properties (id, name, created_by) VALUES ($1, $2, $3)`,
		uuid.New(),
		name,
		currentUser.Email,
	)
	if err != nil {
		return shared.PropertyX{}, fmt.Errorf("error inserting into DB: %s", err.Error())
	}

	entity, ok, err := GetByName(ctx, name)
	if err != nil {
		return shared.PropertyX{}, err
	}
	if !ok {
		return shared.PropertyX{}, fmt.Errorf("couldn't find entity after insert")
	}

	return entity, nil
}
