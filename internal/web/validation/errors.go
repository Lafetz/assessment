package customvalidator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateModel(err validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)

	for _, err := range err {

		errors[strings.ToLower(err.Field())] = errorMsgs(err.Tag(), err.Param())

	}
	return errors

}

func errorMsgs(tag string, value string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "lte":
		return "can not be greater than " + value
	case "gte":
		return "can not be less than " + value
	case "min":
		return "can not be less than " + value
	case "len":
		return "length invalid"
	case "oneof":
		return "sort value must be one of: " + value
	}
	return ""
}
