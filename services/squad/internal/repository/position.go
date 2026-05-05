// internal/repository/position.go
package repository

import (
	"context"

	"github.com/AxilTH/scout-backend/services/squad/internal/model"
)

// PositionRepository определяет интерфейс для работы с должностями
type PositionRepository interface {
	BaseRepository[model.Position]

	GetByTitle(ctx context.Context, title string) (*model.Position, error)
}