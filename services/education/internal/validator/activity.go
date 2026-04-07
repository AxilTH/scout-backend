package validator

import (
	"fmt"
	"time"
)

// CreateActivityRequest — структура ТОЛЬКО для парсинга JSON (без валидации)
type CreateActivityRequest struct {
	CourseID       int64     `json:"course_id"`
	ActivityTypeID int64     `json:"activity_type_id"`
	Title          string    `json:"title"`
	Description    *string   `json:"description"`
	StartsAt       time.Time `json:"starts_at"`
	EndsAt         time.Time `json:"ends_at"`
	MaxScore       *int      `json:"max_score"`
	IsPublished    bool      `json:"is_published"`
}

// CreateActivityInput — структура ДЛЯ валидации (с тегами)
type CreateActivityInput struct {
	CourseID       int64     `validate:"required,min=1"`
	ActivityTypeID int64     `validate:"required,min=1"`
	Title          string    `validate:"required,min=1,max=255"`
	Description    *string   `validate:"omitempty"`
	StartsAt       time.Time `validate:"required"`
	EndsAt         time.Time `validate:"required"`
	MaxScore       *int      `validate:"omitempty,min=0"`
	IsPublished    bool      `validate:"-"`
}

// Validate выполняет кастомную бизнес-валидацию
func (c *CreateActivityInput) Validate() error {
	if !c.EndsAt.After(c.StartsAt) {
		return fmt.Errorf("ends_at must be after starts_at")
	}
	return nil
}
