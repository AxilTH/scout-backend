package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// validate — глобальный экземпляр валидатора (синглтон)
var validate = validator.New()

// GetValidator возвращает экземпляр валидатора
func GetValidator() *validator.Validate {
	return validate
}

// FormatError преобразует ошибки валидации в человекочитаемые сообщения
func FormatError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := e.Field()
			switch e.Tag() {
			case "required":
				if field == "Title" {
					return "Title must not be empty"
				}
				return fmt.Sprintf("%s is required", field)
			case "min":
				if field == "Year" {
					return "Year out of range"
				}
				if field == "Title" {
					return "Title must not be empty"
				}
				if field == "SquadID" {
					return "SquadID is required"
				}
				if field == "CourseID" {
					return "CourseID is required"
				}
				if field == "ActivityTypeID" {
					return "ActivityTypeID is required"
				}
				return fmt.Sprintf("%s is too small", field)
			case "max":
				if field == "Year" {
					return "Year out of range"
				}
				return fmt.Sprintf("%s is too large", field)
			default:
				return fmt.Sprintf("Invalid value for %s", field)
			}
		}
	}
	return "Invalid input data"
}
