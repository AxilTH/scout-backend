// internal/repository/course.go
package repository

import (
	"database/sql"
	"time"

	"github.com/AxilTH/scout-backend/services/education/internal/model"
	"github.com/jmoiron/sqlx"
)

type CourseRepository struct {
	db *sqlx.DB
}

func NewCourseRepository(db *sqlx.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

// GetAll возвращает все курсы, принадлежащие отряду с указанным squad_id
func (r *CourseRepository) GetAll(squadID int64) ([]model.Course, error) {
	var courses []model.Course
	err := r.db.Select(&courses, "SELECT * FROM courses WHERE squad_id = $1 ORDER BY year DESC", squadID)
	return courses, err
}

// GetByID возвращает курс по id, только если он принадлежит указанному отряду
// Возвращает nil, nil если курс не найден или не принадлежит отряду
func (r *CourseRepository) GetByID(id, squadID int64) (*model.Course, error) {
	var course model.Course
	err := r.db.Get(&course, "SELECT * FROM courses WHERE id = $1 AND squad_id = $2", id, squadID)
	if err == sql.ErrNoRows {
		return nil, nil // курс не существует или не принадлежит отряду
	}
	return &course, err
}

// Create создаёт новый курс. Устанавливает created_at и updated_at в текущее время UTC
func (r *CourseRepository) Create(course *model.Course) error {
	course.CreatedAt = time.Now().UTC()
	course.UpdatedAt = time.Now().UTC()

	query := `
		INSERT INTO courses (title, description, starts_at, ends_at, year, 
		   is_active, created_at, updated_at, squad_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		course.Title,
		course.Description,
		course.StartsAt,
		course.EndsAt,
		course.Year,
		course.IsActive,
		course.CreatedAt,
		course.UpdatedAt,
		course.SquadID,
	).Scan(&course.ID)
	
	return err
}

// Update обновляет существующий курс. Обновляет поле updated_at. Не изменяет squad_id и id
func (r *CourseRepository) Update(course *model.Course) error {
	course.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE courses
		SET title = $1, description = $2, starts_at = $3, ends_at = $4, year = $5, 
		   is_active = $6, updated_at = $7
		WHERE id = $8 AND squad_id = $9
	`

	_, err := r.db.Exec(
		query,
		course.Title,
		course.Description,
		course.StartsAt,
		course.EndsAt,
		course.Year,
		course.IsActive,
		course.UpdatedAt,
		course.ID,
		course.SquadID,
	)

	return err
}

// Delete удаляет курс по ID, только если он принадлежит указанному отряду
func (r *CourseRepository) Delete(id, squadID int64) error {
	_, err := r.db.Exec("DELETE FROM courses WHERE id = $1 AND squad_id = $2", id, squadID)
	return err
}

// ExistsForSquad проверяет, существует ли курс с указанным ID и принадлежит ли он отряду.
// Возвращает (true, nil) если курс найден и принадлежит отряду.
// Возвращает (false, nil) если курс не найден или не принадлежит отряду.
// Возвращает (false, err) при ошибке БД.
func (r *CourseRepository) ExistsForSquad(courseID, squadID int64) (bool, error) {
	var id int64
	err := r.db.Get(&id, `
		SELECT id FROM courses 
		WHERE id = $1 AND squad_id = $2
	`, courseID, squadID)
	
	if err == sql.ErrNoRows {
		return false, nil // курс не найден или не принадлежит отряду
	}
	if err != nil {
		return false, err // ошибка БД
	}
	return true, nil // курс существует и принадлежит отряду
}
