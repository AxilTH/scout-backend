package handler

import (
	"log"
	"net/http"

	"github.com/AxilTH/scout-backend/services/education/internal/model"
	"github.com/AxilTH/scout-backend/services/education/internal/repository"
	"github.com/AxilTH/scout-backend/services/education/internal/validator"
	"github.com/gin-gonic/gin"
)

var validate = validator.NewValidator()

func GetCourses(c *gin.Context, repo *repository.CourseRepository) {
	squadID := c.Query("squad_id")
	if squadID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "squad_id query parameter is required"})
		return
	}

	if err := validator.ValidateSquadID(squadID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format for squad_id"})
		return
	}

	courses, err := repo.List(squadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusOK, courses)
}

func GetCourse(c *gin.Context, repo *repository.CourseRepository) {
	id := c.Param("id")
	course, err := repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if course == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	c.JSON(http.StatusOK, course)
}

func CreateCourse(c *gin.Context, repo *repository.CourseRepository) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication middleware not configured"})
		return
	}

	log.Printf("User %s is creating a course", userID)

	var req validator.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	input := validator.CreateCourseInput{
		SquadID:     req.SquadID,
		Year:        req.Year,
		Title:       req.Title,
		Description: req.Description,
		StartsAt:    req.StartsAt,
		EndsAt:      req.EndsAt,
	}

	if err := validate.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.FormatError(err)})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course := &model.Course{
		SquadID:     input.SquadID,
		Year:        input.Year,
		Title:       input.Title,
		Description: input.Description,
		StartsAt:    input.StartsAt,
		EndsAt:      input.EndsAt,
	}

	if err := repo.Create(course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, course)
}
