package main

import (
	"log"
	"os"

	"github.com/AxilTH/scout-backend/services/education/internal/handler"
	"github.com/AxilTH/scout-backend/services/education/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	db, err := repository.NewPostgreSQLDB(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migrate driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Println("✅ Migrations applied successfully")

	r := gin.Default()

	r.GET("/health", handler.HealthCheck)

	// TODO: Запустить HTTP-сервер
	log.Println("🚀 Starting Education Service on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
