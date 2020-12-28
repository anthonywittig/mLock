package datastore

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

type User struct {
	ID    string
	Email string
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func GetUserByEmail(db *sql.DB, email string) (User, error) {
	if db == nil {
		var err error
		db, err = GetDB()
		if err != nil {
			return User{}, fmt.Errorf("error getting DB: %s", err.Error())
		}
	}

	row := db.QueryRow("SELECT id, email FROM users WHERE email = $1", email)
	var idResult string
	var emailResult string
	if err := row.Scan(&idResult, &emailResult); err != nil {
		return User{}, fmt.Errorf("error scanning row: %s", err.Error())
	}

	// This is redundant as the `Scan` will fail if there aren't any results too.
	if idResult == "" {
		return User{}, fmt.Errorf("error finding user")
	}

	return User{ID: idResult, Email: emailResult}, nil
}

func GetUsers(db *sql.DB) ([]User, error) {
	if db == nil {
		var err error
		db, err = GetDB()
		if err != nil {
			return []User{}, fmt.Errorf("error getting DB: %s", err.Error())
		}
	}

	rows, err := db.Query("SELECT id, email FROM users")
	if err != nil {
		return []User{}, fmt.Errorf("error doing query: %s", err.Error())
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var id string
		var email string
		if err := rows.Scan(&id, &email); err != nil {
			return []User{}, fmt.Errorf("error scanning row: %s", err.Error())
		}
		users = append(users, User{ID: id, Email: email})
	}

	return users, nil
}

func InsertUser(db *sql.DB, email string) error {
	if db == nil {
		var err error
		db, err = GetDB()
		if err != nil {
			return fmt.Errorf("error getting DB: %s", err.Error())
		}
	}

	if !isEmailValid(email) {
		// Should indicate it's a 4xx; we should probably do some validation on the frontend too.
		return fmt.Errorf("email isn't formatted correctly")
	}

	_, err := db.Exec(
		`INSERT INTO users (id, email) VALUES ($1, $2)`,
		uuid.New(),
		email,
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
