// cmd/auth/routes.go
package main

import (
	"github.com/AxilTH/scout-backend/services/auth/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func setupRoutes(r *gin.Engine, db *sqlx.DB, jwtSecret string) {
	// Health check
	r.GET("/health", handler.HealthCheck)
	r.HEAD("/health", handler.HealthCheck)

	// API v1 routes
	// v1 := r.Group("/api/v1")
	// {
	// 	// Auth routes
	// 	auth := v1.Group("/auth")
	// 	{
	// 		// TODO: Добавить маршруты для аутентификации
	// 		// auth.POST("/login", h.Login)
	// 		// auth.POST("/register", h.Register)
	// 		// auth.POST("/refresh", h.RefreshToken)
	// 		// auth.POST("/logout", h.Logout)
	// 	}

	// 	// User routes
	// 	users := v1.Group("/users")
	// 	{
	// 		// TODO: Добавить маршруты для пользователей
	// 		// users.GET("/me", h.GetCurrentUser)
	// 	}
	// }
}
