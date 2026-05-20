// internal/model/invitation.go
package model

import (
	"time"
)

// Invitation представляет модель приглашения в отряд
type Invitation struct {
	ID        int64     `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	SquadID   int64     `db:"squad_id" json:"squad_id"`
	RoleID    int64     `db:"role_id" json:"role_id"` // Роль пользователя в отряде
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	UsedAt    *time.Time `db:"used_at" json:"used_at"` // NULL пока не использовано
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	// Внешние ключи
	CreatedBy int64 `db:"created_by" json:"created_by"` // User.id
}