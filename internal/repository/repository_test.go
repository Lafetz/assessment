package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lafetz/assessment/internal/core/domain"
	person "github.com/lafetz/assessment/internal/core/service"
	"github.com/stretchr/testify/assert"
)

func TestAddPerson(t *testing.T) {
	repo := NewRepository()
	person := domain.NewPerson("John D", 30, []string{"Reading", "Swimming"})
	err := repo.AddPerson(context.Background(), person)
	assert.NoError(t, err, "expected no error when adding a person")

	_, err = repo.GetPerson(context.Background(), person.ID)
	assert.NoError(t, err, "expected to get person, got error")
}

func TestAddPerson_AlreadyExists(t *testing.T) {
	repo := NewRepository()
	personID := uuid.New()
	person := domain.Person{ID: personID, Name: "John Doe", Age: 30, Hobbies: []string{"Reading", "Swimming"}}

	err := repo.AddPerson(context.Background(), person)
	assert.NoError(t, err, "expected no error when adding a person")

	err = repo.AddPerson(context.Background(), person)
	assert.Equal(t, ErrDuplicatePk, err, "expected duplicate primary key error")
}

func TestGetPerson(t *testing.T) {
	repo := NewRepository()
	personID := uuid.New()
	person := domain.Person{ID: personID, Name: "Jane Doe", Age: 25, Hobbies: []string{"Cycling"}}

	err := repo.AddPerson(context.Background(), person)
	assert.NoError(t, err, "expected no error when adding a person")

	retrievedPerson, err := repo.GetPerson(context.Background(), personID)
	assert.NoError(t, err, "expected to get person, got error")
	assert.Equal(t, person.ID, retrievedPerson.ID, "expected retrieved person's ID to match the original")
}

func TestGetPerson_NotFound(t *testing.T) {
	repo := NewRepository()
	personID := uuid.New()

	_, err := repo.GetPerson(context.Background(), personID)
	assert.Error(t, err, person.ErrNotFound, "expected error when getting a non-existent person")
}

func TestUpdatePerson(t *testing.T) {
	repo := NewRepository()
	personID := uuid.New()
	person := domain.Person{ID: personID, Name: "Alice", Age: 28, Hobbies: []string{"Traveling"}}

	err := repo.AddPerson(context.Background(), person)
	assert.NoError(t, err, "expected no error when adding a person")

	person.Age = 29
	err = repo.UpdatePerson(context.Background(), person)
	assert.NoError(t, err, "expected no error when updating a person")

	updatedPerson, err := repo.GetPerson(context.Background(), personID)
	assert.NoError(t, err, "expected to get person, got error")
	assert.Equal(t, int32(29), updatedPerson.Age, "expected age to be 29")
}

func TestUpdatePerson_NotFound(t *testing.T) {
	repo := NewRepository()
	personID := uuid.New()
	p := domain.Person{ID: personID, Name: "Charlie", Age: 32, Hobbies: []string{"Gaming"}}

	err := repo.UpdatePerson(context.Background(), p)
	assert.Error(t, err, person.ErrNotFound, "expected error when updating a non-existent person")
}

func TestDeletePerson(t *testing.T) {
	repo := NewRepository()
	personID := uuid.New()
	p := domain.Person{ID: personID, Name: "Bob", Age: 40, Hobbies: []string{"Cooking"}}

	err := repo.AddPerson(context.Background(), p)
	assert.NoError(t, err, "expected no error when adding a person")

	err = repo.DeletePerson(context.Background(), personID)
	assert.NoError(t, err, "expected no error when deleting a person")

	err = repo.DeletePerson(context.Background(), personID)
	assert.ErrorIs(t, err, person.ErrNotFound, "expected error when getting a deleted person")
}

func TestDeletePerson_NotFound(t *testing.T) {
	repo := NewRepository()
	personID := uuid.New()

	err := repo.DeletePerson(context.Background(), personID)
	assert.ErrorIs(t, err, person.ErrNotFound, "expected error when deleting a non-existent person")
}

func TestGetPersons(t *testing.T) {
	repo := NewRepository()
	persons := []domain.Person{
		{ID: uuid.New(), Name: "Anna", Age: 22, Hobbies: []string{"Dancing"}},
		{ID: uuid.New(), Name: "Tom", Age: 35, Hobbies: []string{"Writing"}},
		{ID: uuid.New(), Name: "Sam", Age: 27, Hobbies: []string{"Biking"}},
		{ID: uuid.New(), Name: "Alice", Age: 30, Hobbies: []string{"Traveling"}},
		{ID: uuid.New(), Name: "Bob", Age: 40, Hobbies: []string{"Cooking"}},
	}

	for _, person := range persons {
		err := repo.AddPerson(context.Background(), person)
		assert.NoError(t, err, "expected no error when adding a person")
	}

	tests := []struct {
		name                 string
		page                 int32
		size                 int32
		expectedLen          int
		expectedPage         int32
		expectedTotalRecords int32
		expectedLastPage     int32
	}{
		{"Page 1, Size 2", 0, 2, 2, 1, 5, 3},
		{"Page 2, Size 2", 1, 2, 2, 2, 5, 3},
		{"Page 3, Size 2", 2, 2, 1, 3, 5, 3},
		{"Page 1, Size 3", 0, 3, 3, 1, 5, 2},
		{"Page 2, Size 3", 1, 3, 2, 2, 5, 2},
		{"Page 1, Size 5", 0, 5, 5, 1, 5, 1},
		{"Page 1, Size 10", 0, 10, 5, 1, 5, 1},
		{"Page 4, Size 2", 3, 2, 0, 4, 5, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retrievedPersons, metadata, err := repo.GetPersons(context.Background(), tt.page, tt.size)
			assert.NoError(t, err, "expected no error when getting persons")
			assert.Len(t, retrievedPersons, tt.expectedLen, "expected length of retrieved persons to match")
			assert.Equal(t, tt.expectedTotalRecords, int32(metadata.TotalRecords), "expected total records to match the added persons count")
			assert.Equal(t, tt.expectedPage, metadata.CurrentPage, "expected current page to match")
			assert.Equal(t, tt.expectedLastPage, metadata.LastPage, "expected last page to match")
		})
	}
}
