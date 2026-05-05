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
}