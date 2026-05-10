// internal/repository/repository.go
package repository

import "context"

// BaseRepository определяет базовые методы для всех репозиториев
type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id int64) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*T, error)
}