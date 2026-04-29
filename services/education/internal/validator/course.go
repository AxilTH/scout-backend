// internal/validator/course.go
package validator

import (
	"errors"
	"fmt"
	"time"
)

// CreateCourseRequest — структура только для десериализации JSON
type CreateCourseRequest struct {
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	Year        int       `json:"year"`
}

// CreateCourseInput — структура для структурной валидации (теги `validate`)
type CreateCourseInput struct {
	Title       string    `validate:"required,min=1,max=255"`
	Description *string   `validate:"omitempty"`
	StartsAt    time.Time `validate:"required"`
	EndsAt      time.Time `validate:"required"`
	Year        int       `validate:"required"`
	SquadID     int64     `validate:"required,min=1"`
}

// Validate выполняет бизнес-валидацию и возвращает ВСЕ ошибки.
func (c *CreateCourseInput) Validate() ValidationErrors {
	var errs ValidationErrors

	// 1. Проверка временных зон
	if c.StartsAt.Location() != time.UTC {
		errs = append(errs, errors.New("starts at must be in UTC"))
	}
	if c.EndsAt.Location() != time.UTC {
		errs = append(errs, errors.New("ends at must be in UTC"))
	}

	// 2. Логические отношения
	if !c.EndsAt.After(c.StartsAt) {
		errs = append(errs, errors.New("ends at must be after starts at"))
	}

	// 3. Согласованность года и starts_at
	if c.Year != c.StartsAt.Year() {
		errs = append(errs, fmt.Errorf("year must match the year of starts at (%d)", c.StartsAt.Year()))
	}

	// 4. Динамическая проверка диапазона года
	currentYear := time.Now().UTC().Year()
	if c.Year < 2020 || c.Year > currentYear + 5 {
		errs = append(errs, fmt.Errorf("year must be between 2020 and %d", currentYear + 5))
	}

	return errs
}
