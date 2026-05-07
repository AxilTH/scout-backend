// cmd/squad/routes.go
package main

import (
	"github.com/AxilTH/scout-backend/services/squad/internal/handler"
	"github.com/AxilTH/scout-backend/services/squad/internal/middleware"
	"github.com/AxilTH/scout-backend/services/squad/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// setupRoutes настраивает все маршруты приложения
func setupRoutes(r *gin.Engine, db *sqlx.DB, jwtSecret string) {
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
		// Region routes — admin only
		regions := v1.Group("/regions")
		regions.Use(middleware.AuthRequired(jwtSecret))
		regions.Use(middleware.AdminOnly())
		{
			regions.POST("", h.CreateRegion)
			regions.GET("", h.GetRegions)
			regions.GET("/:id", h.GetRegion)
			regions.PUT("/:id", h.UpdateRegion)
			regions.DELETE("/:id", h.DeleteRegion)
		}

		// Region squads — auth required (fighter sees own squads in region)
		regionSquads := v1.Group("/regions/:id/squads")
		regionSquads.Use(middleware.AuthRequired(jwtSecret))
		{
			regionSquads.GET("", h.GetSquadsByRegion)
		}

		// Role routes — admin only
		roles := v1.Group("/roles")
		roles.Use(middleware.AuthRequired(jwtSecret))
		roles.Use(middleware.AdminOnly())
		{
			roles.POST("", h.CreateRole)
			roles.GET("", h.GetRoles)
			roles.GET("/:id", h.GetRole)
			roles.PUT("/:id", h.UpdateRole)
			roles.DELETE("/:id", h.DeleteRole)
		}

		// Position routes — admin only
		positions := v1.Group("/positions")
		positions.Use(middleware.AuthRequired(jwtSecret))
		positions.Use(middleware.AdminOnly())
		{
			positions.POST("", h.CreatePosition)
			positions.GET("", h.GetPositions)
			positions.GET("/:id", h.GetPosition)
			positions.PUT("/:id", h.UpdatePosition)
			positions.DELETE("/:id", h.DeletePosition)
		}

		// Squad routes — auth required, role-based filtering in handlers
		squads := v1.Group("/squads")
		squads.Use(middleware.AuthRequired(jwtSecret))
		{
			squads.POST("", h.CreateSquad)
			squads.GET("", h.GetSquads)
			squads.GET("/:id", h.GetSquad)
			squads.PUT("/:id", h.UpdateSquad)
			squads.DELETE("/:id", h.DeleteSquad)

			// Squad members — read only for all authenticated users
			squads.GET("/:id/members", h.GetSquadMembers)
			squads.GET("/:id/members/active", h.GetActiveSquadMembers)
			squads.GET("/:id/members/:user_id/role", h.GetUserRole)

			// Squad members management — commander only
			squadMembers := squads.Group("/:id")
			squadMembers.Use(middleware.CommanderOnly(squadLeadershipRepo, positionRepo))
			{
				squadMembers.POST("/memberships", h.AddMembership)
				squadMembers.PUT("/members/:user_id", h.UpdateMembership)
				squadMembers.DELETE("/members/:user_id", h.RemoveMembership)
			}

			// Squad leadership — read only for all authenticated users
			squads.GET("/:id/leadership", h.GetSquadLeadership)

			// Squad leadership management — commander only
			squadLeadership := squads.Group("/:id")
			squadLeadership.Use(middleware.CommanderOnly(squadLeadershipRepo, positionRepo))
			{
				squadLeadership.POST("/leadership", h.AssignPosition)
				squadLeadership.PUT("/leadership/:user_id", h.UpdateLeadershipPosition)
				squadLeadership.DELETE("/leadership/:user_id", h.DismissPosition)
			}
		}

		// User routes — auth required
		users := v1.Group("/users")
		users.Use(middleware.AuthRequired(jwtSecret))
		{
			users.GET("/:user_id/squads", h.GetUserSquads)
			users.GET("/:user_id/squads/active", h.GetUserActiveSquad)
			users.GET("/:user_id/leadership", h.GetUserLeadership)
		}
	}
}
