// internal/repository/db.go
package repository

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func NewPostgreSQLDB(host, port, user, password, dbname string) (*sqlx.DB, error) {
	data_source_name := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", data_source_name)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL")
	return db, nil
}
