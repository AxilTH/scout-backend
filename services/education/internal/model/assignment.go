package model

import "time"

type Assignment struct {
	ID               int64     `db:"id" json:"id"`
	ActivityID       *int64    `db:"activity_id" json:"activity_id,omitempty"`
	AssignmentTypeID int64     `db:"assignment_type_id" json:"assignment_type_id"`
	Title            string    `db:"title" json:"title"`
	Description      *string   `db:"description" json:"description,omitempty"`
	Deadline         time.Time `db:"deadline" json:"deadline"`
	MaxScore         int       `db:"max_score" json:"max_score"`
	IsPublished      bool      `db:"is_published" json:"is_published"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}
