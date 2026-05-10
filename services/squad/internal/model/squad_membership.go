// internal/model/squad_membership.go
package model

import "time"

// SquadMembership представляет членство пользователя в отряде
type SquadMembership struct {
	ID       int64     `db:"id" json:"id"`
	IsActive bool      `db:"is_active" json:"is_active"`

	// Timestamps
	JoinedAt  time.Time  `db:"joined_at" json:"joined_at"`
	LeftAt    *time.Time `db:"left_at" json:"left_at,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`

	// Внешние ключи
	UserID  int64 `db:"user_id" json:"user_id"`
	SquadID int64 `db:"squad_id" json:"squad_id"`
	RoleID  int64 `db:"role_id" json:"role_id"`
}

// TableName specifies the table name for SquadMembership
func (SquadMembership) TableName() string {
	return "squad_memberships"
}
