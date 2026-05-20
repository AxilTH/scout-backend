// cmd/auth/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AxilTH/scout-backend/services/auth/internal/middleware"
	"github.com/AxilTH/scout-backend/services/auth/internal/repository"
	"github.com/AxilTH/scout-backend/services/auth/internal/squadclient"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	ServicePort string
	JWTSecret   string
}

func loadConfig() (*config, error) {
	cfg := &config{
		DBHost:      getEnv("DB_HOST", "postgresql"),
		DBPort:      getEnv("DB_PORT", "5432"),
		ServicePort: getEnv("AUTH_SERVICE_PORT", "8083"),
	}

	required := map[string]*string{
		"DB_USER":     &cfg.DBUser,
		"DB_PASSWORD": &cfg.DBPassword,
		"DB_NAME":     &cfg.DBName,
		"JWT_SECRET":  &cfg.JWTSecret,
	}

	var missing []string
	for key, target := range required {
		val := os.Getenv(key)
		if val == "" {
			missing = append(missing, key)
		} else {
			*target = val
		}
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variable: %v", missing)
	}

	return cfg, nil
}

func main() {
	// Загрузка .env файлов
	if _, err := os.Stat("../../../.env"); err == nil {
		if err := godotenv.Load("../../../.env"); err != nil {
			log.Println("Warning: Root .env file not found, using system environment variables")
		}
	}
	if _, err := os.Stat("../../.env"); err == nil {
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("Warning: Service-specific .env file not found, using system environment variables")
		}
	}

	// Загрузка конфигурации
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("FATAL: Configuration error: %v", err)
	}

	// Инициализация БД
	db, err := repository.NewPostgreSQLDB(
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize DB: %v", err)
	}
	defer db.Close()

	// Применение миграций
	if err := applyMigrations(db); err != nil {
		log.Fatalf("FATAL: Failed to apply migrations: %v", err)
	}
	log.Println("INFO: Migrations applied successfully")

	// TODO: исключить эту логику через потоки
	// Инициализация клиента Squad Service
	squadServiceURL := os.Getenv("SQUAD_SERVICE_URL")
	if squadServiceURL == "" {
		squadServiceURL = "http://squad-service:8082" // Значение по умолчанию из docker-compose
	}
	squadTimeout := 10 * time.Second
	if v := os.Getenv("SQUAD_SERVICE_TIMEOUT"); v != "" {
		if parsed, err := time.ParseDuration(v); err == nil {
			squadTimeout = parsed
		}
	}
	squadClient := squadclient.NewSquadClient(squadServiceURL, squadTimeout)

	// Инициализация репозиториев
	authRepo, err := repository.NewAuthRepository(db, squadClient)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize repository: %v", err)
	}
	// Using blank identifier to avoid unused variable error while we implement handlers
	_ = authRepo

	// Настройка Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLoggerMiddleware(middleware.DefaultLoggingConfig()))

	// Настройка маршрутов
	setupRoutes(r, db, cfg.JWTSecret)

	// Graceful shutdown
	srv := &http.Server{
		Addr: ":" + cfg.ServicePort,
		Handler: r,
	}

	go func() {
		log.Printf("INFO: Starting Auth & Identity Service on port %s", cfg.ServicePort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("FATAL: Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("INFO: Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("ERROR: Server forced to shutdown: %v", err)
	}
	log.Println("INFO: Server exited gracefully")
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func applyMigrations(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}