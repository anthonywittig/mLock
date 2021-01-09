package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mlock/shared"
	"mlock/shared/postgres"
	"regexp"

	"github.com/google/uuid"
)

type UserServiceImpl struct{}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{}
}

func GetByID(ctx context.Context, id string) (shared.User, bool, error) {
	db, err := postgres.GetDB(ctx)
	if err != nil {
		return shared.User{}, false, fmt.Errorf("error getting DB: %s", err.Error())
	}

	// Verify id is a uuid.
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return shared.User{}, false, fmt.Errorf("error parsing ID: %s", err.Error())
	}

	row := db.QueryRowContext(ctx, "SELECT id, email, created_by FROM users WHERE id = $1", parsedID)
	var idResult string
	var emailResult string
	var createdByResult string
	if err := row.Scan(&idResult, &emailResult, &createdByResult); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shared.User{}, false, nil // Not really an error.
		}

		return shared.User{}, false, fmt.Errorf("error scanning row: %s", err.Error())
	}

	return shared.User{ID: idResult, Email: emailResult, CreatedBy: createdByResult}, true, nil
}

func (u *UserServiceImpl) GetByEmail(ctx context.Context, email string) (shared.User, bool, error) {
	db, err := postgres.GetDB(ctx)
	if err != nil {
		return shared.User{}, false, fmt.Errorf("error getting DB: %s", err.Error())
	}

	row := db.QueryRowContext(ctx, "SELECT id, email, created_by FROM users WHERE email = $1", email)
	var idResult string
	var emailResult string
	var createdByResult string
	if err := row.Scan(&idResult, &emailResult, &createdByResult); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shared.User{}, false, nil // Not really an error.
		}

		return shared.User{}, false, fmt.Errorf("error scanning row: %s", err.Error())
	}

	return shared.User{ID: idResult, Email: emailResult, CreatedBy: createdByResult}, true, nil
}

func GetAll(ctx context.Context) ([]shared.User, error) {
	db, err := postgres.GetDB(ctx)
	if err != nil {
		return []shared.User{}, fmt.Errorf("error getting DB: %s", err.Error())
	}

	rows, err := db.QueryContext(ctx, "SELECT id, email, created_by FROM users ORDER BY email")
	if err != nil {
		return []shared.User{}, fmt.Errorf("error doing query: %s", err.Error())
	}
	defer rows.Close()

	users := []shared.User{}
	for rows.Next() {
		var id string
		var email string
		var createdBy string
		if err := rows.Scan(&id, &email, &createdBy); err != nil {
			return []shared.User{}, fmt.Errorf("error scanning row: %s", err.Error())
		}
		users = append(users, shared.User{ID: id, Email: email, CreatedBy: createdBy})
	}

	return users, nil
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
		`DELETE FROM users WHERE id = $1`,
		parsedID,
	)
	if err != nil {
		return fmt.Errorf("error deleting from DB: %s", err.Error())
	}

	return nil
}

func Insert(ctx context.Context, email string) error {
	db, err := postgres.GetDB(ctx)
	if err != nil {
		return fmt.Errorf("error getting DB: %s", err.Error())
	}

	if !isEmailValid(email) {
		// Should indicate it's a 4xx; we should probably do some validation on the frontend too.
		return fmt.Errorf("email isn't formatted correctly")
	}

	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return fmt.Errorf("can't get context data: %s", err.Error())
	}

	currentUser := cd.User
	if currentUser == nil {
		return fmt.Errorf("no current user")
	}

	_, err = db.ExecContext(
		ctx,
		`INSERT INTO users (id, email, created_by) VALUES ($1, $2, $3)`,
		uuid.New(),
		email,
		currentUser.Email,
	)
	if err != nil {
		return fmt.Errorf("error inserting into DB: %s", err.Error())
	}

	return nil
}

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}
