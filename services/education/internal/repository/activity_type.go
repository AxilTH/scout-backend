// internal/repository/activity_type.go
package repository

import (
	"database/sql"
	"time"

	"github.com/AxilTH/scout-backend/services/education/internal/model"
	"github.com/jmoiron/sqlx"
)

type ActivityTypeRepository struct {
	db *sqlx.DB
}

func NewActivityTypeRepository(db *sqlx.DB) *ActivityTypeRepository {
	return &ActivityTypeRepository{db: db}
}

// GetAll возвращает все типы активностей, доступные указанному отряду:
// - глобальные (squad_id IS NULL) — доступны всем;
// - локальные этого отряда (squad_id = squadID).
func (r *ActivityTypeRepository) GetAll(squadID int64) ([]model.ActivityType, error) {
	var activityTypes []model.ActivityType
	err := r.db.Select(&activityTypes, `
		SELECT * FROM activity_types 
		WHERE squad_id IS NULL OR squad_id = $1 
		ORDER BY title ASC
	`, squadID)
	return activityTypes, err
}

// GetByID возвращает тип активности по ID, только если он доступен отряду:
// - глобальный (squad_id IS NULL) — доступен всем;
// - локальный (squad_id = squadID) — доступен только этому отряду.
// Возвращает (nil, nil), если тип не найден или не принадлежит отряду.
func (r *ActivityTypeRepository) GetByID(id, squadID int64) (*model.ActivityType, error) {
	var activityType model.ActivityType
	err := r.db.Get(&activityType, `
		SELECT * FROM activity_types 
		WHERE id = $1 AND (squad_id IS NULL OR squad_id = $2)
	`, id, squadID)
	if err == sql.ErrNoRows {
		return nil, nil // не найден или не принадлежит отряду
	}
	return &activityType, err
}

// Create создаёт новый тип активности.
// Если activityType.SquadID == nil - создаётся глобальный тип (только через сидирование).
// Если activityType.SquadID != nil - создаётся локальный тип для указанного отряда.
func (r *ActivityTypeRepository) Create(activityType *model.ActivityType) error {
	activityType.CreatedAt = time.Now().UTC()
	activityType.UpdatedAt = time.Now().UTC()

	query := `
	   INSERT INTO activity_types (title, description, created_at, updated_at, squad_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		activityType.Title,
		activityType.Description,
		activityType.CreatedAt,
		activityType.UpdatedAt,
		activityType.SquadID,
	).Scan(&activityType.ID)

	return err
}

// Update обновляет тип активности, только если он принадлежит отряду или является глобальным.
// Глобальные типы (squad_id = NULL) не должны обновляться через API — только через миграции.
// Этот метод предназначен для обновления локальных типов.
func (r *ActivityTypeRepository) Update(activityType *model.ActivityType) error {
	activityType.UpdatedAt = time.Now().UTC()

	query := `
	   UPDATE activity_types 
		SET title = $1, description = $2, updated_at = $3
		WHERE id = $4 AND squad_id = $5
	`

	_, err := r.db.Exec(
		query,
		activityType.Title,
		activityType.Description,
		activityType.UpdatedAt,
		activityType.ID,
		activityType.SquadID,
	)

	return err
}

// Delete удаляет тип активности, только если он является локальным для указанного отряда.
// Глобальные типы (squad_id = NULL) не могут быть удалены через этот метод.
// Для удаления локального типа необходимо передать корректный squadID.
func (r *ActivityTypeRepository) Delete(id, squadID int64) error {
	_, err := r.db.Exec(`
		DELETE FROM activity_types 
		WHERE id = $1 AND squad_id = $2
	`, id, squadID)
	return err
}
