// internal/repository/education_institution_impl.go
package repository

import (
	"context"
	"database/sql"

	"github.com/AxilTH/scout-backend/services/auth/internal/model"
	"github.com/jmoiron/sqlx"
)

// educationInstitutionRepositoryImpl implements EducationInstitutionRepository
type educationInstitutionRepositoryImpl struct {
	db *sqlx.DB
}

// Create creates a new education institution
func (r *educationInstitutionRepositoryImpl) Create(ctx context.Context, entity *model.EducationInstitution) error {
	query := `
		INSERT INTO education_institutions (title, short_name, created_at, updated_at)
		VALUES (:title, :short_name, :created_at, :updated_at)
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

// GetByID returns an education institution by ID
func (r *educationInstitutionRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.EducationInstitution, error) {
	var edu model.EducationInstitution
	err := r.db.GetContext(ctx, &edu, "SELECT id, title, short_name, created_at, updated_at FROM education_institutions WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &edu, nil
}

// Update updates an existing education institution
func (r *educationInstitutionRepositoryImpl) Update(ctx context.Context, entity *model.EducationInstitution) error {
	query := `
		UPDATE education_institutions
		SET title = :title,
		    short_name = :short_name,
		    created_at = :created_at,
		    updated_at = :updated_at
		WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, entity)
	return err
}

// Delete deletes an education institution by ID
func (r *educationInstitutionRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM education_institutions WHERE id = $1", id)
	return err
}

// List returns a list of education institutions with limit and offset
func (r *educationInstitutionRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*model.EducationInstitution, error) {
	var eds []*model.EducationInstitution
	err := r.db.SelectContext(ctx, &eds, "SELECT id, title, short_name, created_at, updated_at FROM education_institutions LIMIT $1 OFFSET $2", limit, offset)
	return eds, err
}

// GetByTitle returns an education institution by title
func (r *educationInstitutionRepositoryImpl) GetByTitle(ctx context.Context, title string) (*model.EducationInstitution, error) {
	var edu model.EducationInstitution
	err := r.db.GetContext(ctx, &edu, "SELECT id, title, short_name, created_at, updated_at FROM education_institutions WHERE title = $1", title)
	if err != nil {
		return nil, err
	}
	return &edu, nil
}