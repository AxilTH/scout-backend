// internal/repository/auth_repository.go
package repository

import (
	"github.com/jmoiron/sqlx"
)

// AuthRepository содержит все репозитории для Auth сервиса
type AuthRepository struct {
	UserRepository                 UserRepository
	EducationInstitutionRepository EducationInstitutionRepository
	InvitationRepository           InvitationRepository
}

// NewAuthRepository создает новый AuthRepository со всеми репозиториями
func NewAuthRepository(db *sqlx.DB) (*AuthRepository, error) {
	userRepo := &userRepositoryImpl{db: db}
	eduRepo := &educationInstitutionRepositoryImpl{db: db}
	invRepo := &invitationRepositoryImpl{db: db}

	return &AuthRepository{
		UserRepository:                 userRepo,
		EducationInstitutionRepository: eduRepo,
		InvitationRepository:           invRepo,
	}, nil
}
