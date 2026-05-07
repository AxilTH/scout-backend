// internal/validator/leadership.go
package validator

// AssignPositionRequest — структура только для десериализации JSON
type AssignPositionRequest struct {
	UserID     int64 `json:"user_id"`
	SquadID    int64 `json:"squad_id"`
	PositionID int64 `json:"position_id"`
}

// AssignPositionInput — структура для структурной валидации (теги `validate`)
type AssignPositionInput struct {
	UserID     int64 `validate:"required,min=1"`
	SquadID    int64 `validate:"required,min=1"`
	PositionID int64 `validate:"required,min=1"`
}

// Validate выполняет бизнес-валидацию и возвращает ВСЕ ошибки.
func (a *AssignPositionInput) Validate() ValidationErrors {
	var errs ValidationErrors

	// Бизнес-правила могут быть добавлены здесь
	// Например: проверка, что пользователь является fighter

	return errs
}

// UpdateLeadershipPositionRequest — структура только для десериализации JSON
type UpdateLeadershipPositionRequest struct {
	PositionID  int64  `json:"position_id"`
	DismissedAt *int64 `json:"dismissed_at,omitempty"` // timestamp
}

// UpdateLeadershipPositionInput — структура для структурной валидации (теги `validate`)
type UpdateLeadershipPositionInput struct {
	PositionID  int64  `validate:"required,min=1"`
	DismissedAt *int64 `validate:"omitempty"`
}

// Validate выполняет бизнес-валидацию и возвращает ВСЕ ошибки.
func (u *UpdateLeadershipPositionInput) Validate() ValidationErrors {
	var errs ValidationErrors

	// Бизнес-правила могут быть добавлены здесь

	return errs
}