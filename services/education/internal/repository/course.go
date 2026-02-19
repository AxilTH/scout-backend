package repository

import (
	"database/sql"

	"github.com/AxilTH/scout-backend/services/education/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CourseRepository struct {
	db *sqlx.DB
}

func NewCourseRepository(db *sqlx.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) Create(course *model.Course) error {
	course.ID = uuid.New().String()
	course.IsActive = true
	query := `
		INSERT INTO courses (id, squad_id, year, title, description, starts_at, ends_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(
		query,
		course.ID,
		course.SquadID,
		course.Year,
		course.Title,
		course.Description,
		course.StartsAt,
		course.EndsAt,
		course.IsActive,
	)
	return err
}

func (r *CourseRepository) GetByID(id string) (*model.Course, error) {
	var course model.Course
	err := r.db.Get(&course, "SELECT * FROM courses WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &course, err
}

func (r *CourseRepository) List(squadID string) ([]model.Course, error) {
	var courses []model.Course
	err := r.db.Select(&courses, "SELECT * FROM courses WHERE squad_id = $1 ORDER BY year DESC", squadID)
	return courses, err
}