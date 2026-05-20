// internal/handler/switch_squad.go
package handler

import (
	"net/http"
	"time"

	"github.com/AxilTH/scout-backend/services/auth/internal/token"
	"github.com/AxilTH/scout-backend/services/auth/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SwitchSquadRequest defines the expected JSON body for switching squad.
type SwitchSquadRequest struct {
	SquadID int64 `json:"squad_id" binding:"required"`
}

// SwitchSquadHandler allows an authenticated user to change their current squad.
// It validates that the requested squad_id is present in the user's squad_ids claim.
// On success, it returns a new JWT token with the same squad_ids and the new current_squad_id.
func SwitchSquadHandler(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SwitchSquadRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
			return
		}

		// Extract claims from context (set by AuthRequired middleware)
		userID := middleware.GetUserIDFromContext(c)
		squadIDs := middleware.GetSquadIDsFromContext(c)

		if userID == 0 || len(squadIDs) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated or no squad info"})
			return
		}

		// Verify that the requested squad_id is among the user's squads
		allowed := false
		for _, id := range squadIDs {
			if id == req.SquadID {
				allowed = true
				break
			}
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "squad not accessible to user"})
			return
		}

		// Generate a new token with updated current_squad_id.
		token, err := token.GenerateJWT(userID, squadIDs, req.SquadID, jwtSecret, 2*time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token": token,
			"token_type":   "Bearer",
			"expires_in":   int64(2 * time.Hour.Seconds()),
		})
	}
}