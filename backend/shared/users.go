package shared

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

type User struct {
	ID        string
	Email     string
	CreatedBy string
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func GetUserByEmail(ctx context.Context, email string) (User, bool, error) {
	db, err := GetDB(ctx)
	if err != nil {
		return User{}, false, fmt.Errorf("error getting DB: %s", err.Error())
	}

	row := db.QueryRowContext(ctx, "SELECT id, email, created_by FROM users WHERE email = $1", email)
	var idResult string
	var emailResult string
	var createdByResult string
	if err := row.Scan(&idResult, &emailResult, &createdByResult); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, false, nil // Not really an error.
		}

		return User{}, false, fmt.Errorf("error scanning row: %s", err.Error())
	}

	return User{ID: idResult, Email: emailResult, CreatedBy: createdByResult}, true, nil
}

func GetUsers(ctx context.Context) ([]User, error) {
	db, err := GetDB(ctx)
	if err != nil {
		return []User{}, fmt.Errorf("error getting DB: %s", err.Error())
	}

	rows, err := db.QueryContext(ctx, "SELECT id, email, created_by FROM users ORDER BY email")
	if err != nil {
		return []User{}, fmt.Errorf("error doing query: %s", err.Error())
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var id string
		var email string
		var createdBy string
		if err := rows.Scan(&id, &email, &createdBy); err != nil {
			return []User{}, fmt.Errorf("error scanning row: %s", err.Error())
		}
		users = append(users, User{ID: id, Email: email, CreatedBy: createdBy})
	}

	return users, nil
}

func InsertUser(ctx context.Context, email string) error {
	db, err := GetDB(ctx)
	if err != nil {
		return fmt.Errorf("error getting DB: %s", err.Error())
	}

	if !isEmailValid(email) {
		// Should indicate it's a 4xx; we should probably do some validation on the frontend too.
		return fmt.Errorf("email isn't formatted correctly")
	}

	cd, err := GetContextData(ctx)
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
