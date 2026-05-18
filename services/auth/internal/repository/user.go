// internal/repository/user.go
package repository

import (
	"context"

	"github.com/AxilTH/scout-backend/services/auth/internal/model"
)

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	BaseRepository[model.User]
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByIDs(ctx context.Context, ids []int64) ([]*model.User, error)
	// GetSquadIDsByUserID returns squad IDs that the user belongs to.
	GetSquadIDsByUserID(ctx context.Context, userID int64) ([]int64, error)
}