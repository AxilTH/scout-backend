// internal/handler/invitation.go
package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/AxilTH/scout-backend/services/auth/internal/repository"
	"github.com/gin-gonic/gin"
)

// CreateInvitationHandler handles creating a new invitation.
// It expects JSON with email, squadID, roleID, and createdBy.
func CreateInvitationHandler(invitationRepo repository.InvitationRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" binding:"required,email"`
			SquadID  int64  `json:"squad_id" binding:"required"`
			RoleID   int64  `json:"role_id" binding:"required"`
			CreatedBy int64  `json:"created_by" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
			return
		}

		invitation, err := invitationRepo.CreateInvitation(c, input.Email, input.SquadID, input.RoleID, input.CreatedBy)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create invitation"})
			return
		}

		c.JSON(http.StatusCreated, invitation)
	}
}

// GetInvitationHandler returns an invitation by ID.
func GetInvitationHandler(invitationRepo repository.InvitationRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		var id int64
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invitation ID"})
			return
		}

		invitation, err := invitationRepo.GetInvitation(c, id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "invitation not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch invitation"})
			}
			return
		}
		c.JSON(http.StatusOK, invitation)
	}
}

// GetInvitationsHandler returns a list of invitations for a squad with pagination.
// It accepts squad_id as query parameter.
func GetInvitationsHandler(invitationRepo repository.InvitationRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		squadIDStr := c.Query("squad_id")
		if squadIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "squad_id query parameter is required"})
			return
		}

		var squadID int64
		if _, err := fmt.Sscanf(squadIDStr, "%d", &squadID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid squad_id"})
			return
		}

		var limit, offset int
		if l := c.Query("limit"); l != "" {
			fmt.Sscanf(l, "%d", &limit)
		} else {
			limit = 100 // default limit
		}
		if o := c.Query("offset"); o != "" {
			fmt.Sscanf(o, "%d", &offset)
		} else {
			offset = 0 // default offset
		}

		invitations, err := invitationRepo.GetValidInvitations(c, squadID, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch invitations"})
			return
		}
		c.JSON(http.StatusOK, invitations)
	}
}