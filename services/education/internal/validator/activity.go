// internal/validator/activity.go
package validator

import (
	"errors"
	"time"
)

// CreateActivityRequest — структура только для десериализации JSON
type CreateActivityRequest struct {
	Title          string    `json:"title"`
	Description    *string   `json:"description"`
	StartsAt       time.Time `json:"starts_at"`
	EndsAt         time.Time `json:"ends_at"`
	MaxScore       int       `json:"max_score"`
	IsPublished    bool      `json:"is_published"`
	ActivityTypeID int64     `json:"activity_type_id"`
}

// CreateActivityInput — структура для структурной валидации (теги `validate`)
type CreateActivityInput struct {
	Title          string    `validate:"required,min=1,max=255"`
	Description    *string   `validate:"omitempty"`
	StartsAt       time.Time `validate:"required"`
	EndsAt         time.Time `validate:"required"`
	MaxScore       int       `validate:"min=0"`
	IsPublished    bool      `validate:"-"`
	CourseID       int64     `validate:"required,min=1"`
	ActivityTypeID int64     `validate:"required,min=1"`
}

// Validate выполняет бизнес-валидацию и возвращает ВСЕ ошибки.
func (c *CreateActivityInput) Validate() ValidationErrors {
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

	return errs
}