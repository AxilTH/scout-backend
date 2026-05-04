// internal/model/squad.go
package model

import "time"

// Squad представляет студенческий педагогический отряд
type Squad struct {
	ID        int64     `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`

	// Timestamps
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	// Внешние ключи
	RegionID int64 `db:"region_id" json:"region_id"`
}

// TableName specifies the table name for Squad
func (Squad) TableName() string {
	return "squads"
}
