// internal/repository/role.go
package repository

import (
	"context"

	"github.com/AxilTH/scout-backend/services/squad/internal/model"
)

// RoleRepository определяет интерфейс для работы с ролями
type RoleRepository interface {
	BaseRepository[model.Role]

	GetByTitle(ctx context.Context, title string) (*model.Role, error)
}