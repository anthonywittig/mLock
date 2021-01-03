package shared

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Property struct {
	ID        string
	Name      string
	CreatedBy string
}

func GetPropertyByID(ctx context.Context, id string) (Property, bool, error) {
	db, err := GetDB(ctx)
	if err != nil {
		return Property{}, false, fmt.Errorf("error getting DB: %s", err.Error())
	}

	// Verify id is a uuid.
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return Property{}, false, fmt.Errorf("error parsing ID: %s", err.Error())
	}

	row := db.QueryRowContext(ctx, "SELECT id, name, created_by FROM properties WHERE id = $1", parsedID)
	var idResult string
	var nameResult string
	var createdByResult string
	if err := row.Scan(&idResult, &nameResult, &createdByResult); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Property{}, false, nil // Not really an error.
		}

		return Property{}, false, fmt.Errorf("error scanning row: %s", err.Error())
	}

	return Property{ID: idResult, Name: nameResult, CreatedBy: createdByResult}, true, nil
}

func GetPropertyByName(ctx context.Context, name string) (Property, bool, error) {
	name = strings.TrimSpace(name)

	db, err := GetDB(ctx)
	if err != nil {
		return Property{}, false, fmt.Errorf("error getting DB: %s", err.Error())
	}

	row := db.QueryRowContext(ctx, "SELECT id, name, created_by FROM properties WHERE name = $1", name)
	var idResult string
	var nameResult string
	var createdByResult string
	if err := row.Scan(&idResult, &nameResult, &createdByResult); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Property{}, false, nil // Not really an error.
		}

		return Property{}, false, fmt.Errorf("error scanning row: %s", err.Error())
	}

	return Property{ID: idResult, Name: nameResult, CreatedBy: createdByResult}, true, nil
}

func GetProperties(ctx context.Context) ([]Property, error) {
	db, err := GetDB(ctx)
	if err != nil {
		return []Property{}, fmt.Errorf("error getting DB: %s", err.Error())
	}

	rows, err := db.QueryContext(ctx, "SELECT id, name, created_by FROM properties ORDER BY name")
	if err != nil {
		return []Property{}, fmt.Errorf("error doing query: %s", err.Error())
	}
	defer rows.Close()

	entities := []Property{}
	for rows.Next() {
		var id string
		var name string
		var createdBy string
		if err := rows.Scan(&id, &name, &createdBy); err != nil {
			return []Property{}, fmt.Errorf("error scanning row: %s", err.Error())
		}
		entities = append(entities, Property{ID: id, Name: name, CreatedBy: createdBy})
	}

	return entities, nil
}

func DeleteProperty(ctx context.Context, id string) error {
	db, err := GetDB(ctx)
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

func InsertProperty(ctx context.Context, name string) (Property, error) {
	name = strings.TrimSpace(name)

	db, err := GetDB(ctx)
	if err != nil {
		return Property{}, fmt.Errorf("error getting DB: %s", err.Error())
	}

	cd, err := GetContextData(ctx)
	if err != nil {
		return Property{}, fmt.Errorf("can't get context data: %s", err.Error())
	}

	currentUser := cd.User
	if currentUser == nil {
		return Property{}, fmt.Errorf("no current user")
	}

	_, err = db.ExecContext(
		ctx,
		`INSERT INTO properties (id, name, created_by) VALUES ($1, $2, $3)`,
		uuid.New(),
		name,
		currentUser.Email,
	)
	if err != nil {
		return Property{}, fmt.Errorf("error inserting into DB: %s", err.Error())
	}

	entity, ok, err := GetPropertyByName(ctx, name)
	if err != nil {
		return Property{}, err
	}
	if !ok {
		return Property{}, fmt.Errorf("couldn't find entity after insert")
	}

	return entity, nil
}
