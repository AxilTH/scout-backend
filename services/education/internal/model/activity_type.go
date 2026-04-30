// internal/model/activity_type.go
package model

import "time"

// Модель "ActivityTypes"
// Связь с моделью "Squad" через ключ "squad_id"
// Если SquadID == nil - тип глобальный (доступен всем отрядам)
// Если SquadID != nil - тип привязан к конкретному отряду
type ActivityType struct {
	ID          int64   `db:"id" json:"id"`
	Title       string  `db:"title" json:"title"`
	Description *string `db:"description" json:"description,omitempty"`

	// Timestamps
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	// Внешние ключи
	SquadID *int64 `db:"squad_id" json:"squad_id,omitempty"`
}
