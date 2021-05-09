package shared

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"-"` // deprecated
	Email     string    `json:"email"`
	CreatedBy string    `json:"-"` // deprecated in favor of UpdatedBy
	UpdatedBy string    `json:"updatedBy"`
}

type UserService interface {
	GetByEmail(ctx context.Context, email string) (User, bool, error)
}
