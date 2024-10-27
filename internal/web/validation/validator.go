package customvalidator

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validate *validator.Validate
}

func NewCustomValidator(validate *validator.Validate) *CustomValidator {
	return &CustomValidator{
		validate: validate,
	}
}

func (v *CustomValidator) ValidateAndRespond(w http.ResponseWriter, input interface{}) bool {
	err := v.validate.Struct(input)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := ValidateModel(validationErrors)
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ValidationErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Errors:     errors,
			})
			return true
		}
	}
	return false
}

type ValidationErrorResponse struct {
	StatusCode int         `json:"statusCode"`
	Errors     interface{} `json:"errors"`
}
