// internal/handler/auth.go
package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AxilTH/scout-backend/services/auth/internal/model"
	"github.com/AxilTH/scout-backend/services/auth/internal/token"
	"github.com/AxilTH/scout-backend/services/auth/internal/repository"
	"github.com/AxilTH/scout-backend/services/auth/internal/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest defines the expected JSON body for the login endpoint.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginHandler handles user login: checks credentials and returns a JWT token.
// It expects JSON with email and password.
func LoginHandler(userRepo repository.UserRepository, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("LoginHandler called")
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
			return
		}

		// Find user by email
		log.Printf("Looking up user by email: %s", req.Email)
		user, err := userRepo.GetByEmail(c, req.Email)
		if err != nil {
			log.Printf("Error getting user by email: %v", err)
			// Do not reveal whether the email exists
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		log.Printf("User found: ID=%d", user.ID)

		// Debug logging
		log.Printf("Login attempt for email: %s", req.Email)
		log.Printf("User found: ID=%d, Email=%s, HashLen=%d", user.ID, user.Email, len(user.PasswordHash))
		log.Printf("Password from request: %s", req.Password)
		log.Printf("Stored hash: %s", user.PasswordHash)

		// Compare password hash
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			log.Printf("Password comparison failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		log.Printf("Password comparison succeeded")

		// Fetch squad IDs for the user from repository
		squadIDs, err := userRepo.GetSquadIDsByUserID(c, user.ID)
		if err != nil {
			// If we cannot fetch squad IDs, we treat as no squads (but could also error)
			squadIDs = []int64{}
		}
		// Determine current squad: if there is at least one squad, set current to the first one.
		// In a more advanced system, we might store the last used squad in user profile or let client specify.
		currentSquadID := int64(0)
		if len(squadIDs) > 0 {
			currentSquadID = squadIDs[0]
		}

		// Generate JWT token (valid for 2 hours)
		token, err := token.GenerateJWT(user.ID, squadIDs, currentSquadID, jwtSecret, 2*time.Hour)
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

// RefreshRequest defines the expected JSON body for refresh endpoint.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshHandler handles token refresh: given a valid refresh token, returns a new access token.
// For simplicity, we reuse the same token as both access and refresh; in production you may want separate.
func RefreshHandler(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
			return
		}

		// Parse the refresh token (we accept the same token; ensure it's not expired)
		claims, err := token.ParseJWT(req.RefreshToken, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
			return
		}

		// Generate new token with same claims but new expiration (e.g., 2 hours)
		newToken, err := token.GenerateJWT(claims.UserID, claims.SquadIDs, claims.CurrentSquad, jwtSecret, 2*time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate new token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token": newToken,
			"token_type":   "Bearer",
			"expires_in":   int64(2 * time.Hour.Seconds()),
		})
	}
}

// LogoutHandler simply returns success; client should discard the token.
func LogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
	}
}

// GetUserHandler returns a user by ID.
func GetUserHandler(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		var id int64
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		user, err := userRepo.GetByID(c, id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
			}
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// UpdateUserHandler updates an existing user.
func UpdateUserHandler(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		var id int64
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		var req model.User
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
			return
		}
		// Ensure ID from path is used
		req.ID = id

		if err := userRepo.Update(c, &req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user updated"})
	}
}

// GetCurrentUserHandler returns the currently authenticated user from context.
func GetCurrentUserHandler(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := middleware.GetUserIDFromContext(c)
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		user, err := userRepo.GetByID(c, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
			}
			return
		}
		c.JSON(http.StatusOK, user)
	}
}