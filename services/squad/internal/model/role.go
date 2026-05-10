// internal/model/role.go
package model

import "time"

// Role представляет роль пользователя в отряде
type Role struct {
	ID        int64     `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`

	// Timestamps
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName specifies the table name for Role
func (Role) TableName() string {
	return "roles"
}
