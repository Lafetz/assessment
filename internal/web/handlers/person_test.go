package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/lafetz/assessment/internal/core/domain"
	"github.com/lafetz/assessment/internal/web/dto"
	"github.com/lafetz/assessment/internal/web/handlers"
	customvalidator "github.com/lafetz/assessment/internal/web/validation"
)

type MockPersonSvc struct {
}

func NewMockPersonSvc() *MockPersonSvc {
	return &MockPersonSvc{}
}

func (m *MockPersonSvc) AddPerson(ctx context.Context, person domain.Person) (domain.Person, error) {

	person.ID = uuid.New()
	return person, nil
}

func (m *MockPersonSvc) GetPerson(ctx context.Context, id uuid.UUID) (domain.Person, error) {

	return domain.Person{
		ID:      id,
		Name:    "Test Person",
		Age:     30,
		Hobbies: []string{"Reading", "Gaming"},
	}, nil
}

func (m *MockPersonSvc) GetPersons(ctx context.Context, page, size int32) ([]domain.Person, domain.Metadata, error) {
	persons := []domain.Person{
		{ID: uuid.New(), Name: "Alice", Age: 25, Hobbies: []string{"Dancing"}},
		{ID: uuid.New(), Name: "Bob", Age: 28, Hobbies: []string{"Cycling"}},
	}

	metadata := domain.Metadata{
		TotalRecords: int32(len(persons)),
		CurrentPage:  page,
		LastPage:     int32(int32(len(persons)) / size),
	}

	return persons, metadata, nil
}
func (m *MockPersonSvc) DeletePerson(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *MockPersonSvc) UpdatePerson(ctx context.Context, person domain.Person) (domain.Person, error) {
	return person, nil
}
func TestAddPerson(t *testing.T) {
	mockSvc := NewMockPersonSvc()
	handler := handlers.AddPerson(mockSvc, slog.Default(), customvalidator.NewCustomValidator(validator.New()))

	createPerson := dto.CreatePerson{
		Name:    "John Doe",
		Age:     30,
		Hobbies: []string{"Reading", "Swimming"},
	}

	reqBody, _ := json.Marshal(createPerson)
	req := httptest.NewRequest(http.MethodPost, "/persons", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}
	fmt.Printf("%s", w.Body)
	var response dto.JSONPerson
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}
	if response.Name != createPerson.Name {
		t.Errorf("Expected person name %s, got %s", createPerson.Name, response.Name)
	}
}

func TestGetPersonByID(t *testing.T) {
	mockSvc := NewMockPersonSvc()
	handler := handlers.GetPersonByID(mockSvc, slog.Default())

	personID := uuid.New()
	expectedPerson := domain.Person{
		ID:      personID,
		Name:    "Test Person",
		Age:     30,
		Hobbies: []string{"Reading", "Gaming"},
	}

	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID.String(), nil)
	req.SetPathValue("personId", personID.String())
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response dto.JSONPerson
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}
	if response.Name != expectedPerson.Name {
		t.Errorf("Expected person name %s, got %s", expectedPerson.Name, response.Name)
	}
}

func TestGetPersons(t *testing.T) {
	mockSvc := NewMockPersonSvc()
	handler := handlers.GetPersons(mockSvc, slog.Default())

	req := httptest.NewRequest(http.MethodGet, "/persons?page=0&size=10", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response dto.GetPersonsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}
	if len(response.Persons) != 2 {
		t.Errorf("Expected 2 persons, got %d", len(response.Persons))
	}
}

func TestUpdatePerson(t *testing.T) {
	mockSvc := NewMockPersonSvc()
	handler := handlers.UpdatePerson(mockSvc, slog.Default(), customvalidator.NewCustomValidator(validator.New()))

	personID := uuid.New()
	updatePerson := dto.UpdatePerson{
		Name:    "John Updated",
		Age:     31,
		Hobbies: []string{"Reading", "Jogging"},
	}
	person := domain.NewPerson(updatePerson.Name, updatePerson.Age, updatePerson.Hobbies)
	person.ID = personID

	reqBody, _ := json.Marshal(updatePerson)
	req := httptest.NewRequest(http.MethodPut, "/persons?personId="+personID.String(), bytes.NewBuffer(reqBody))
	req.SetPathValue("personId", personID.String())
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response dto.JSONPerson
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}
	if response.Name != updatePerson.Name {
		t.Errorf("Expected updated person name %s, got %s", updatePerson.Name, response.Name)
	}
}

func TestDeletePerson(t *testing.T) {
	mockSvc := NewMockPersonSvc()
	handler := handlers.DeletePerson(mockSvc, slog.Default())

	personID := uuid.New()

	req := httptest.NewRequest(http.MethodDelete, "/persons?personId="+personID.String(), nil)
	req.SetPathValue("personId", personID.String())
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
	}
}
