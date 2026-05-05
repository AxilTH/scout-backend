// internal/repository/squad_membership_impl.go
package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/AxilTH/scout-backend/services/squad/internal/model"
)

type squadMembershipRepository struct {
	db *sqlx.DB
}

func NewSquadMembershipRepository(db *sqlx.DB) SquadMembershipRepository {
	return &squadMembershipRepository{db: db}
}

// GetByID возвращает членство в отряде по его идентификатору
func (r *squadMembershipRepository) GetByID(ctx context.Context, id int64) (*model.SquadMembership, error) {
	var membership model.SquadMembership
	query := `SELECT id, is_active, joined_at, left_at, created_at, updated_at, user_id, squad_id, role_id FROM squad_memberships WHERE id = $1`
	err := r.db.GetContext(ctx, &membership, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get squad membership by id: %w", err)
	}
	return &membership, nil
}

// GetByUserID возвращает список членств пользователя в отрядах с пагинацией
func (r *squadMembershipRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.SquadMembership, error) {
	var memberships []*model.SquadMembership
	query := `SELECT id, is_active, joined_at, left_at, created_at, updated_at, user_id, squad_id, role_id FROM squad_memberships WHERE user_id = $1 ORDER BY id LIMIT $2 OFFSET $3`
	err := r.db.SelectContext(ctx, &memberships, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get squad memberships by user id: %w", err)
	}
	return memberships, nil
}

// GetBySquadID возвращает список членств в отряде с пагинацией
func (r *squadMembershipRepository) GetBySquadID(ctx context.Context, squadID int64, limit, offset int) ([]*model.SquadMembership, error) {
	var memberships []*model.SquadMembership
	query := `SELECT id, is_active, joined_at, left_at, created_at, updated_at, user_id, squad_id, role_id FROM squad_memberships WHERE squad_id = $1 ORDER BY id LIMIT $2 OFFSET $3`
	err := r.db.SelectContext(ctx, &memberships, query, squadID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get squad memberships by squad id: %w", err)
	}
	return memberships, nil
}

// GetActiveByUserID возвращает активное членство пользователя в отряде
func (r *squadMembershipRepository) GetActiveByUserID(ctx context.Context, userID int64) (*model.SquadMembership, error) {
	var membership model.SquadMembership
	query := `SELECT id, is_active, joined_at, left_at, created_at, updated_at, user_id, squad_id, role_id FROM squad_memberships WHERE user_id = $1 AND is_active = true LIMIT 1`
	err := r.db.GetContext(ctx, &membership, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active squad membership by user id: %w", err)
	}
	return &membership, nil
}

// GetActiveBySquadID возвращает список активных членств в отряде с пагинацией
func (r *squadMembershipRepository) GetActiveBySquadID(ctx context.Context, squadID int64, limit, offset int) ([]*model.SquadMembership, error) {
	var memberships []*model.SquadMembership
	query := `SELECT id, is_active, joined_at, left_at, created_at, updated_at, user_id, squad_id, role_id FROM squad_memberships WHERE squad_id = $1 AND is_active = true ORDER BY id LIMIT $2 OFFSET $3`
	err := r.db.SelectContext(ctx, &memberships, query, squadID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get active squad memberships by squad id: %w", err)
	}
	return memberships, nil
}

// List возвращает список всех членств в отрядах с пагинацией
func (r *squadMembershipRepository) List(ctx context.Context, limit, offset int) ([]*model.SquadMembership, error) {
	var memberships []*model.SquadMembership
	query := `SELECT id, is_active, joined_at, left_at, created_at, updated_at, user_id, squad_id, role_id FROM squad_memberships ORDER BY id LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &memberships, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list squad memberships: %w", err)
	}
	return memberships, nil
}

// Create создает новое членство в отряде в базе данных
func (r *squadMembershipRepository) Create(ctx context.Context, entity *model.SquadMembership) error {
	query := `INSERT INTO squad_memberships (is_active, joined_at, left_at, created_at, updated_at, user_id, squad_id, role_id) VALUES (:is_active, :joined_at, :left_at, :created_at, :updated_at, :user_id, :squad_id, :role_id) RETURNING id`
	rows, err := r.db.NamedQueryContext(ctx, query, entity)
	if err != nil {
		return fmt.Errorf("failed to create squad membership: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&entity.ID)
	}

	return fmt.Errorf("failed to get inserted id")
}

// Update обновляет существующее членство в отряде в базе данных
func (r *squadMembershipRepository) Update(ctx context.Context, entity *model.SquadMembership) error {
	query := `UPDATE squad_memberships SET is_active = :is_active, joined_at = :joined_at, left_at = :left_at, updated_at = :updated_at, user_id = :user_id, squad_id = :squad_id, role_id = :role_id WHERE id = :id`
	result, err := r.db.NamedExecContext(ctx, query, entity)
	if err != nil {
		return fmt.Errorf("failed to update squad membership: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("squad membership not found")
	}

	return nil
}

// Delete удаляет членство в отряде из базы данных по его идентификатору
func (r *squadMembershipRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM squad_memberships WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete squad membership: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("squad membership not found")
	}

	return nil
}