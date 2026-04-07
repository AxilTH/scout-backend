package repository

import (
	"database/sql"
	"time"

	"github.com/AxilTH/scout-backend/services/education/internal/model"
	"github.com/jmoiron/sqlx"
)

type ActivityRepository struct {
	db *sqlx.DB
}

func NewActivityRepository(db *sqlx.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) Create(activity *model.Activity) error {
	activity.CreatedAt = time.Now()
	activity.UpdatedAt = time.Now()

	query := `
		INSERT INTO activities
		(course_id, activity_type_id, title, description, starts_at, ends_at, max_score, is_published, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		activity.CourseID,
		activity.ActivityTypeID,
		activity.Title,
		activity.Description,
		activity.StartsAt,
		activity.EndsAt,
		activity.MaxScore,
		activity.IsPublished,
		activity.CreatedAt,
		activity.UpdatedAt,
	).Scan(&activity.ID)

	return err
}

func (r *ActivityRepository) GetByID(id int64) (*model.Activity, error) {
	var activity model.Activity

	err := r.db.Get(&activity, "SELECT * FROM activities WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &activity, err
}

func (r *ActivityRepository) ListByCourse(courseId int64) ([]model.Activity, error) {
	var activities []model.Activity
	err := r.db.Select(&activities, "SELECT * FROM activities WHERE course_id = $1 ORDER BY starts_at ASC", courseId)

	return activities, err
}

func (r *ActivityRepository) Update(activity *model.Activity) error {
	activity.UpdatedAt = time.Now()

	query := `
		UPDATE activities
		SET title = $1, description = $2, starts_at = $3, ends_at = $4,
			 max_score = $5, is_published = $6, updated_at = $7
		WHERE id = $8
	`

	_, err := r.db.Exec(
		query,
		activity.Title,
		activity.Description,
		activity.StartsAt,
		activity.EndsAt,
		activity.MaxScore,
		activity.IsPublished,
		activity.UpdatedAt,
		activity.ID,
	)

	return err
}

func (r *ActivityRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM activities WHERE id = $1", id)
	return err
}
