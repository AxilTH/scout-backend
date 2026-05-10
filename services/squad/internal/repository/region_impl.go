// internal/repository/region_impl.go
package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/AxilTH/scout-backend/services/squad/internal/model"
)

type regionRepository struct {
	db *sqlx.DB
}

func NewRegionRepository(db *sqlx.DB) RegionRepository {
	return &regionRepository{db: db}
}

// GetByID возвращает регион по его идентификатору
func (r *regionRepository) GetByID(ctx context.Context, id int64) (*model.Region, error) {
	var region model.Region
	query := `SELECT id, title, created_at, updated_at FROM regions WHERE id = $1`
	err := r.db.GetContext(ctx, &region, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get region by id: %w", err)
	}
	return &region, nil
}

// GetByTitle возвращает регион по его названию
func (r *regionRepository) GetByTitle(ctx context.Context, title string) (*model.Region, error) {
	var region model.Region
	query := `SELECT id, title, created_at, updated_at FROM regions WHERE title = $1`
	err := r.db.GetContext(ctx, &region, query, title)
	if err != nil {
		return nil, fmt.Errorf("failed to get region by title: %w", err)
	}
	return &region, nil
}

// List возвращает список регионов с пагинацией
func (r *regionRepository) List(ctx context.Context, limit, offset int) ([]*model.Region, error) {
	var regions []*model.Region
	query := `SELECT id, title, created_at, updated_at FROM regions ORDER BY id LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &regions, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list regions: %w", err)
	}
	return regions, nil
}

// Create создает новый регион в базе данных
func (r *regionRepository) Create(ctx context.Context, entity *model.Region) error {
	query := `INSERT INTO regions (title, created_at, updated_at) VALUES (:title, :created_at, :updated_at) RETURNING id`
	rows, err := r.db.NamedQueryContext(ctx, query, entity)
	if err != nil {
		return fmt.Errorf("failed to create region: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&entity.ID)
	}

	return fmt.Errorf("failed to get inserted id")
}

// Update обновляет существующий регион в базе данных
func (r *regionRepository) Update(ctx context.Context, entity *model.Region) error {
	query := `UPDATE regions SET title = :title, updated_at = :updated_at WHERE id = :id`
	result, err := r.db.NamedExecContext(ctx, query, entity)
	if err != nil {
		return fmt.Errorf("failed to update region: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("region not found")
	}

	return nil
}

// Delete удаляет регион из базы данных по его идентификатору
func (r *regionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM regions WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete region: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("region not found")
	}

	return nil
}