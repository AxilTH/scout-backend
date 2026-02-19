package validator

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateCourseRequest struct {
	SquadID     string    `json:"squad_id"`
	Year        int       `json:"year"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
}

type CreateCourseInput struct {
	SquadID     string    `validate:"required,uuid4"`
	Year        int       `validate:"required,min=2020,max=2100"`
	Title       string    `validate:"required,min=1,max=255"`
	Description *string   `validate:"-"`
	StartsAt    time.Time `validate:"required"`
	EndsAt      time.Time `validate:"required"`
}

func NewValidator() *validator.Validate {
	v := validator.New()
	return v
}

func (c *CreateCourseInput) Validate() error {
	if !c.EndsAt.After(c.StartsAt) {
		return fmt.Errorf("ends_at must be after starts_at")
	}
	return nil
}

func ValidateSquadID(squadID string) error {
	_, err := uuid.Parse(squadID)
	return err
}

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
				return fmt.Sprintf("%s is too small", field)
			case "max":
				if field == "Year" {
					return "Year out of range"
				}
				return fmt.Sprintf("%s is too large", field)
			case "uuid4":
				return "Invalid UUID format for squad_id"
			default:
				return fmt.Sprintf("Invalid value for %s", field)
			}
		}
	}
	return "Invalid input data"
}
