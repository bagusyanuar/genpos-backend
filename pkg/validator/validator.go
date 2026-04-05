package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// RegisterTagNameFunc uses the JSON tag as the field name in error messages.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Validate checks the struct and returns a map of errors in Laravel style (map[string][]string).
func Validate(s any) map[string][]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string][]string)
	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		// If JSON tag is not present, library falls back to struct field name.
		errors[field] = append(errors[field], msgForTag(err.Tag(), field, err.Param()))
	}

	return errors
}

// msgForTag returns a human-readable message based on the validation tag.
func msgForTag(tag, field, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("The %s field is required.", field)
	case "email":
		return fmt.Sprintf("The %s must be a valid email address.", field)
	case "min":
		return fmt.Sprintf("The %s must be at least %s characters.", field, param)
	case "max":
		return fmt.Sprintf("The %s may not be greater than %s characters.", field, param)
	case "unique":
		return fmt.Sprintf("The %s has already been taken.", field)
	}
	return fmt.Sprintf("The %s field is invalid.", field)
}
