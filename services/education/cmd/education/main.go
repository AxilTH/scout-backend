package main

import (
	"log"
	"os"

	"github.com/AxilTH/scout-backend/services/education/internal/handler"
	"github.com/AxilTH/scout-backend/services/education/internal/middleware"
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
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

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
	log.Println("Migrations applied successfully")

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	courseRepo := repository.NewCourseRepository(db)
	activityRepo := repository.NewActivityRepository(db)
	assignmentRepo := repository.NewAssignmentRepository(db)
	activityResultRepo := repository.NewActivityResultRepository(db)

	r.GET("/courses", func(c *gin.Context) { handler.GetCourses(c, courseRepo) })
	r.GET("/courses/:id", func(c *gin.Context) { handler.GetCourse(c, courseRepo) })
	r.POST("/courses", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.CreateCourse(c, courseRepo) })

	r.GET("/activities", func(c *gin.Context) { handler.GetActivities(c, activityRepo) })
	r.GET("/activities/:id", func(c *gin.Context) { handler.GetActivity(c, activityRepo) })
	r.POST("/activities", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.CreateActivity(c, activityRepo) })
	r.PUT("/activities/:id", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.UpdateActivity(c, activityRepo) })
	r.DELETE("/activities/:id", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.DeleteActivity(c, activityRepo) })

	r.POST("/activities/:id/submit", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.SubmitActivity(c, activityResultRepo) })

	r.GET("/assignments", func(c *gin.Context) { handler.GetAssignments(c, assignmentRepo) })
	r.GET("/assignments/:id", func(c *gin.Context) { handler.GetAssignment(c, assignmentRepo) })
	r.POST("/assignments", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.CreateAssignment(c, assignmentRepo) })
	r.PUT("/assignments/:id", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.UpdateAssignment(c, assignmentRepo) })
	r.DELETE("/assignments/:id", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.DeleteAssignment(c, assignmentRepo) })

	r.GET("/activity-results/:id", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.GetActivityResult(c, activityResultRepo) })
	r.PUT("/activity-results/:id/review", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.ReviewActivityResult(c, activityResultRepo) })
	r.GET("/courses/:id/results", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.GetCourseResults(c, activityResultRepo) })
	r.GET("/users/:id/results", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) { handler.GetUserResults(c, activityResultRepo) })

	r.GET("/health", handler.HealthCheck)

	log.Println("Starting Education Service on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
