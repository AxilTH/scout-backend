package validator

import (
	"fmt"
	"time"
)

// CreateCourseRequest — структура ТОЛЬКО для парсинга JSON (без валидации)
type CreateCourseRequest struct {
	SquadID     int64     `json:"squad_id"`
	Year        int       `json:"year"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
}

// CreateCourseInput — структура ДЛЯ валидации (с тегами)
type CreateCourseInput struct {
	SquadID     int64     `validate:"required,min=1"`
	Year        int       `validate:"required,min=2020,max=2100"`
	Title       string    `validate:"required,min=1,max=255"`
	Description *string   `validate:"omitempty"`
	StartsAt    time.Time `validate:"required"`
	EndsAt      time.Time `validate:"required"`
}

// Validate выполняет кастомную бизнес-валидацию
func (c *CreateCourseInput) Validate() error {
	if !c.EndsAt.After(c.StartsAt) {
		return fmt.Errorf("ends_at must be after starts_at")
	}
	return nil
}
