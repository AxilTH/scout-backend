// internal/model/course.go
package model

import "time"

// Модель "Course"
// Связь с моделью "Squad" через ключ "squad_id"
type Course struct {
	ID          int64     `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description *string   `db:"description" json:"description,omitempty"`
	StartsAt    time.Time `db:"starts_at" json:"starts_at"`
	EndsAt      time.Time `db:"ends_at" json:"ends_at"`
	Year        int       `db:"year" json:"year"`
	IsActive    bool      `db:"is_active" json:"is_active"`

	// Timestamps
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`

	// Внешние ключи
	SquadID int64 `db:"squad_id" json:"squad_id"`
}
