// internal/repository/db.go
package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func NewPostgreSQLDB(host, port, user, password, dbname string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
