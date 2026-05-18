// internal/repository/user_impl.go
package repository

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/AxilTH/scout-backend/services/auth/internal/model"
	"github.com/jmoiron/sqlx"
)

// userRepositoryImpl implements UserRepository
type userRepositoryImpl struct {
	db *sqlx.DB
}

// Create creates a new user
func (r *userRepositoryImpl) Create(ctx context.Context, entity *model.User) error {
	query := `
		INSERT INTO users (first_name, last_name, middle_name, phone_number, email, date_of_birth, vk_profile_url, password_hash, created_at, updated_at)
		VALUES (:first_name, :last_name, :middle_name, :phone_number, :email, :date_of_birth, :vk_profile_url, :password_hash, :created_at, :updated_at)
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

// GetByID returns a user by ID
func (r *userRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT id, first_name, last_name, middle_name, phone_number, email, date_of_birth, vk_profile_url, password_hash, created_at, updated_at FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user
func (r *userRepositoryImpl) Update(ctx context.Context, entity *model.User) error {
	query := `
		UPDATE users
		SET first_name = :first_name,
		    last_name = :last_name,
		    middle_name = COALESCE(:middle_name, middle_name),
		    phone_number = COALESCE(:phone_number, phone_number),
		    email = :email,
		    date_of_birth = COALESCE(:date_of_birth, date_of_birth),
		    vk_profile_url = COALESCE(:vk_profile_url, vk_profile_url),
		    password_hash = :password_hash,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, entity)
	return err
}

// Delete deletes a user by ID
func (r *userRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}

// List returns a list of users with limit and offset
func (r *userRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	var users []*model.User
	err := r.db.SelectContext(ctx, &users, "SELECT id, first_name, last_name, middle_name, phone_number, email, date_of_birth, vk_profile_url, password_hash, created_at, updated_at FROM users LIMIT $1 OFFSET $2", limit, offset)
	return users, err
}

// GetByEmail returns a user by email
func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT id, first_name, last_name, middle_name, phone_number, email, date_of_birth, vk_profile_url, password_hash, created_at, updated_at FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByIDs returns users by a slice of IDs
func (r *userRepositoryImpl) GetByIDs(ctx context.Context, ids []int64) ([]*model.User, error) {
	if len(ids) == 0 {
		return []*model.User{}, nil
	}
	// Build query with placeholders
	query := "SELECT id, first_name, last_name, middle_name, phone_number, email, date_of_birth, vk_profile_url, password_hash, created_at, updated_at FROM users WHERE id IN ("
	args := []interface{}{}
	for i, id := range ids {
		if i > 0 {
			query += ", "
		}
		query += "$" + strconv.Itoa(i+1)
		args = append(args, id)
	}
	query += ")"

	var users []*model.User
	err := r.db.SelectContext(ctx, &users, query, args...)
	return users, err
}

// GetSquadIDsByUserID returns squad IDs that the user belongs to,
// based on invitations that have been used by the user (used_by = userID).
func (r *userRepositoryImpl) GetSquadIDsByUserID(ctx context.Context, userID int64) ([]int64, error) {
	var squadIDs []int64
	query := `SELECT DISTINCT squad_id FROM invitations WHERE used_by = $1`
	err := r.db.SelectContext(ctx, &squadIDs, query, userID)
	if err != nil {
		return nil, err
	}
	return squadIDs, nil
}