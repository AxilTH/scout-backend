// internal/repository/squad_membership.go
package repository

import (
	"context"

	"github.com/AxilTH/scout-backend/services/squad/internal/model"
)

// SquadMembershipRepository определяет интерфейс для работы с членством в отрядах
type SquadMembershipRepository interface {
	BaseRepository[model.SquadMembership]

	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.SquadMembership, error)
	GetBySquadID(ctx context.Context, squadID int64, limit, offset int) ([]*model.SquadMembership, error)
	GetActiveByUserID(ctx context.Context, userID int64) (*model.SquadMembership, error)
	GetActiveBySquadID(ctx context.Context, squadID int64, limit, offset int) ([]*model.SquadMembership, error)
}