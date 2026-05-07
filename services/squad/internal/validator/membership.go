// internal/validator/membership.go
package validator

import (
	"errors"
)

// AddMembershipRequest — структура только для десериализации JSON
type AddMembershipRequest struct {
	UserID  int64 `json:"user_id"`
	SquadID int64 `json:"squad_id"`
	RoleID  int64 `json:"role_id"`
}

// AddMembershipInput — структура для структурной валидации (теги `validate`)
type AddMembershipInput struct {
	UserID  int64 `validate:"required,min=1"`
	SquadID int64 `validate:"required,min=1"`
	RoleID  int64 `validate:"required,min=1"`
}

// Validate выполняет бизнес-валидацию и возвращает ВСЕ ошибки.
func (a *AddMembershipInput) Validate() ValidationErrors {
	var errs ValidationErrors

	// Бизнес-правила могут быть добавлены здесь
	// Например: проверка, что пользователь не состоит уже в этом отряде

	return errs
}

// UpdateMembershipRequest — структура только для десериализации JSON
type UpdateMembershipRequest struct {
	RoleID   int64  `json:"role_id"`
	IsActive bool   `json:"is_active"`
	LeftAt   *int64 `json:"left_at,omitempty"` // timestamp
}

// UpdateMembershipInput — структура для структурной валидации (теги `validate`)
type UpdateMembershipInput struct {
	RoleID   int64  `validate:"required,min=1"`
	IsActive bool   `validate:"required"`
	LeftAt   *int64 `validate:"omitempty"`
}

// Validate выполняет бизнес-валидацию и возвращает ВСЕ ошибки.
func (u *UpdateMembershipInput) Validate() ValidationErrors {
	var errs ValidationErrors

	// Если деактивируем членство, LeftAt должен быть установлен
	if !u.IsActive && u.LeftAt == nil {
		errs = append(errs, errors.New("left at is required when deactivating membership"))
	}

	// Если активируем членство, LeftAt должен быть nil
	if u.IsActive && u.LeftAt != nil {
		errs = append(errs, errors.New("left at must be nil when activating membership"))
	}

	return errs
}