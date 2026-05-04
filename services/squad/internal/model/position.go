// internal/model/position.go
package model

import "time"

// Position представляет должность в командном составе отряда
type Position struct {
	ID        int64     `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`

	// Timestamps
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName specifies the table name for Position
func (Position) TableName() string {
	return "positions"
}
