package shared

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Type      string // deprecated
	Email     string
	CreatedBy string // deprecated in favor of UpdatedBy
	UpdatedBy string
}

type UserService interface {
	GetByEmail(ctx context.Context, email string) (User, bool, error)
}
