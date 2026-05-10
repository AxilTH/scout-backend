// internal/model/region.go
package model

import "time"

// Region представляет регион
type Region struct {
	ID        int64     `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`

	// Timestamps
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName specifies the table name for Region
func (Region) TableName() string {
	return "regions"
}