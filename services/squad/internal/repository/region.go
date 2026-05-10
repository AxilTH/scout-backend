// internal/repository/region.go
package repository

import (
	"context"

	"github.com/AxilTH/scout-backend/services/squad/internal/model"
)

// RegionRepository определяет интерфейс для работы с регионами
type RegionRepository interface {
	BaseRepository[model.Region]

	GetByTitle(ctx context.Context, title string) (*model.Region, error)
}