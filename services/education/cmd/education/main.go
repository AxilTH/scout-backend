// cmd/education/main.go — точка входа для сервиса Education
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

	"github.com/AxilTH/scout-backend/services/education/internal/handler"
	"github.com/AxilTH/scout-backend/services/education/internal/middleware"
	"github.com/AxilTH/scout-backend/services/education/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

// config хранит конфигурацию сервиса
type config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	ServicePort string
}

// loadConfig загружает и валидирует конфигурацию из окружения
func loadConfig() (*config, error) {
	cfg := &config{
		// Только инфраструктурные параметры
		DBHost:      getEnv("DB_HOST", "postgresql"),
		DBPort:      getEnv("DB_PORT", "5432"),      
		ServicePort: getEnv("SERVICE_PORT", "8081"),
	}

	// Критичные переменные - обязательны к установке
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
		return nil, fmt.Errorf("missing required environment variables: %v", missing)
	}

	return cfg, nil
}

func main() {
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

	// 5. Инициализация репозиториев
	courseRepo := repository.NewCourseRepository(db)
	activityTypeRepo := repository.NewActivityTypeRepository(db)
	activityRepo := repository.NewActivityRepository(db)
	// assignmentRepo := repository.NewAssignmentRepository(db)
	// activityResultRepo := repository.NewActivityResultRepository(db)
	// courseMentorshipRepo := repository.NewCourseMentorshipRepository(db)

	// 6. Настройка маршрутов
	secured := r.Group("")
	secured.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		secured.GET("/courses/:id/activities/:activity_id", func(c *gin.Context) {
			handler.GetActivity(c, activityRepo)
		})
		secured.PUT("/courses/:id/activities/:activity_id", func(c *gin.Context) {
			handler.UpdateActivity(c, activityRepo)
		})
		secured.DELETE("/courses/:id/activities/:activity_id", func(c *gin.Context) {
			handler.DeleteActivity(c, activityRepo)
		})

		secured.GET("/courses/:id/activities", func(c *gin.Context) {
			handler.GetActivities(c, activityRepo)
		})
		secured.POST("/courses/:id/activities", func(c *gin.Context) {
			handler.CreateActivity(c, activityRepo, courseRepo)
		})

		secured.GET("/courses/:id", func(c *gin.Context) {
			handler.GetCourse(c, courseRepo)
		})
		secured.PUT("/courses/:id", func(c *gin.Context) {
			handler.UpdateCourse(c, courseRepo)
		})
		secured.DELETE("/courses/:id", func(c *gin.Context) {
			handler.DeleteCourse(c, courseRepo)
		})

		secured.GET("/courses", func(c *gin.Context) {
			handler.GetCourses(c, courseRepo)
		})
		secured.POST("/courses", func(c *gin.Context) {
			handler.CreateCourse(c, courseRepo)
		})

		secured.GET("/activity-types", func(c *gin.Context) {
			handler.GetActivityTypes(c, activityTypeRepo)
		})
		secured.GET("/activity-types/:id", func(c *gin.Context) {
			handler.GetActivityType(c, activityTypeRepo)
		})
		secured.POST("/activity-types", func(c *gin.Context) {
			handler.CreateActivityType(c, activityTypeRepo)
		})
		secured.PUT("/activity-types/:id", func(c *gin.Context) {
			handler.UpdateActivityType(c, activityTypeRepo)
		})
		secured.DELETE("/activity-types/:id", func(c *gin.Context) {
			handler.DeleteActivityType(c, activityTypeRepo)
		})
	}
	// r.POST("/activities/:id/submit", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.SubmitActivity(c, activityResultRepo) })

	// r.GET("/assignments", func(c *gin.Context) { handler.GetAssignments(c, assignmentRepo) })
	// r.GET("/assignments/:id", func(c *gin.Context) { handler.GetAssignment(c, assignmentRepo) })
	// r.POST("/assignments", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.CreateAssignment(c, assignmentRepo) })
	// r.PUT("/assignments/:id", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.UpdateAssignment(c, assignmentRepo) })
	// r.DELETE("/assignments/:id", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.DeleteAssignment(c, assignmentRepo) })

	// r.GET("/activity-results/:id", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.GetActivityResult(c, activityResultRepo) })
	// r.PUT("/activity-results/:id/review", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.ReviewActivityResult(c, activityResultRepo) })
	// r.GET("/courses/:id/results", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.GetCourseResults(c, activityResultRepo) })
	// r.GET("/users/:id/results", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.GetUserResults(c, activityResultRepo) })

	// r.GET("/courses/:id/mentors", func(c *gin.Context) { handler.GetCourseMentors(c, courseMentorshipRepo) })
	// r.POST("/courses/:id/mentors", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.AssignMentor(c, courseMentorshipRepo) })
	// r.DELETE("/courses/:id/mentors/:user_id", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.RemoveMentor(c, courseMentorshipRepo) })

	r.GET("/health", handler.HealthCheck)

	// 7. Graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.ServicePort,
		Handler: r,
	}

	go func() {
		log.Printf("INFO: Starting Education Service on port %s", cfg.ServicePort)
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

