// internal/handler/register.go
package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/AxilTH/scout-backend/services/auth/internal/model"
	"github.com/AxilTH/scout-backend/services/auth/internal/repository"
	"github.com/AxilTH/scout-backend/services/auth/internal/squadclient"
	"github.com/AxilTH/scout-backend/services/auth/internal/token"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest defines the expected JSON body for the register endpoint.
type RegisterRequest struct {
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	InvitationID int64  `json:"invitation_id" binding:"required"`
}

// RegisterHandler handles user registration via invitation.
// It validates the invitation, creates the user, marks invitation as used, and returns JWT token.
func RegisterHandler(userRepo repository.UserRepository, invitationRepo repository.InvitationRepository, squadClient squadclient.SquadClient, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
			return
		}

		// Get invitation by ID
		invitation, err := invitationRepo.GetInvitation(c, req.InvitationID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "invitation not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch invitation"})
			}
			return
		}

		// Validate invitation
		if invitation.UsedAt != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invitation already used"})
			return
		}

		if invitation.Email != req.Email {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invitation email does not match"})
			return
		}

		// Check if user with this email already exists
		existingUser, _ := userRepo.GetByEmail(c, req.Email)
		if existingUser != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "user with this email already exists"})
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		// Create user
		user := &model.User{
			Email:        req.Email,
			PasswordHash: string(hashedPassword),
		}

		if err := userRepo.Create(c, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}

		// Mark invitation as used
		if err := invitationRepo.UseInvitation(c, req.InvitationID, user.ID); err != nil {
			// In a production system, we might want to rollback the user creation here
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark invitation as used"})
			return
		}

		// Create squad membership via Squad Service
		squadID := invitation.SquadID
		roleID := invitation.RoleID
		if err := squadClient.CreateMembership(c, user.ID, squadID, roleID); err != nil {
			// If we failed to create membership, we should delete the user to keep consistency
			if delErr := userRepo.Delete(c, user.ID); delErr != nil {
				// Log the error but return the original error
				// In a production system, we might want to use a more sophisticated rollback mechanism
				// For now, we just return the membership creation error
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create squad membership"})
			return
		}

		// Get user's squad IDs (now should include the newly created one)
		squadIDs, err := userRepo.GetSquadIDsByUserID(c, user.ID)
		if err != nil {
			squadIDs = []int64{}
		}
		currentSquadID := int64(0)
		if len(squadIDs) > 0 {
			currentSquadID = squadIDs[0]
		}

		// Generate JWT token (valid for 2 hours)
		accessToken, err := token.GenerateJWT(user.ID, squadIDs, currentSquadID, jwtSecret, 2*time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"access_token": accessToken,
			"token_type":   "Bearer",
			"expires_in":   int64(2 * time.Hour.Seconds()),
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
			},
		})
	}
}
