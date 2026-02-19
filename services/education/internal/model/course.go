package model

import "time"

type Course struct {
	ID          string    `db:"id" json:"id"`
	SquadID     string    `db:"squad_id" json:"squad_id"`
	Year        int       `db:"year" json:"year"`
	Title       string    `db:"title" json:"title"`
	Description *string   `db:"description" json:"description,omitempty"`
	StartsAt    time.Time `db:"starts_at" json:"starts_at"`
	EndsAt      time.Time `db:"ends_at" json:"ends_at"`
	IsActive    bool      `db:"is_active" json:"is_active"`
}