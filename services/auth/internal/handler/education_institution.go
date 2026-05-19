// internal/handler/education_institution.go
package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/AxilTH/scout-backend/services/auth/internal/model"
	"github.com/AxilTH/scout-backend/services/auth/internal/repository"
	"github.com/gin-gonic/gin"
)

// CreateEducationInstitutionHandler handles creating a new education institution.
// It expects JSON with title.
func CreateEducationInstitutionHandler(eduRepo repository.EducationInstitutionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var edu model.EducationInstitution
		if err := c.ShouldBindJSON(&edu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
			return
		}

		if err := eduRepo.Create(c, &edu); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create education institution"})
			return
		}

		c.JSON(http.StatusCreated, edu)
	}
}

// GetEducationInstitutionHandler returns an education institution by ID.
func GetEducationInstitutionHandler(eduRepo repository.EducationInstitutionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		var id int64
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid education institution ID"})
			return
		}

		edu, err := eduRepo.GetByID(c, id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "education institution not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch education institution"})
			}
			return
		}
		c.JSON(http.StatusOK, edu)
	}
}

// GetEducationInstitutionsHandler returns a list of education institutions with pagination.
func GetEducationInstitutionsHandler(eduRepo repository.EducationInstitutionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		institutions, err := eduRepo.List(c, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch education institutions"})
			return
		}
		c.JSON(http.StatusOK, institutions)
	}
}