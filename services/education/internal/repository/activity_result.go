package repository

import (
	"database/sql"
	"time"

	"github.com/AxilTH/scout-backend/services/education/internal/model"
	"github.com/jmoiron/sqlx"
)

type ActivityResultRepository struct {
	db *sqlx.DB
}

func NewActivityResultRepository(db *sqlx.DB) *ActivityResultRepository {
	return &ActivityResultRepository{db: db}
}

func (r *ActivityResultRepository) Create(result *model.ActivityResult) error {
	result.CreatedAt = time.Now()
	result.UpdatedAt = time.Now()

	query := `
		INSERT INTO activity_results
		(user_id, activity_id, status, submitted_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		result.UserID,
		result.ActivityID,
		result.Status,
		result.SubmittedAt,
		result.CreatedAt,
		result.UpdatedAt,
	).Scan(&result.ID)

	return err
}

func (r *ActivityResultRepository) GetByID(id int64) (*model.ActivityResult, error) {
	var result model.ActivityResult
	err := r.db.Get(&result, "SELECT * FROM activity_results WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (r *ActivityResultRepository) GetByUserAndActivity(userID, activityID int64) (*model.ActivityResult, error) {
	var result model.ActivityResult
	err := r.db.Get(&result, "SELECT * FROM activity_results WHERE user_id = $1 AND activity_id = $2", userID, activityID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (r *ActivityResultRepository) ListByCourse(courseID int64) ([]model.ActivityResult, error) {
	var results []model.ActivityResult
	query := `
		SELECT ar.* FROM activity_results ar
		JOIN activities a ON ar.activity_id = a.id
		WHERE a.course_id = $1
		ORDER BY ar.updated_at DESC
	`
	err := r.db.Select(&results, query, courseID)
	return results, err
}

func (r *ActivityResultRepository) ListByUser(userID int64) ([]model.ActivityResult, error) {
	var results []model.ActivityResult
	err := r.db.Select(&results, "SELECT * FROM activity_results WHERE user_id = $1 ORDER BY updated_at DESC", userID)
	return results, err
}

func (r *ActivityResultRepository) Update(result *model.ActivityResult) error {
	result.UpdatedAt = time.Now()

	query := `
		UPDATE activity_results
		SET status = $1, score = $2, feedback = $3,
			submitted_at = $4, reviewed_at = $5, reviewer_id = $6, updated_at = $7
		WHERE id = $8
	`

	_, err := r.db.Exec(
		query,
		result.Status,
		result.Score,
		result.Feedback,
		result.SubmittedAt,
		result.ReviewedAt,
		result.ReviewerID,
		result.UpdatedAt,
		result.ID,
	)

	return err
}
