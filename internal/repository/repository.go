package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/lafetz/assessment/internal/core/domain"
	person "github.com/lafetz/assessment/internal/core/service"
)

var (
	ErrDuplicatePk = errors.New("duplicate pk id")
)

type Repository struct {
	mu      sync.RWMutex
	storage map[uuid.UUID]domain.Person
}

func NewRepository() *Repository {
	return &Repository{
		storage: make(map[uuid.UUID]domain.Person),
	}
}

func (r *Repository) AddPerson(ctx context.Context, person domain.Person) (domain.Person, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.storage[person.ID]; exists {
		return domain.Person{}, ErrDuplicatePk
	}

	r.storage[person.ID] = person
	return domain.Person{}, nil
}

func (r *Repository) GetPerson(ctx context.Context, id uuid.UUID) (domain.Person, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.storage[id]
	if !exists {
		return domain.Person{}, person.ErrNotFound
	}
	return p, nil
}

func (r *Repository) GetPersons(ctx context.Context, page, size int32) ([]domain.Person, domain.Metadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	persons := make([]domain.Person, 0, len(r.storage))
	for _, person := range r.storage {
		persons = append(persons, person)
	}

	totalRecords := int32(len(persons))
	offset := page * size
	limit := size

	if offset > totalRecords {
		return []domain.Person{}, domain.CalculateMetadata(totalRecords, offset, limit), nil
	}

	end := offset + limit
	if end > totalRecords {
		end = totalRecords
	}

	return persons[offset:end], domain.CalculateMetadata(totalRecords, offset, limit), nil
}

func (r *Repository) DeletePerson(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.storage[id]; !exists {
		return person.ErrNotFound
	}

	delete(r.storage, id)
	return nil
}

func (r *Repository) UpdatePerson(ctx context.Context, p domain.Person) (domain.Person, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.storage[p.ID]; !exists {
		return domain.Person{}, person.ErrNotFound
	}

	r.storage[p.ID] = p
	return p, nil
}
