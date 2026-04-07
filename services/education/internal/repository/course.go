package repository

import (
	"database/sql"

	"github.com/AxilTH/scout-backend/services/education/internal/model"
	"github.com/jmoiron/sqlx"
)

type CourseRepository struct {
	db *sqlx.DB
}

func NewCourseRepository(db *sqlx.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) Create(course *model.Course) error {
	query := `
		INSERT INTO courses (squad_id, year, title, description, starts_at, ends_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		course.SquadID,
		course.Year,
		course.Title,
		course.Description,
		course.StartsAt,
		course.EndsAt,
		course.IsActive,
	).Scan(&course.ID)
	return err
}

func (r *CourseRepository) GetByID(id int64) (*model.Course, error) {
	var course model.Course
	err := r.db.Get(&course, "SELECT * FROM courses WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &course, err
}

func (r *CourseRepository) ListBySquad(squadID int64) ([]model.Course, error) {
	var courses []model.Course
	err := r.db.Select(&courses, "SELECT * FROM courses WHERE squad_id = $1 ORDER BY year DESC", squadID)
	return courses, err
}
