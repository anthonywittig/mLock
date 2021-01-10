package shared

import "github.com/google/uuid"

type Unit struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	PropertyID uuid.UUID `json:"propertyId"`
	UpdatedBy  string    `json:"updatedBy"`
}

/*
type UnitService interface {
	getByID(ctx context.Context, id string) (Unit, bool, error)
	GetByName(ctx context.Context, name string) (Unit, bool, error)
	GetAll(ctx context.Context) ([]Unit, error)
	Delete(ctx context.Context, id string) error
	Insert(ctx context.Context, name string, propertyID uuid.UUID) (Unit, error)
}
*/
