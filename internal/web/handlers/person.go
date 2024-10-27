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

// AddPerson godoc
//
//	@Summary		Add a new person
//	@Description	Add a new person to the database
//	@Tags			Persons
//	@Accept			json
//	@Produce		json
//	@Param			person	body		dto.CreatePerson	true	"Person data"
//	@Success		201		{object}	domain.Person
//	@Failure		400		{object}	string	"Invalid input"
//	@Router			/api/v1/persons [post]
func AddPerson(personSvc person.PersonSvcApi, logger *slog.Logger, v *customvalidator.CustomValidator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var createPerson dto.CreatePerson
		if err := json.NewDecoder(r.Body).Decode(&createPerson); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		if v.ValidateAndRespond(w, createPerson) {
			return
		}
		person := domain.NewPerson(createPerson.Name, createPerson.Age, createPerson.Hobbies)
		if _, err := personSvc.AddPerson(r.Context(), person); err != nil {
			HandleError(err, w, logger)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(dto.ConvertToJSONPerson(person)); err != nil {
			HandleError(err, w, logger)
		}
	}
}

// GetPersonByID godoc
// @Summary		Get person by ID
// @Description	Retrieve a person by their ID
// @Tags			Persons
// @Accept			json
// @Produce		json
// @Param			personId	path		string	true	"ID of the person"
// @Success		200			{object}	domain.Person
// @Failure		404			{object}	string	"Person not found"
// @Router			/api/v1/persons/{personId} [get]
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(dto.ConvertToJSONPerson(person)); err != nil {
			HandleError(err, w, logger)
		}
	}
}

// GetPersons godoc
//
//	@Summary		Get all persons
//	@Description	Retrieve a list of persons
//	@Tags			Persons
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	domain.Person
//	@Router			/api/v1/persons [get]
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
		if err := json.NewEncoder(w).Encode(response); err != nil {
			HandleError(err, w, logger)
		}
	}
}

// UpdatePerson godoc

// @Summary		Update an existing person
// @Description	Update a person by their ID
// @Tags			Persons
// @Accept			json
// @Produce		json
// @Param			personId	path		string			true	"ID of the person"
// @Param			person		body		dto.CreatePerson	true	"Updated person data"
// @Success		200			{object}	domain.Person
// @Failure		404			{object}	string	"Person not found"
// @Failure		400			{object}	string	"Invalid input"
// @Router			/api/v1/persons/{personId} [put]
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
		if v.ValidateAndRespond(w, updatePerson) {
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			HandleError(err, w, logger)
		}
	}
}

// DeletePerson godoc
//
//	@Summary		Delete a person
//	@Description	Delete a person by their ID
//	@Tags			Persons
//	@Accept			json
//	@Produce		json
//	@Param			personId	path	string	true	"ID of the person"
//	@Success		204			"No Content"
//	@Failure		404			{object}	string	"Person not found"
//	@Router			/api/v1/persons/{personId} [delete]
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
