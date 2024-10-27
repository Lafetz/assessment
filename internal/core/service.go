package person

import (
	"context"

	"github.com/google/uuid"
)

type PersonSvc struct {
	repo repository
}

func NewPersonSvc(repo repository) PersonSvc {
	return PersonSvc{
		repo: repo,
	}
}

func (s *PersonSvc) AddPerson(ctx context.Context, id uuid.UUID) error {
	return s.repo.AddPerson(ctx, id)
}

func (s *PersonSvc) GetPerson(ctx context.Context, id uuid.UUID) (Person, error) {
	return s.repo.GetPerson(ctx, id)
}

func (s *PersonSvc) GetPersons(ctx context.Context) ([]Person, error) {
	return s.repo.GetPersons(ctx)
}

func (s *PersonSvc) DeletePerson(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeletePerson(ctx, id)
}

func (s *PersonSvc) UpdatePerson(ctx context.Context, person Person) error {
	return s.repo.UpdatePerson(ctx, person)
}
