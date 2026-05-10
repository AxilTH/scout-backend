// internal/validator/validator.go
package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ValidationErrors представляет список ошибок валидации
type ValidationErrors []error

// Error реализует интерфейс error (возвращает первую ошибку)
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	return ve[0].Error()
}

// Глобальный экземпляр структурного валидатора
var validate = validator.New()

// toHumanField преобразует имя поля Go в человекочитаемую форму.
func toHumanField(field string) string {
	switch field {
	case "SquadID":
		return "squad ID"
	case "UserID":
		return "user ID"
	case "RoleID":
		return "role ID"
	case "PositionID":
		return "position ID"
	case "Title":
		return "title"
	case "IsActive":
		return "is active"
	case "LeftAt":
		return "left at"
	case "DismissedAt":
		return "dismissed at"
	default:
		// fallback: lowercase first letter (простой вариант)
		if len(field) == 0 {
			return field
		}
		return string(field[0]|32) + field[1:] // 'A' → 'a'
	}
}

// ToHumanReadable возвращает все ошибки в виде человекочитаемых строк
func (ve ValidationErrors) ToHumanReadable() []string {
	messages := make([]string, len(ve))
	for i, err := range ve {
		messages[i] = err.Error()
	}
	return messages
}

// ValidateStruct выполняет структурную валидацию (по тегам `validate:"..."`)
// и возвращает все найденные ошибки.
func ValidateStruct(v interface{}) ValidationErrors {
	err := validate.Struct(v)
	if err == nil {
		return nil
	}

	var errs ValidationErrors
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			fieldName := toHumanField(fieldErr.Field())
			tag := fieldErr.Tag()
			switch tag {
			case "required":
				errs = append(errs, fmt.Errorf("%s is required", fieldName))
			case "min", "max":
				errs = append(errs, fmt.Errorf("%s is out of range", fieldName))
			default:
				errs = append(errs, fmt.Errorf("invalid value for %s", fieldName))
			}
		}
	} else {
		errs = append(errs, err)
	}
	return errs
}