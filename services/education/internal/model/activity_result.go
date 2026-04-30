package model

// import "time"

// // Модель "ActivityResult"
// // Связь с моделью "User" через ключ "user_id" (Кандидат)
// // Связь с моделью "Activity" через ключ "activity_id"
// // Связь с моделью "User" через ключ "reviewer_id" (Куратор)
// type ActivityResult struct {
// 	ID          int64      `db:"id" json:"id"`
// 	Status      string     `db:"status" json:"status"`
// 	Score       *int       `db:"score" json:"score,omitempty"`
// 	Feedback    *string    `db:"feedback" json:"feedback,omitempty"`
// 	SubmittedAt *time.Time `db:"submitted_at" json:"submitted_at,omitempty"`
// 	ReviewedAt  *time.Time `db:"reviewed_at" json:"reviewed_at,omitempty"`
// 	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
// 	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`

// 	// Внешние ключи
// 	UserID     int64  `db:"user_id" json:"user_id"`
// 	ActivityID int64  `db:"activity_id" json:"activity_id"`
// 	ReviewerID *int64 `db:"reviewer_id" json:"reviewer_id,omitempty"`
// }
