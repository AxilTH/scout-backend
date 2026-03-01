package handler

import (
	"log"
	"net/http"

	"github.com/AxilTH/scout-backend/services/education/internal/model"
	"github.com/AxilTH/scout-backend/services/education/internal/repository"
	"github.com/AxilTH/scout-backend/services/education/internal/validator"
	"github.com/gin-gonic/gin"
)

func GetActivities(c *gin.Context, repo *repository.ActivityRepository) {
	courseID := c.Query("course_id")
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "course_id query parameter is required"})
		return
	}

	activities, err := repo.ListByCourse(courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, activities)
}

func GetActivity(c *gin.Context, repo *repository.ActivityRepository) {
	id := c.Param("id")
	activity, err := repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if activity == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	c.JSON(http.StatusOK, activity)
}

func CreateActivity(c *gin.Context, repo *repository.ActivityRepository) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication middleware not configured"})
		return
	}
	log.Printf("User %s is creating an activity", userID)

	var req validator.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	input := validator.CreateActivityInput{
		CourseID:       req.CourseID,
		ActivityTypeID: req.ActivityTypeID,
		Title:          req.Title,
		Description:    req.Description,
		StartsAt:       req.StartsAt,
		EndsAt:         req.EndsAt,
		MaxScore:       req.MaxScore,
		IsPublished:    req.IsPublished,
	}

	if err := validate.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.FormatError(err)})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activity := &model.Activity{
		CourseID:       input.CourseID,
		ActivityTypeID: input.ActivityTypeID,
		Title:          input.Title,
		Description:    input.Description,
		StartsAt:       input.StartsAt,
		EndsAt:         input.EndsAt,
		MaxScore:       input.MaxScore,
		IsPublished:    input.IsPublished,
	}

	if err := repo.Create(activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity"})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

func UpdateActivity(c *gin.Context, repo *repository.ActivityRepository) {
	id := c.Param("id")
	
	existing, err := repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	var req validator.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	input := validator.CreateActivityInput{
		CourseID:       req.CourseID,
		ActivityTypeID: req.ActivityTypeID,
		Title:          req.Title,
		Description:    req.Description,
		StartsAt:       req.StartsAt,
		EndsAt:         req.EndsAt,
		MaxScore:       req.MaxScore,
		IsPublished:    req.IsPublished,
	}

	if err := validate.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.FormatError(err)})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activity := &model.Activity{
		ID:             id,
		CourseID:       input.CourseID,
		ActivityTypeID: input.ActivityTypeID,
		Title:          input.Title,
		Description:    input.Description,
		StartsAt:       input.StartsAt,
		EndsAt:         input.EndsAt,
		MaxScore:       input.MaxScore,
		IsPublished:    input.IsPublished,
	}

	if err := repo.Update(activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update activity"})
		return
	}

	c.JSON(http.StatusOK, activity)
}

func DeleteActivity(c *gin.Context, repo *repository.ActivityRepository) {
	id := c.Param("id")

	existing, err := repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	if err := repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete activity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity deleted"})
}