package shared

import "context"

type User struct {
	ID        string
	Email     string
	CreatedBy string
}

type UserService interface {
	GetByEmail(ctx context.Context, email string) (User, bool, error)
}
