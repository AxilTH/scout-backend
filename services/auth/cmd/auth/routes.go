// cmd/auth/routes.go
package main

import (
	"github.com/AxilTH/scout-backend/services/auth/internal/handler"
	"github.com/AxilTH/scout-backend/services/auth/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/AxilTH/scout-backend/services/auth/internal/repository"
	"log"
)

func setupRoutes(r *gin.Engine, db *sqlx.DB, jwtSecret string) {
	log.Println("Setting up routes")
	// Health check
	r.GET("/health", handler.HealthCheck)
	r.HEAD("/health", handler.HealthCheck)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			authRepo, err := repository.NewAuthRepository(db)
			if err != nil {
				panic(err)
			}
			log.Println("Registering auth routes")
			auth.POST("/login", handler.LoginHandler(authRepo.UserRepository, jwtSecret))
			auth.POST("/refresh", handler.RefreshHandler(jwtSecret))
			auth.POST("/logout", handler.LogoutHandler())
			auth.POST("/switch-squad", middleware.AuthRequired(jwtSecret), handler.SwitchSquadHandler(jwtSecret))
			// User profile routes (protected)
			auth.GET("/me", middleware.AuthRequired(jwtSecret), handler.GetCurrentUserHandler(authRepo.UserRepository))
			auth.GET("/users/:id", middleware.AuthRequired(jwtSecret), handler.GetUserHandler(authRepo.UserRepository))
			auth.PUT("/users/:id", middleware.AuthRequired(jwtSecret), handler.UpdateUserHandler(authRepo.UserRepository))
		}

		// User routes (placeholder for future)
		// users := v1.Group("/users")
		// {
		// 		// TODO: Добавить маршруты для пользователей
		// 		// users.GET("/me", h.GetCurrentUser)
		// }
	}
}
