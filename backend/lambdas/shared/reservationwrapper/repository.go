package reservationwrapper

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"

	"github.com/google/uuid"
)

type ReservationRepository interface {
	GetForUnits(ctx context.Context, units []shared.Unit) (map[uuid.UUID][]shared.Reservation, error)
}

type Repository struct {
	reservationRepositories []ReservationRepository
}

func NewRepository(reservationRepositories []ReservationRepository) *Repository {
	return &Repository{
		reservationRepositories: reservationRepositories,
	}
}

func (r *Repository) GetForUnits(ctx context.Context, units []shared.Unit) (map[uuid.UUID][]shared.Reservation, error) {
	reservationsByUnit := make(map[uuid.UUID][]shared.Reservation)

	for i, rr := range r.reservationRepositories {
		reservations, err := rr.GetForUnits(ctx, units)
		if err != nil {
			return nil, err
		}

		for unitID, unitReservations := range reservations {
			existingReservations := reservationsByUnit[unitID]
			for _, unitReservation := range unitReservations {
				found := false
				for _, existingReservation := range existingReservations {
					if existingReservation.TransactionNumber == unitReservation.TransactionNumber {
						found = true
						fmt.Printf("Duplicate reservation %s\n", unitReservation.TransactionNumber)
						break
					}
				}

				if !found {
					existingReservations = append(existingReservations, unitReservation)
					if i > 0 {
						fmt.Printf("Missing reservation %s\n", unitReservation.TransactionNumber)
					}
				}
			}
			reservationsByUnit[unitID] = existingReservations
		}
	}

	return reservationsByUnit, nil
}
