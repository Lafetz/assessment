package person

import (
	"context"

	"github.com/google/uuid"
	"github.com/lafetz/assessment/internal/core/domain"
)

type PersonSvc struct {
	repo repository
}

func NewPersonSvc(repo repository) PersonSvc {
	return PersonSvc{
		repo: repo,
	}
}

func (s *PersonSvc) AddPerson(ctx context.Context, person domain.Person) error {
	return s.repo.AddPerson(ctx, person)
}

func (s *PersonSvc) GetPerson(ctx context.Context, id uuid.UUID) (domain.Person, error) {
	return s.repo.GetPerson(ctx, id)
}

func (s *PersonSvc) GetPersons(ctx context.Context, page, size int32) ([]domain.Person, domain.Metadata, error) {
	return s.repo.GetPersons(ctx, page, size)
}

func (s *PersonSvc) DeletePerson(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeletePerson(ctx, id)
}

func (s *PersonSvc) UpdatePerson(ctx context.Context, person domain.Person) error {
	return s.repo.UpdatePerson(ctx, person)
}
