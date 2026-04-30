//internal/model/activity.go
package model

import "time"

// Модель "Activity"
// Связь с моделью "Course" через ключ "course_id"
// Связь с моделью "ActivityType" через ключ "activity_type_id"
type Activity struct {
	ID          int64     `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description *string   `db:"description" json:"description,omitempty"`
	StartsAt    time.Time `db:"starts_at" json:"starts_at"`
	EndsAt      time.Time `db:"ends_at" json:"ends_at"`
	MaxScore    int      `db:"max_score" json:"max_score"`
	IsPublished bool      `db:"is_published" json:"is_published"`

	// Timestamps
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`

	// Внешние ключи
	CourseID       int64 `db:"course_id" json:"course_id"`
	ActivityTypeID int64 `db:"activity_type_id" json:"activity_type_id"`
}
