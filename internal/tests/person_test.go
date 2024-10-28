package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/lafetz/assessment/internal/core/domain"
	person "github.com/lafetz/assessment/internal/core/service"
	"github.com/lafetz/assessment/internal/repository"
	"github.com/lafetz/assessment/internal/web"
	"github.com/lafetz/assessment/internal/web/dto"
	customvalidator "github.com/lafetz/assessment/internal/web/validation"
	"github.com/stretchr/testify/assert"
)

func TestAddPerson(t *testing.T) {
	repo := repository.NewRepository()
	personSvc := person.NewPersonSvc(repo)
	val := validator.New()
	custonmVal := customvalidator.NewCustomValidator(val)

	web := web.NewApp(8080, slog.Default(), personSvc, custonmVal)

	server := httptest.NewServer(web.Router)
	defer server.Close()

	t.Run("success", func(t *testing.T) {
		payload := `{"name":"John","age":30,"hobbies":["Reading"]}`
		req, _ := http.NewRequest(http.MethodPost, server.URL+"/api/v1/persons", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()
		var createdPerson domain.Person
		err = json.NewDecoder(resp.Body).Decode(&createdPerson)
		assert.NoError(t, err)
		assert.Equal(t, "John", createdPerson.Name)
	})

	t.Run("validation failure", func(t *testing.T) {
		payload := `{"name":"","age":30}`
		req, _ := http.NewRequest(http.MethodPost, server.URL+"/api/v1/persons", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("invalid input", func(t *testing.T) {
		payload := `invalid json`
		req, _ := http.NewRequest(http.MethodPost, server.URL+"/api/v1/persons", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// /
func TestGetPersons(t *testing.T) {
	repo := repository.NewRepository()
	personSvc := person.NewPersonSvc(repo)
	val := validator.New()
	custonmVal := customvalidator.NewCustomValidator(val)

	web := web.NewApp(8080, slog.Default(), personSvc, custonmVal)

	server := httptest.NewServer(web.Router)
	defer server.Close()

	// Add a person to retrieve later
	_, _ = personSvc.AddPerson(context.Background(), domain.Person{Name: "Alice", Age: 28, Hobbies: []string{"Writing"}})

	t.Run("success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/persons?page=0&size=10", nil)
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()
		var personsResponse dto.GetPersonsResponse
		err = json.NewDecoder(resp.Body).Decode(&personsResponse)
		assert.NoError(t, err)
		assert.Len(t, personsResponse.Persons, 1)
	})

	t.Run("internal server error", func(t *testing.T) {

		t.Skip("Implement error handling in person service for test")
	})
}

func TestGetPersonByID(t *testing.T) {
	repo := repository.NewRepository()
	personSvc := person.NewPersonSvc(repo)
	val := validator.New()
	custonmVal := customvalidator.NewCustomValidator(val)

	web := web.NewApp(8080, slog.Default(), personSvc, custonmVal)

	server := httptest.NewServer(web.Router)
	defer server.Close()

	newPerson, _ := personSvc.AddPerson(context.Background(), domain.Person{Name: "Bob", Age: 35, Hobbies: []string{"Gaming"}})

	t.Run("success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/persons/"+newPerson.ID.String(), nil)
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()
		var person domain.Person
		err = json.NewDecoder(resp.Body).Decode(&person)
		assert.NoError(t, err)
		assert.Equal(t, newPerson.Name, person.Name)
	})

	t.Run("person not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/persons/invalid-id", nil)
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("invalid id format", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/persons/invalid-id", nil)
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestUpdatePerson(t *testing.T) {
	repo := repository.NewRepository()
	personSvc := person.NewPersonSvc(repo)
	val := validator.New()
	custonmVal := customvalidator.NewCustomValidator(val)

	web := web.NewApp(8080, slog.Default(), personSvc, custonmVal)

	server := httptest.NewServer(web.Router)
	defer server.Close()

	personToUpdate, _ := personSvc.AddPerson(context.Background(), domain.Person{Name: "Charlie", Age: 40, Hobbies: []string{"Cooking"}})

	t.Run("success", func(t *testing.T) {
		payload := `{"name":"Charlie Updated","age":41,"hobbies":["Cooking","Traveling"]}`
		req, _ := http.NewRequest(http.MethodPut, server.URL+"/api/v1/persons/"+personToUpdate.ID.String(), bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()
		var updatedPerson domain.Person
		err = json.NewDecoder(resp.Body).Decode(&updatedPerson)
		assert.NoError(t, err)
		assert.Equal(t, "Charlie Updated", updatedPerson.Name)
	})

	t.Run("person not found", func(t *testing.T) {
		payload := `{"name":"Non-existent","age":30}`
		req, _ := http.NewRequest(http.MethodPut, server.URL+"/api/v1/persons/invalid-id", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
	t.Run("name missing", func(t *testing.T) {
		payload := `{"name":"","age":30}`
		req, _ := http.NewRequest(http.MethodPut, server.URL+"/api/v1/persons/"+personToUpdate.ID.String(), bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
	t.Run("incorrect age", func(t *testing.T) {
		payload := `{"name":"wrorld","age":-1}`
		req, _ := http.NewRequest(http.MethodPut, server.URL+"/api/v1/persons/"+personToUpdate.ID.String(), bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
	t.Run("invalid input", func(t *testing.T) {
		payload := ``
		req, _ := http.NewRequest(http.MethodPut, server.URL+"/api/v1/persons/"+personToUpdate.ID.String(), bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestDeletePerson(t *testing.T) {
	repo := repository.NewRepository()
	personSvc := person.NewPersonSvc(repo)
	val := validator.New()
	custonmVal := customvalidator.NewCustomValidator(val)

	web := web.NewApp(8080, slog.Default(), personSvc, custonmVal)

	server := httptest.NewServer(web.Router)
	defer server.Close()

	personToDelete, _ := personSvc.AddPerson(context.Background(), domain.Person{Name: "David", Age: 50, Hobbies: []string{"Hiking"}})

	t.Run("success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, server.URL+"/api/v1/persons/"+personToDelete.ID.String(), nil)
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("person not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, server.URL+"/api/v1/persons/invalid-id", nil)
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}
