// internal/repository/activity.go
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

// GetAll возвращает все активности, принадлежащие курсу с указанным course_id и отряду squadID.
// Проверяет что курс принадлежит отряду (multi-tenancy protection).
func (r *ActivityRepository) GetAll(courseID, squadID int64) ([]model.Activity, error) {
	var activities []model.Activity
	err := r.db.Select(&activities, `
		SELECT a.* FROM activities a
		INNER JOIN courses c ON a.course_id = c.id
		WHERE a.course_id = $1 AND c.squad_id = $2
		ORDER BY a.starts_at ASC
	`, courseID, squadID)
	return activities, err
}

// GetByID возвращает активность по ID, только если она принадлежит курсу и курс принадлежит отряду.
// Это предотвращает доступ пользователя из одного отряда к активностям другого отряда.
// Возвращает (nil, nil), если активность не найдена или не принадлежит отряду.
func (r *ActivityRepository) GetByID(activityID, courseID, squadID int64) (*model.Activity, error) {
	var activity model.Activity
	err := r.db.Get(&activity, `
		SELECT a.* FROM activities a
		INNER JOIN courses c ON a.course_id = c.id
		WHERE a.id = $1 AND a.course_id = $2 AND c.squad_id = $3
	`, activityID, courseID, squadID)
	if err == sql.ErrNoRows {
		return nil, nil // не найдена или не принадлежит отряду
	}
	return &activity, err
}

// Create создает новую активность для курса.
// Убедитесь что courseID принадлежит squadID ПЕРЕД вызовом этого метода!
// Это должно быть проверено на уровне handler'а.
// Устанавливает CreatedAt и UpdatedAt в текущее время UTC.
func (r *ActivityRepository) Create(activity *model.Activity) error {
	activity.CreatedAt = time.Now().UTC()
	activity.UpdatedAt = time.Now().UTC()

	query := `
		INSERT INTO activities
		(title, description, starts_at, ends_at, max_score, is_published, created_at, updated_at, course_id, activity_type_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		activity.Title,
		activity.Description,
		activity.StartsAt,
		activity.EndsAt,
		activity.MaxScore,
		activity.IsPublished,
		activity.CreatedAt,
		activity.UpdatedAt,
		activity.CourseID,
		activity.ActivityTypeID,
	).Scan(&activity.ID)

	return err
}

// Update обновляет активность, только если она принадлежит курсу и курс принадлежит отряду.
// Обновляет поле UpdatedAt в текущее время UTC.
// Не позволяет изменить CourseID (привязку к курсу) - активность привязана к одному курсу.
func (r *ActivityRepository) Update(activity *model.Activity, squadID int64) error {
	activity.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE activities
		SET title = $1, description = $2, starts_at = $3, ends_at = $4,
			max_score = $5, is_published = $6, updated_at = $7
		WHERE id = $8 AND course_id = $9 AND course_id IN (
			SELECT id FROM courses WHERE squad_id = $10
		)
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
		activity.CourseID,
		squadID,
	)

	return err
}

// Delete удаляет активность, только если она принадлежит курсу и курс принадлежит отряду.
// Это предотвращает удаление активностей в курсах других отрядов.
func (r *ActivityRepository) Delete(activityID, courseID, squadID int64) error {
	_, err := r.db.Exec(`
		DELETE FROM activities
		WHERE id = $1 AND course_id = $2 AND course_id IN (
			SELECT id FROM courses WHERE squad_id = $3
		)
	`, activityID, courseID, squadID)
	return err
}
