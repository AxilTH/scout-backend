// internal/repository/squad_leadership.go
package repository

import (
	"context"

	"github.com/AxilTH/scout-backend/services/squad/internal/model"
)

// SquadLeadershipRepository определяет интерфейс для работы с командным составом отрядов
type SquadLeadershipRepository interface {
	BaseRepository[model.SquadLeadership]

	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.SquadLeadership, error)
	GetBySquadID(ctx context.Context, squadID int64, limit, offset int) ([]*model.SquadLeadership, error)
	GetByIDAndSquadID(ctx context.Context, id, squadID int64) (*model.SquadLeadership, error)
	UpdateAndSquadID(ctx context.Context, entity *model.SquadLeadership, squadID int64) error
	DeleteAndSquadID(ctx context.Context, id, squadID int64) error
}