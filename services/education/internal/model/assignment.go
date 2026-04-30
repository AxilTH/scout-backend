package model

// import "time"

// // Модель "Assignment"
// // Связь с моделью "Activity" через ключ "activity_id"
// // Связь с моделью "AssignmentType" через ключ "assignment_type_id"
// type Assignment struct {
// 	ID          int64     `db:"id" json:"id"`
// 	Title       string    `db:"title" json:"title"`
// 	Description *string   `db:"description" json:"description,omitempty"`
// 	Deadline    time.Time `db:"deadline" json:"deadline"`
// 	MaxScore    int       `db:"max_score" json:"max_score"`
// 	IsPublished bool      `db:"is_published" json:"is_published"`
// 	CreatedAt   time.Time `db:"created_at" json:"created_at"`
// 	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`

// 	// Внешние ключи
// 	ActivityID       *int64 `db:"activity_id" json:"activity_id,omitempty"`
// 	AssignmentTypeID int64  `db:"assignment_type_id" json:"assignment_type_id"`
// }
