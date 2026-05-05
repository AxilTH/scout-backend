// internal/repository/position_impl.go
package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/AxilTH/scout-backend/services/squad/internal/model"
)

type positionRepository struct {
	db *sqlx.DB
}

func NewPositionRepository(db *sqlx.DB) PositionRepository {
	return &positionRepository{db: db}
}

// GetByID возвращает должность по её идентификатору
func (r *positionRepository) GetByID(ctx context.Context, id int64) (*model.Position, error) {
	var position model.Position
	query := `SELECT id, title, created_at, updated_at FROM positions WHERE id = $1`
	err := r.db.GetContext(ctx, &position, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get position by id: %w", err)
	}
	return &position, nil
}

// GetByTitle возвращает должность по её названию
func (r *positionRepository) GetByTitle(ctx context.Context, title string) (*model.Position, error) {
	var position model.Position
	query := `SELECT id, title, created_at, updated_at FROM positions WHERE title = $1`
	err := r.db.GetContext(ctx, &position, query, title)
	if err != nil {
		return nil, fmt.Errorf("failed to get position by title: %w", err)
	}
	return &position, nil
}

// List возвращает список должностей с пагинацией
func (r *positionRepository) List(ctx context.Context, limit, offset int) ([]*model.Position, error) {
	var positions []*model.Position
	query := `SELECT id, title, created_at, updated_at FROM positions ORDER BY id LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &positions, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list positions: %w", err)
	}
	return positions, nil
}

// Create создает новую должность в базе данных
func (r *positionRepository) Create(ctx context.Context, entity *model.Position) error {
	query := `INSERT INTO positions (title, created_at, updated_at) VALUES (:title, :created_at, :updated_at) RETURNING id`
	rows, err := r.db.NamedQueryContext(ctx, query, entity)
	if err != nil {
		return fmt.Errorf("failed to create position: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&entity.ID)
	}

	return fmt.Errorf("failed to get inserted id")
}

// Update обновляет существующую должность в базе данных
func (r *positionRepository) Update(ctx context.Context, entity *model.Position) error {
	query := `UPDATE positions SET title = :title, updated_at = :updated_at WHERE id = :id`
	result, err := r.db.NamedExecContext(ctx, query, entity)
	if err != nil {
		return fmt.Errorf("failed to update position: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("position not found")
	}

	return nil
}

// Delete удаляет должность из базы данных по её идентификатору
func (r *positionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM positions WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete position: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("position not found")
	}

	return nil
}