// internal/repository/squad.go
package repository

import (
	"context"

	"github.com/AxilTH/scout-backend/services/squad/internal/model"
)

// SquadRepository определяет интерфейс для работы с отрядами
type SquadRepository interface {
	BaseRepository[model.Squad]

	GetByTitle(ctx context.Context, title string) (*model.Squad, error)
	GetByRegionID(ctx context.Context, regionID int64, limit, offset int) ([]*model.Squad, error)
}