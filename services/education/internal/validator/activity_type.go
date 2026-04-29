// internal/validator/activity_type.go
package validator

// CreateActivityTypeRequest — структура только для десериализации JSON
type CreateActivityTypeRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
}

// CreateActivityTypeInput — структура для структурной валидации (теги `validate`)
type CreateActivityTypeInput struct {
	Title       string  `validate:"required,min=1,max=255"`
	Description *string `validate:"omitempty"`
	SquadID     int64   `validate:"required,min=1"`
}
