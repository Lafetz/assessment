package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/lafetz/assessment/internal/core/domain"
	person "github.com/lafetz/assessment/internal/core/service"
	"github.com/lafetz/assessment/internal/web/dto"
	customvalidator "github.com/lafetz/assessment/internal/web/validation"
)

func AddPerson(personSvc person.PersonSvcApi, logger *slog.Logger, v *customvalidator.CustomValidator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var createPerson dto.CreatePerson
		if err := json.NewDecoder(r.Body).Decode(&createPerson); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		person := domain.NewPerson(createPerson.Name, createPerson.Age, createPerson.Hobbies)
		if _, err := personSvc.AddPerson(r.Context(), person); err != nil {
			logger.Error(err.Error())
			http.Error(w, "Unable to add person", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(dto.ConvertToJSONPerson(person))
	}
}

func GetPersonByID(personSvc person.PersonSvcApi, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		personIDStr := r.PathValue("personId")
		personID, err := uuid.Parse(personIDStr)
		if err != nil {
			http.Error(w, "Invalid person ID", http.StatusUnprocessableEntity)
			return
		}

		person, err := personSvc.GetPerson(r.Context(), personID)
		if err != nil {
			HandleError(err, w, logger)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dto.ConvertToJSONPerson(person))
	}
}

func GetPersons(personSvc person.PersonSvcApi, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 32)
		if err != nil {
			page = 0
		}
		size, err := strconv.ParseInt(r.URL.Query().Get("size"), 10, 32)
		if err != nil {
			size = 10
		}

		persons, metadata, err := personSvc.GetPersons(r.Context(), int32(page), int32(size))
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Unable to get persons", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := dto.ConvertToGetPersonsResponse(persons, metadata)
		json.NewEncoder(w).Encode(response)
	}
}

func UpdatePerson(personSvc person.PersonSvcApi, logger *slog.Logger, v *customvalidator.CustomValidator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		personIDStr := r.PathValue("personId")

		personID, err := uuid.Parse(personIDStr)
		if err != nil {
			http.Error(w, "Invalid person ID", http.StatusUnprocessableEntity)
			return
		}

		var updatePerson dto.UpdatePerson
		if err := json.NewDecoder(r.Body).Decode(&updatePerson); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		person := domain.NewPerson(updatePerson.Name, updatePerson.Age, updatePerson.Hobbies)
		person.ID = personID

		updatedPerson, err := personSvc.UpdatePerson(r.Context(), person)
		if err != nil {
			HandleError(err, w, logger)
			return
		}
		response := dto.ConvertToJSONPerson(updatedPerson)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func DeletePerson(personSvc person.PersonSvcApi, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		personIDStr := r.PathValue("personId")
		personID, err := uuid.Parse(personIDStr)
		if err != nil {
			http.Error(w, "Invalid person ID", http.StatusUnprocessableEntity)
			return
		}
		if err := personSvc.DeletePerson(r.Context(), personID); err != nil {
			HandleError(err, w, logger)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
