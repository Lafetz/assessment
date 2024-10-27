package person

import (
	"context"

	"github.com/google/uuid"
)

type repository interface {
	AddPerson(ctx context.Context, id uuid.UUID) error
	GetPerson(ctx context.Context, id uuid.UUID) (Person, error)
	GetPersons(ctx context.Context) ([]Person, error)
	DeletePerson(ctx context.Context, id uuid.UUID) error
	UpdatePerson(ctx context.Context, person Person) error
}
