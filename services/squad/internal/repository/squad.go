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
	GetByIDAndSquadID(ctx context.Context, id, squadID int64) (*model.Squad, error)
	UpdateAndSquadID(ctx context.Context, entity *model.Squad, squadID int64) error
	DeleteAndSquadID(ctx context.Context, id, squadID int64) error
}