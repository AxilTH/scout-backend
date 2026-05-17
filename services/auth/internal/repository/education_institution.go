// internal/repository/education_institution.go
package repository

import (
	"context"

	"github.com/AxilTH/scout-backend/services/auth/internal/model"
)

// EducationInstitutionRepository определяет интерфейс для работы с учебными заведениями
type EducationInstitutionRepository interface {
	BaseRepository[model.EducationInstitution]
	GetByTitle(ctx context.Context, title string) (*model.EducationInstitution, error)
}