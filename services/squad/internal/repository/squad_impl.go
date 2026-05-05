// internal/repository/squad_impl.go
package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/AxilTH/scout-backend/services/squad/internal/model"
)

type squadRepository struct {
	db *sqlx.DB
}

func NewSquadRepository(db *sqlx.DB) SquadRepository {
	return &squadRepository{db: db}
}

// GetByID возвращает отряд по его идентификатору
func (r *squadRepository) GetByID(ctx context.Context, id int64) (*model.Squad, error) {
	var squad model.Squad
	query := `SELECT id, title, region_id, created_at, updated_at FROM squads WHERE id = $1`
	err := r.db.GetContext(ctx, &squad, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get squad by id: %w", err)
	}
	return &squad, nil
}

// GetByTitle возвращает отряд по его названию
func (r *squadRepository) GetByTitle(ctx context.Context, title string) (*model.Squad, error) {
	var squad model.Squad
	query := `SELECT id, title, region_id, created_at, updated_at FROM squads WHERE title = $1`
	err := r.db.GetContext(ctx, &squad, query, title)
	if err != nil {
		return nil, fmt.Errorf("failed to get squad by title: %w", err)
	}
	return &squad, nil
}

// GetByRegionID возвращает список отрядов по идентификатору региона с пагинацией
func (r *squadRepository) GetByRegionID(ctx context.Context, regionID int64, limit, offset int) ([]*model.Squad, error) {
	var squads []*model.Squad
	query := `SELECT id, title, region_id, created_at, updated_at FROM squads WHERE region_id = $1 ORDER BY id LIMIT $2 OFFSET $3`
	err := r.db.SelectContext(ctx, &squads, query, regionID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get squads by region id: %w", err)
	}
	return squads, nil
}

// List возвращает список отрядов с пагинацией
func (r *squadRepository) List(ctx context.Context, limit, offset int) ([]*model.Squad, error) {
	var squads []*model.Squad
	query := `SELECT id, title, region_id, created_at, updated_at FROM squads ORDER BY id LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &squads, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list squads: %w", err)
	}
	return squads, nil
}

// Create создает новый отряд в базе данных
func (r *squadRepository) Create(ctx context.Context, entity *model.Squad) error {
	query := `INSERT INTO squads (title, region_id, created_at, updated_at) VALUES (:title, :region_id, :created_at, :updated_at) RETURNING id`
	rows, err := r.db.NamedQueryContext(ctx, query, entity)
	if err != nil {
		return fmt.Errorf("failed to create squad: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&entity.ID)
	}

	return fmt.Errorf("failed to get inserted id")
}

// Update обновляет существующий отряд в базе данных
func (r *squadRepository) Update(ctx context.Context, entity *model.Squad) error {
	query := `UPDATE squads SET title = :title, region_id = :region_id, updated_at = :updated_at WHERE id = :id`
	result, err := r.db.NamedExecContext(ctx, query, entity)
	if err != nil {
		return fmt.Errorf("failed to update squad: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("squad not found")
	}

	return nil
}

// Delete удаляет отряд из базы данных по его идентификатору
func (r *squadRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM squads WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete squad: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("squad not found")
	}

	return nil
}