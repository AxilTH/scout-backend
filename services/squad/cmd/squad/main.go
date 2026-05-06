// cmd/squad/main.go
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

	"github.com/AxilTH/scout-backend/services/squad/internal/middleware"
	"github.com/AxilTH/scout-backend/services/squad/internal/repository"
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
}

func loadConfig() (*config, error) {
	cfg := &config{
		// Только инфраструктурные параметры
		DBHost:      getEnv("DB_HOST", "postgresql"),
		DBPort:      getEnv("DB_PORT", "5432"),
		ServicePort: getEnv("SQUAD_SERVICE_PORT", "8082"),
	}

	// Критичные переменные - обязательны к установке
	required := map[string]*string{
		"DB_USER":     &cfg.DBUser,
		"DB_PASSWORD": &cfg.DBPassword,
		"DB_NAME":     &cfg.DBName,
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
		return nil, fmt.Errorf("missing required environment variables: %v", missing)
	}

	return cfg, nil
}

func main() {
	// 0. Загрузка .env файлов (только локально)
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

	// 1. Загрузка конфигурации
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("FATAL: Configuration error: %v", err)
	}

	// 2. Инициализация БД
	db, err := repository.NewPostgreSQLDB(
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize DB: %v", err)
	}
	defer db.Close()

	// 3. Применение миграций
	if err := applyMigrations(db); err != nil {
		log.Fatalf("FATAL: Failed to apply migrations: %v", err)
	}
	log.Println("INFO: Migrations applied successfully")

	// 4. Настройка Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLoggerMiddleware(middleware.DefaultLoggingConfig()))

	// 5. Настройка маршрутов
	setupRoutes(r, db)

	// 6. Graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.ServicePort,
		Handler: r,
	}

	go func() {
		log.Printf("INFO: Starting Squad & Membership Service on port %s", cfg.ServicePort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("FATAL: Failed to start server: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("INFO: Shutting down server...")

	// Контекст с таймаутом для завершения текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("ERROR: Server forced to shutdown: %v", err)
	}
	log.Println("INFO: Server exited gracefully")
}

// =============================================================================
// Вспомогательные функции
// =============================================================================

// getEnv возвращает значение переменной окружения или дефолт (только для безопасных параметров)
func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// applyMigrations применяет миграции к БД
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
