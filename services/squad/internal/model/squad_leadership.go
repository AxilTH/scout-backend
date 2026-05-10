// internal/model/squad_leadership.go
package model

import "time"

// SquadLeadership представляет должность в командном составе отряда
type SquadLeadership struct {
	ID int64 `db:"id" json:"id"`

	// Timestamps
	AppointedAt time.Time  `db:"appointed_at" json:"appointed_at"`
	DismissedAt *time.Time `db:"dismissed_at" json:"dismissed_at,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`

	// Внешние ключи
	UserID     int64 `db:"user_id" json:"user_id"`
	SquadID    int64 `db:"squad_id" json:"squad_id"`
	PositionID int64 `db:"position_id" json:"position_id"`
}

// TableName specifies the table name for SquadLeadership
func (SquadLeadership) TableName() string {
	return "squad_leaderships"
}
