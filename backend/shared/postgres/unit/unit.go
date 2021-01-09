package unit

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

const (
	table = "units"
)

func GetByID(ctx context.Context, id string) (shared.Unit, bool, error) {
	db, err := postgres.GetDB(ctx)
	if err != nil {
		return shared.Unit{}, false, fmt.Errorf("error getting DB: %s", err.Error())
	}

	// Verify id is a uuid.
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return shared.Unit{}, false, fmt.Errorf("error parsing ID: %s", err.Error())
	}

	row := db.QueryRowContext(ctx, "SELECT id, name, property_id, updated_by FROM "+table+" WHERE id = $1", parsedID)
	var idResult string
	var nameResult string
	var propertyIDResult string
	var updatedByResult string
	if err := row.Scan(&idResult, &nameResult, &propertyIDResult, &updatedByResult); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shared.Unit{}, false, nil // Not really an error.
		}

		return shared.Unit{}, false, fmt.Errorf("error scanning row: %s", err.Error())
	}

	return shared.Unit{ID: idResult, Name: nameResult, PropertyID: propertyIDResult, UpdatedBy: updatedByResult}, true, nil
}

func GetByName(ctx context.Context, name string) (shared.Unit, bool, error) {
	name = strings.TrimSpace(name)

	db, err := postgres.GetDB(ctx)
	if err != nil {
		return shared.Unit{}, false, fmt.Errorf("error getting DB: %s", err.Error())
	}

	row := db.QueryRowContext(ctx, "SELECT id, name, property_id, updated_by FROM "+table+" WHERE name = $1", name)
	var idResult string
	var nameResult string
	var propertyIDResult string
	var updatedByResult string
	if err := row.Scan(&idResult, &nameResult, &propertyIDResult, &updatedByResult); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shared.Unit{}, false, nil // Not really an error.
		}

		return shared.Unit{}, false, fmt.Errorf("error scanning row: %s", err.Error())
	}

	return shared.Unit{ID: idResult, Name: nameResult, PropertyID: propertyIDResult, UpdatedBy: updatedByResult}, true, nil
}

func GetAll(ctx context.Context) ([]shared.Unit, error) {
	db, err := postgres.GetDB(ctx)
	if err != nil {
		return []shared.Unit{}, fmt.Errorf("error getting DB: %s", err.Error())
	}

	rows, err := db.QueryContext(ctx, "SELECT id, name, property_id, updated_by FROM "+table+" ORDER BY name")
	if err != nil {
		return []shared.Unit{}, fmt.Errorf("error doing query: %s", err.Error())
	}
	defer rows.Close()

	entities := []shared.Unit{}
	for rows.Next() {
		var id string
		var name string
		var propertyID string
		var updatedBy string
		if err := rows.Scan(&id, &name, &propertyID, &updatedBy); err != nil {
			return []shared.Unit{}, fmt.Errorf("error scanning row: %s", err.Error())
		}
		entities = append(entities, shared.Unit{ID: id, Name: name, PropertyID: propertyID, UpdatedBy: updatedBy})
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
		`DELETE FROM `+table+` WHERE id = $1`,
		parsedID,
	)
	if err != nil {
		return fmt.Errorf("error deleting from DB: %s", err.Error())
	}

	return nil
}

func Insert(ctx context.Context, name string, propertyID uuid.UUID) (shared.Unit, error) {
	name = strings.TrimSpace(name)

	db, err := postgres.GetDB(ctx)
	if err != nil {
		return shared.Unit{}, fmt.Errorf("error getting DB: %s", err.Error())
	}

	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return shared.Unit{}, fmt.Errorf("can't get context data: %s", err.Error())
	}

	currentUser := cd.User
	if currentUser == nil {
		return shared.Unit{}, fmt.Errorf("no current user")
	}

	_, err = db.ExecContext(
		ctx,
		`INSERT INTO `+table+` (id, name, property_id, updated_by) VALUES ($1, $2, $3, $4)`,
		uuid.New(),
		name,
		propertyID,
		currentUser.Email,
	)
	if err != nil {
		return shared.Unit{}, fmt.Errorf("error inserting into DB: %s", err.Error())
	}

	unit, ok, err := GetByName(ctx, name)
	if err != nil {
		return shared.Unit{}, err
	}
	if !ok {
		return shared.Unit{}, fmt.Errorf("couldn't find unit after insert")
	}

	return unit, nil
}
