package shared

import "context"

type User struct {
	Type      string
	Email     string
	CreatedBy string
}

type UserService interface {
	Get(ctx context.Context, email string) (User, bool, error)
}
