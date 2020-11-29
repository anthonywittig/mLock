package datastore

import (
	"fmt"
)

type User struct {
	ID    string
	Email string
}

func GetUsers() ([]User, error) {
	db, err := GetDB()
	if err != nil {
		return []User{}, fmt.Errorf("error getting DB: %s", err.Error())
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
