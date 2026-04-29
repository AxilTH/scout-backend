// internal/repository/course_tx.go
package repository

import (
	"database/sql"
	"time"

	"github.com/AxilTH/scout-backend/services/education/internal/model"
	"github.com/jmoiron/sqlx"
)

// CourseTxRepository представляет репозиторий курсов с поддержкой транзакций
type CourseTxRepository struct {
	*TransactionalRepository
}

// NewCourseTxRepository создает новый репозиторий курсов с поддержкой транзакций
func NewCourseTxRepository(db *sqlx.DB) *CourseTxRepository {
	return &CourseTxRepository{
		TransactionalRepository: NewTransactionalRepository(db),
	}
}

// GetAll возвращает все курсы, принадлежащие отряду с указанным squad_id
func (r *CourseTxRepository) GetAll(squadID int64) ([]model.Course, error) {
	var courses []model.Course
	err := r.DB().Select(&courses, "SELECT * FROM courses WHERE squad_id = $1 ORDER BY year DESC", squadID)
	return courses, err
}

// GetByID возвращает курс по id, только если он принадлежит указанному отряду
// Возвращает nil, nil если курс не найден или не принадлежит отряду
func (r *CourseTxRepository) GetByID(id, squadID int64) (*model.Course, error) {
	var course model.Course
	err := r.DB().Get(&course, "SELECT * FROM courses WHERE id = $1 AND squad_id = $2", id, squadID)
	if err == sql.ErrNoRows {
		return nil, nil // курс не существует или не принадлежит отряду
	}
	return &course, err
}

// Create создает новый курс в транзакции
func (r *CourseTxRepository) Create(course *model.Course) error {
	course.CreatedAt = time.Now().UTC()
	course.UpdatedAt = time.Now().UTC()

	query := `
		INSERT INTO courses (title, description, starts_at, ends_at, year,
		   is_active, created_at, updated_at, squad_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	err := r.DB().QueryRow(
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

// CreateInTx создает новый курс в указанной транзакции
func (r *CourseTxRepository) CreateInTx(tx *sqlx.Tx, course *model.Course) error {
	course.CreatedAt = time.Now().UTC()
	course.UpdatedAt = time.Now().UTC()

	query := `
		INSERT INTO courses (title, description, starts_at, ends_at, year,
		   is_active, created_at, updated_at, squad_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	err := tx.QueryRow(
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
func (r *CourseTxRepository) Update(course *model.Course) error {
	course.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE courses
		SET title = $1, description = $2, starts_at = $3, ends_at = $4, year = $5,
		   is_active = $6, updated_at = $7
		WHERE id = $8 AND squad_id = $9
	`

	_, err := r.DB().Exec(
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

// UpdateInTx обновляет существующий курс в указанной транзакции
func (r *CourseTxRepository) UpdateInTx(tx *sqlx.Tx, course *model.Course) error {
	course.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE courses
		SET title = $1, description = $2, starts_at = $3, ends_at = $4, year = $5,
		   is_active = $6, updated_at = $7
		WHERE id = $8 AND squad_id = $9
	`

	_, err := tx.Exec(
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
func (r *CourseTxRepository) Delete(id, squadID int64) error {
	_, err := r.DB().Exec("DELETE FROM courses WHERE id = $1 AND squad_id = $2", id, squadID)
	return err
}

// DeleteInTx удаляет курс в указанной транзакции
func (r *CourseTxRepository) DeleteInTx(tx *sqlx.Tx, id, squadID int64) error {
	_, err := tx.Exec("DELETE FROM courses WHERE id = $1 AND squad_id = $2", id, squadID)
	return err
}

// ExistsForSquad проверяет, существует ли курс с указанным ID и принадлежит ли он отряду.
// Возвращает (true, nil) если курс найден и принадлежит отряду.
// Возвращает (false, nil) если курс не найден или не принадлежит отряду.
// Возвращает (false, err) при ошибке БД.
func (r *CourseTxRepository) ExistsForSquad(courseID, squadID int64) (bool, error) {
	var id int64
	err := r.DB().Get(&id, `
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

// ExistsForSquadInTx проверяет существование курса в указанной транзакции
func (r *CourseTxRepository) ExistsForSquadInTx(tx *sqlx.Tx, courseID, squadID int64) (bool, error) {
	var id int64
	err := tx.Get(&id, `
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