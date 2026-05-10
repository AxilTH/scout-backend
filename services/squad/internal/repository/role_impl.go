// internal/repository/role_impl.go
package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/AxilTH/scout-backend/services/squad/internal/model"
)

type roleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) RoleRepository {
	return &roleRepository{db: db}
}

// GetByID возвращает роль по её идентификатору
func (r *roleRepository) GetByID(ctx context.Context, id int64) (*model.Role, error) {
	var role model.Role
	query := `SELECT id, title, created_at, updated_at FROM roles WHERE id = $1`
	err := r.db.GetContext(ctx, &role, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role by id: %w", err)
	}
	return &role, nil
}

// GetByTitle возвращает роль по её названию
func (r *roleRepository) GetByTitle(ctx context.Context, title string) (*model.Role, error) {
	var role model.Role
	query := `SELECT id, title, created_at, updated_at FROM roles WHERE title = $1`
	err := r.db.GetContext(ctx, &role, query, title)
	if err != nil {
		return nil, fmt.Errorf("failed to get role by title: %w", err)
	}
	return &role, nil
}

// List возвращает список ролей с пагинацией
func (r *roleRepository) List(ctx context.Context, limit, offset int) ([]*model.Role, error) {
	var roles []*model.Role
	query := `SELECT id, title, created_at, updated_at FROM roles ORDER BY id LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &roles, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	return roles, nil
}

// Create создает новую роль в базе данных
func (r *roleRepository) Create(ctx context.Context, entity *model.Role) error {
	query := `INSERT INTO roles (title, created_at, updated_at) VALUES (:title, :created_at, :updated_at) RETURNING id`
	rows, err := r.db.NamedQueryContext(ctx, query, entity)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&entity.ID)
	}

	return fmt.Errorf("failed to get inserted id")
}

// Update обновляет существующую роль в базе данных
func (r *roleRepository) Update(ctx context.Context, entity *model.Role) error {
	query := `UPDATE roles SET title = :title, updated_at = :updated_at WHERE id = :id`
	result, err := r.db.NamedExecContext(ctx, query, entity)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}

// Delete удаляет роль из базы данных по её идентификатору
func (r *roleRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM roles WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}