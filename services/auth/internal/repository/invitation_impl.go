// internal/repository/invitation_impl.go
package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/AxilTH/scout-backend/services/auth/internal/model"
	"github.com/jmoiron/sqlx"
)

// invitationRepositoryImpl implements InvitationRepository
type invitationRepositoryImpl struct {
	db *sqlx.DB
}

// Create creates a new invitation (base repository method)
func (r *invitationRepositoryImpl) Create(ctx context.Context, entity *model.Invitation) error {
	now := time.Now().UTC()
	entity.CreatedAt = now
	entity.UpdatedAt = now

	query := `
		INSERT INTO invitations (email, squad_id, role_id, expires_at, used_at, used_by, created_by, created_at, updated_at)
		VALUES (:email, :squad_id, :role_id, :expires_at, :used_at, :used_by, :created_by, :created_at, :updated_at)
		RETURNING id`
	rows, err := r.db.NamedQueryContext(ctx, query, entity)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&entity.ID)
	}
	return sql.ErrNoRows
}

// GetByID returns an invitation by ID (base repository method)
func (r *invitationRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.Invitation, error) {
	var inv model.Invitation
	err := r.db.GetContext(ctx, &inv, "SELECT id, email, squad_id, role_id, expires_at, used_at, used_by, created_by, created_at, updated_at FROM invitations WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

// Update updates an existing invitation
func (r *invitationRepositoryImpl) Update(ctx context.Context, entity *model.Invitation) error {
	query := `
		UPDATE invitations
		SET email = :email,
		    squad_id = :squad_id,
		    role_id = :role_id,
		    expires_at = :expires_at,
		    used_at = :used_at,
		    used_by = :used_by,
		    created_by = :created_by,
		    created_at = :created_at,
		    updated_at = :updated_at
		WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, entity)
	return err
}

// Delete deletes an invitation by ID
func (r *invitationRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM invitations WHERE id = $1", id)
	return err
}

// List returns a list of invitations with limit and offset
func (r *invitationRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*model.Invitation, error) {
	var invs []*model.Invitation
	err := r.db.SelectContext(ctx, &invs, "SELECT id, email, squad_id, role_id, expires_at, used_at, used_by, created_by, created_at, updated_at FROM invitations LIMIT $1 OFFSET $2", limit, offset)
	return invs, err
}

// CreateInvitation creates a new invitation with explicit parameters
func (r *invitationRepositoryImpl) CreateInvitation(ctx context.Context, email string, squadID int64, roleID int64, createdBy int64) (*model.Invitation, error) {
	inv := &model.Invitation{
		Email:     email,
		SquadID:   squadID,
		RoleID:    roleID,
		CreatedBy: createdBy,
	}
	err := r.Create(ctx, inv)
	if err != nil {
		return nil, err
	}
	return inv, nil
}

// GetInvitation returns an invitation by ID (wrapper for GetByID)
func (r *invitationRepositoryImpl) GetInvitation(ctx context.Context, id int64) (*model.Invitation, error) {
	return r.GetByID(ctx, id)
}

// GetInvitationByEmail returns an invitation by email
func (r *invitationRepositoryImpl) GetInvitationByEmail(ctx context.Context, email string) (*model.Invitation, error) {
	var inv model.Invitation
	err := r.db.GetContext(ctx, &inv, "SELECT id, email, squad_id, role_id, expires_at, used_at, used_by, created_by, created_at, updated_at FROM invitations WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

// UseInvitation marks an invitation as used by setting used_at and used_by
func (r *invitationRepositoryImpl) UseInvitation(ctx context.Context, id int64, userID int64) error {
	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, "UPDATE invitations SET used_at = $1, used_by = $2 WHERE id = $3", now, userID, id)
	return err
}

// GetValidInvitations returns valid invitations for a squad (not used and not expired)
func (r *invitationRepositoryImpl) GetValidInvitations(ctx context.Context, squadID int64, limit, offset int) ([]*model.Invitation, error) {
	var invs []*model.Invitation
	err := r.db.SelectContext(ctx, &invs, "SELECT id, email, squad_id, role_id, expires_at, used_at, used_by, created_by, created_at, updated_at FROM invitations WHERE squad_id = $1 AND used_at IS NULL AND expires_at > NOW() ORDER BY created_at LIMIT $2 OFFSET $3", squadID, limit, offset)
	return invs, err
}