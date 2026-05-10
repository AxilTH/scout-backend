// internal/repository/squad_leadership_impl.go
package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/AxilTH/scout-backend/services/squad/internal/model"
)

type squadLeadershipRepository struct {
	db *sqlx.DB
}

func NewSquadLeadershipRepository(db *sqlx.DB) SquadLeadershipRepository {
	return &squadLeadershipRepository{db: db}
}

// GetByID возвращает должность в командном составе отряда по её идентификатору
func (r *squadLeadershipRepository) GetByID(ctx context.Context, id int64) (*model.SquadLeadership, error) {
	var leadership model.SquadLeadership
	query := `SELECT id, appointed_at, dismissed_at, created_at, updated_at, user_id, squad_id, position_id FROM squad_leaderships WHERE id = $1`
	err := r.db.GetContext(ctx, &leadership, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get squad leadership by id: %w", err)
	}
	return &leadership, nil
}

// GetByUserID возвращает список должностей пользователя в командном составе отрядов с пагинацией
func (r *squadLeadershipRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.SquadLeadership, error) {
	var leaderships []*model.SquadLeadership
	query := `SELECT id, appointed_at, dismissed_at, created_at, updated_at, user_id, squad_id, position_id FROM squad_leaderships WHERE user_id = $1 ORDER BY id LIMIT $2 OFFSET $3`
	err := r.db.SelectContext(ctx, &leaderships, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get squad leaderships by user id: %w", err)
	}
	return leaderships, nil
}

// GetBySquadID возвращает список должностей в командном составе отряда с пагинацией
func (r *squadLeadershipRepository) GetBySquadID(ctx context.Context, squadID int64, limit, offset int) ([]*model.SquadLeadership, error) {
	var leaderships []*model.SquadLeadership
	query := `SELECT id, appointed_at, dismissed_at, created_at, updated_at, user_id, squad_id, position_id FROM squad_leaderships WHERE squad_id = $1 ORDER BY id LIMIT $2 OFFSET $3`
	err := r.db.SelectContext(ctx, &leaderships, query, squadID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get squad leaderships by squad id: %w", err)
	}
	return leaderships, nil
}

// List возвращает список всех должностей в командном составе отрядов с пагинацией
func (r *squadLeadershipRepository) List(ctx context.Context, limit, offset int) ([]*model.SquadLeadership, error) {
	var leaderships []*model.SquadLeadership
	query := `SELECT id, appointed_at, dismissed_at, created_at, updated_at, user_id, squad_id, position_id FROM squad_leaderships ORDER BY id LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &leaderships, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list squad leaderships: %w", err)
	}
	return leaderships, nil
}

// Create создает новую должность в командном составе отряда в базе данных
func (r *squadLeadershipRepository) Create(ctx context.Context, entity *model.SquadLeadership) error {
	query := `INSERT INTO squad_leaderships (appointed_at, dismissed_at, created_at, updated_at, user_id, squad_id, position_id) VALUES (:appointed_at, :dismissed_at, :created_at, :updated_at, :user_id, :squad_id, :position_id) RETURNING id`
	rows, err := r.db.NamedQueryContext(ctx, query, entity)
	if err != nil {
		return fmt.Errorf("failed to create squad leadership: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&entity.ID)
	}

	return fmt.Errorf("failed to get inserted id")
}

// Update обновляет существующую должность в командном составе отряда в базе данных
func (r *squadLeadershipRepository) Update(ctx context.Context, entity *model.SquadLeadership) error {
	query := `UPDATE squad_leaderships SET appointed_at = :appointed_at, dismissed_at = :dismissed_at, updated_at = :updated_at, user_id = :user_id, squad_id = :squad_id, position_id = :position_id WHERE id = :id`
	result, err := r.db.NamedExecContext(ctx, query, entity)
	if err != nil {
		return fmt.Errorf("failed to update squad leadership: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("squad leadership not found")
	}

	return nil
}

// Delete удаляет должность в командном составе отряда из базы данных по её идентификатору
func (r *squadLeadershipRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM squad_leaderships WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete squad leadership: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("squad leadership not found")
	}

	return nil
}

// GetByIDAndSquadID возвращает должность по ID, только если она принадлежит указанному отряду
// Возвращает nil, nil если должность не найдена или не принадлежит отряду
func (r *squadLeadershipRepository) GetByIDAndSquadID(ctx context.Context, id, squadID int64) (*model.SquadLeadership, error) {
	var leadership model.SquadLeadership
	query := `SELECT id, appointed_at, dismissed_at, created_at, updated_at, user_id, squad_id, position_id FROM squad_leaderships WHERE id = $1`
	err := r.db.GetContext(ctx, &leadership, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get squad leadership by id: %w", err)
	}

	// Проверяем, что должность принадлежит указанному squad_id
	if leadership.SquadID != squadID {
		return nil, fmt.Errorf("leadership does not belong to specified squad")
	}

	return &leadership, nil
}

// UpdateAndSquadID обновляет существующую должность, только если она принадлежит указанному отряду
func (r *squadLeadershipRepository) UpdateAndSquadID(ctx context.Context, entity *model.SquadLeadership, squadID int64) error {
	// Сначала проверяем, что должность существует и принадлежит squad_id
	existingLeadership, err := r.GetByID(ctx, entity.ID)
	if err != nil {
		return fmt.Errorf("leadership not found: %w", err)
	}

	if existingLeadership.SquadID != squadID {
		return fmt.Errorf("leadership does not belong to specified squad")
	}

	// Обновляем должность
	query := `UPDATE squad_leaderships SET appointed_at = :appointed_at, dismissed_at = :dismissed_at, updated_at = :updated_at, user_id = :user_id, squad_id = :squad_id, position_id = :position_id WHERE id = :id`
	result, err := r.db.NamedExecContext(ctx, query, entity)
	if err != nil {
		return fmt.Errorf("failed to update squad leadership: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("squad leadership not found")
	}

	return nil
}

// DeleteAndSquadID удаляет должность по ID, только если она принадлежит указанному отряду
func (r *squadLeadershipRepository) DeleteAndSquadID(ctx context.Context, id, squadID int64) error {
	// Сначала проверяем, что должность существует и принадлежит squad_id
	existingLeadership, err := r.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("leadership not found: %w", err)
	}

	if existingLeadership.SquadID != squadID {
		return fmt.Errorf("leadership does not belong to specified squad")
	}

	// Удаляем должность
	query := `DELETE FROM squad_leaderships WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete squad leadership: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("squad leadership not found")
	}

	return nil
}