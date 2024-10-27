package person

import (
	"context"

	"github.com/google/uuid"
	"github.com/lafetz/assessment/internal/core/domain"
)

type repository interface {
	AddPerson(ctx context.Context, person domain.Person) error
	GetPerson(ctx context.Context, id uuid.UUID) (domain.Person, error)
	GetPersons(ctx context.Context, page, size int32) ([]domain.Person, domain.Metadata, error)
	DeletePerson(ctx context.Context, id uuid.UUID) error
	UpdatePerson(ctx context.Context, person domain.Person) error
}
