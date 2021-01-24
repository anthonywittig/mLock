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

type scanner interface {
	Scan(dest ...interface{}) error
}

const (
	table      = "units"
	allColumns = "id, name, property_id, calendar_url, updated_by"
)

func GetByID(ctx context.Context, id string) (shared.UnitX, bool, error) {
	db, err := postgres.GetDB(ctx)
	if err != nil {
		return shared.UnitX{}, false, fmt.Errorf("error getting DB: %s", err.Error())
	}

	// Verify id is a uuid.
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return shared.UnitX{}, false, fmt.Errorf("error parsing ID (%s): %s", id, err.Error())
	}

	row := db.QueryRowContext(ctx, "SELECT "+allColumns+" FROM "+table+" WHERE id = $1", parsedID)
	return getEntity(row)
}

func GetByName(ctx context.Context, name string) (shared.UnitX, bool, error) {
	name = strings.TrimSpace(name)

	db, err := postgres.GetDB(ctx)
	if err != nil {
		return shared.UnitX{}, false, fmt.Errorf("error getting DB: %s", err.Error())
	}

	row := db.QueryRowContext(ctx, "SELECT "+allColumns+" FROM "+table+" WHERE name = $1", name)
	return getEntity(row)
}

func GetAll(ctx context.Context) ([]shared.UnitX, error) {
	db, err := postgres.GetDB(ctx)
	if err != nil {
		return []shared.UnitX{}, fmt.Errorf("error getting DB: %s", err.Error())
	}

	rows, err := db.QueryContext(ctx, "SELECT "+allColumns+" FROM "+table+" ORDER BY name")
	if err != nil {
		return []shared.UnitX{}, fmt.Errorf("error doing query: %s", err.Error())
	}
	defer rows.Close()

	entities := []shared.UnitX{}
	for rows.Next() {
		entity, ok, err := getEntity(rows)
		if err != nil {
			return []shared.UnitX{}, fmt.Errorf("error getting entity: %s", err.Error())
		}
		if !ok {
			return []shared.UnitX{}, fmt.Errorf("error getting entity, not found")
		}

		entities = append(entities, entity)
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

func Insert(ctx context.Context, name string, propertyID uuid.UUID) (shared.UnitX, error) {
	name = strings.TrimSpace(name)

	db, err := postgres.GetDB(ctx)
	if err != nil {
		return shared.UnitX{}, fmt.Errorf("error getting DB: %s", err.Error())
	}

	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return shared.UnitX{}, fmt.Errorf("can't get context data: %s", err.Error())
	}

	currentUser := cd.User
	if currentUser == nil {
		return shared.UnitX{}, fmt.Errorf("no current user")
	}

	_, err = db.ExecContext(
		ctx,
		`INSERT INTO `+table+` (`+allColumns+`) VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(),
		name,
		propertyID,
		"", // Current usecases only insert with an empty calendar url.
		currentUser.Email,
	)
	if err != nil {
		return shared.UnitX{}, fmt.Errorf("error inserting into DB: %s", err.Error())
	}

	unit, ok, err := GetByName(ctx, name)
	if err != nil {
		return shared.UnitX{}, err
	}
	if !ok {
		return shared.UnitX{}, fmt.Errorf("couldn't find unit after insert")
	}

	return unit, nil
}

func Update(ctx context.Context, data shared.UnitX) (shared.UnitX, error) {
	db, err := postgres.GetDB(ctx)
	if err != nil {
		return shared.UnitX{}, fmt.Errorf("error getting DB: %s", err.Error())
	}

	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return shared.UnitX{}, fmt.Errorf("can't get context data: %s", err.Error())
	}

	currentUser := cd.User
	if currentUser == nil {
		return shared.UnitX{}, fmt.Errorf("no current user")
	}
	data.UpdatedBy = currentUser.Email

	_, err = db.ExecContext(
		ctx,
		`UPDATE `+table+` SET name = $1, property_id = $2, calendar_url = $3, updated_by = $4 WHERE id = $5`,
		data.Name,
		data.PropertyID,
		data.CalendarURL,
		data.UpdatedBy,
		data.ID,
	)
	if err != nil {
		return shared.UnitX{}, fmt.Errorf("error inserting into DB: %s", err.Error())
	}

	unit, ok, err := GetByID(ctx, data.ID)
	if err != nil {
		return shared.UnitX{}, err
	}
	if !ok {
		return shared.UnitX{}, fmt.Errorf("couldn't find unit after update")
	}

	return unit, nil
}

func getEntity(s scanner) (shared.UnitX, bool, error) {
	var id string
	var name string
	var propertyID uuid.UUID
	var calendarURL string
	var updatedBy string
	if err := s.Scan(&id, &name, &propertyID, &calendarURL, &updatedBy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shared.UnitX{}, false, nil // Not really an error.
		}

		return shared.UnitX{}, false, fmt.Errorf("error scanning row: %s", err.Error())
	}

	return shared.UnitX{
		ID:          id,
		Name:        name,
		PropertyID:  propertyID,
		CalendarURL: calendarURL,
		UpdatedBy:   updatedBy,
	}, true, nil
}
