// cmd/squad/routes.go
package main

import (
	"github.com/AxilTH/scout-backend/services/squad/internal/handler"
	"github.com/AxilTH/scout-backend/services/squad/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// setupRoutes настраивает все маршруты приложения
func setupRoutes(r *gin.Engine, db *sqlx.DB) {
	// Инициализация репозиториев
	regionRepo := repository.NewRegionRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	positionRepo := repository.NewPositionRepository(db)
	squadRepo := repository.NewSquadRepository(db)
	squadMembershipRepo := repository.NewSquadMembershipRepository(db)
	squadLeadershipRepo := repository.NewSquadLeadershipRepository(db)

	// Инициализация handlers
	h := handler.NewHandler(
		regionRepo,
		roleRepo,
		positionRepo,
		squadRepo,
		squadMembershipRepo,
		squadLeadershipRepo,
	)

	// Health check
	r.GET("/health", handler.HealthCheck)
	r.HEAD("/health", handler.HealthCheck)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Region routes
		regions := v1.Group("/regions")
		{
			regions.POST("", h.CreateRegion)
			regions.GET("", h.GetRegions)
			regions.GET("/:id", h.GetRegion)
			regions.PUT("/:id", h.UpdateRegion)
			regions.DELETE("/:id", h.DeleteRegion)
			regions.GET("/:id/squads", h.GetSquadsByRegion)
		}

		// Role routes
		roles := v1.Group("/roles")
		{
			roles.POST("", h.CreateRole)
			roles.GET("", h.GetRoles)
			roles.GET("/:id", h.GetRole)
			roles.PUT("/:id", h.UpdateRole)
			roles.DELETE("/:id", h.DeleteRole)
		}

		// Position routes
		positions := v1.Group("/positions")
		{
			positions.POST("", h.CreatePosition)
			positions.GET("", h.GetPositions)
			positions.GET("/:id", h.GetPosition)
			positions.PUT("/:id", h.UpdatePosition)
			positions.DELETE("/:id", h.DeletePosition)
		}

		// Squad routes
		squads := v1.Group("/squads")
		{
			squads.POST("", h.CreateSquad)
			squads.GET("", h.GetSquads)
			squads.GET("/:id", h.GetSquad)
			squads.PUT("/:id", h.UpdateSquad)
			squads.DELETE("/:id", h.DeleteSquad)

			// Squad members
			squads.GET("/:id/members", h.GetSquadMembers)
			squads.GET("/:id/members/active", h.GetActiveSquadMembers)
			squads.POST("/:id/memberships", h.AddMembership)
			squads.GET("/:id/members/:user_id/role", h.GetUserRole)
			squads.PUT("/:id/members/:user_id", h.UpdateMembership)
			squads.DELETE("/:id/members/:user_id", h.RemoveMembership)

			// Squad leadership
			squads.GET("/:id/leadership", h.GetSquadLeadership)
			squads.POST("/:id/leadership", h.AssignPosition)
			squads.PUT("/:id/leadership/:user_id", h.UpdateLeadershipPosition)
			squads.DELETE("/:id/leadership/:user_id", h.DismissPosition)
		}

		// User routes
		users := v1.Group("/users")
		{
			users.GET("/:user_id/squads", h.GetUserSquads)
			users.GET("/:user_id/squads/active", h.GetUserActiveSquad)
			users.GET("/:user_id/leadership", h.GetUserLeadership)
		}
	}
}