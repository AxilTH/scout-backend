// internal/repository/invitation.go
package repository

import (
	"context"

	"github.com/AxilTH/scout-backend/services/auth/internal/model"
)

// InvitationRepository определяет интерфейс для работы с приглашениями
type InvitationRepository interface {
	BaseRepository[model.Invitation]
	CreateInvitation(ctx context.Context, email string, squadID int64, roleID int64, createdBy int64) (*model.Invitation, error)
	GetInvitation(ctx context.Context, id int64) (*model.Invitation, error)
	GetInvitationByEmail(ctx context.Context, email string) (*model.Invitation, error)
	UseInvitation(ctx context.Context, id int64, userID int64) error
	GetValidInvitations(ctx context.Context, squadID int64, limit, offset int) ([]*model.Invitation, error)
}