package model

import "time"

type Activity struct {
	ID             string    `db:"id" json:"id"`
	CourseID       string    `db:"course_id" json:"course_id"`
	ActivityTypeID string    `db:"activity_type_id" json:"activity_type_id"`
	Title          string    `db:"title" json:"title"`
	Description    *string   `db:"description" json:"description,omitempty"`
	StartsAt       time.Time `db:"starts_at" json:"starts_at"`
	EndsAt         time.Time `db:"ends_at" json:"ends_at"`
	MaxScore       *int      `db:"max_score" json:"max_score,omitempty"`
	IsPublished    bool      `db:"is_published" json:"is_published"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}
