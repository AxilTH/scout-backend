// internal/middleware/auth.go
package middleware

import (
	"net/http"
	"log"

	"github.com/AxilTH/scout-backend/services/auth/internal/token"
	"github.com/gin-gonic/gin"
)

// AuthRequired returns a middleware that validates JWT token and
// stores userID, squadIDs and currentSquadID in the gin context.
// secret is used to sign/verify the token.
func AuthRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		// Expect format "Bearer <token>"
		const bearerPrefix = "Bearer "
		if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}
		tokenString := authHeader[len(bearerPrefix):]

		// Parse and validate token
		claims, err := token.ParseJWT(tokenString, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Debug logging
		log.Printf("AuthMiddleware: UserID=%d, SquadIDs=%v, CurrentSquadID=%d",
			claims.UserID, claims.SquadIDs, claims.CurrentSquad)

		// Store claims in context for handlers
		c.Set("userID", claims.UserID)
		c.Set("squadIDs", claims.SquadIDs)
		c.Set("currentSquadID", claims.CurrentSquad)

		c.Next()
	}
}

// GetUserIDFromContext returns userID stored in the context by AuthRequired middleware.
// Returns 0 if not present.
func GetUserIDFromContext(c *gin.Context) int64 {
	if val, exists := c.Get("userID"); exists {
		if uid, ok := val.(int64); ok {
			return uid
		}
	}
	return 0
}

// GetSquadIDsFromContext returns squadIDs slice stored in the context.
// Returns empty slice if not present.
func GetSquadIDsFromContext(c *gin.Context) []int64 {
	if val, exists := c.Get("squadIDs"); exists {
		if ids, ok := val.([]int64); ok {
			return ids
		}
	}
	return nil
}

// GetCurrentSquadIDFromContext returns currentSquadID stored in the context.
// Returns 0 if not present.
func GetCurrentSquadIDFromContext(c *gin.Context) int64 {
	if val, exists := c.Get("currentSquadID"); exists {
		if csid, ok := val.(int64); ok {
			return csid
		}
	}
	return 0
}
