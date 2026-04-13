package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AxilTH/scout-backend/services/education/internal/model"
	"github.com/AxilTH/scout-backend/services/education/internal/repository"
	"github.com/AxilTH/scout-backend/services/education/internal/validator"
	"github.com/gin-gonic/gin"
)

func SubmitActivity(c *gin.Context, repo *repository.ActivityResultRepository) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication middleware not configured"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in token"})
		return
	}

	activityIDStr := c.Param("id")
	activityID, err := strconv.ParseInt(activityIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID format"})
		return
	}

	// Проверяем, не отправлял ли уже пользователь эту активность
	existing, err := repo.GetByUserAndActivity(userID, activityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if existing != nil && existing.Status == "submitted" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Activity already submitted"})
		return
	}

	// Если результат уже существует, обновляем его
	if existing != nil {
		now := time.Now()
		existing.Status = "submitted"
		existing.SubmittedAt = &now
		existing.Score = nil
		existing.Feedback = nil
		existing.ReviewedAt = nil
		existing.ReviewerID = nil

		if err := repo.Update(existing); err != nil {
			log.Printf("Failed to update activity result: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit activity"})
			return
		}

		c.JSON(http.StatusOK, existing)
		return
	}

	// Создаём новый результат
	result := &model.ActivityResult{
		UserID:      userID,
		ActivityID:  activityID,
		Status:      "submitted",
		SubmittedAt: &[]time.Time{time.Now()}[0],
	}

	if err := repo.Create(result); err != nil {
		log.Printf("Failed to create activity result: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit activity"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func GetActivityResult(c *gin.Context, repo *repository.ActivityResultRepository) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity result ID format"})
		return
	}

	result, err := repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity result not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func ReviewActivityResult(c *gin.Context, repo *repository.ActivityResultRepository) {
	reviewerIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication middleware not configured"})
		return
	}

	reviewerID, err := strconv.ParseInt(reviewerIDStr.(string), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in token"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity result ID format"})
		return
	}

	// Проверяем существование
	existing, err := repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity result not found"})
		return
	}

	// Парсим запрос
	var req validator.ReviewActivityResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Валидация
	input := validator.ReviewActivityResultInput{
		Score:    req.Score,
		Feedback: req.Feedback,
		Status:   req.Status,
	}

	if err := validate.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.FormatError(err)})
		return
	}

	// Обновляем результат
	now := time.Now()
	existing.Status = input.Status
	existing.Score = input.Score
	existing.Feedback = input.Feedback
	existing.ReviewedAt = &now
	existing.ReviewerID = &reviewerID

	if err := repo.Update(existing); err != nil {
		log.Printf("Failed to review activity result: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to review activity"})
		return
	}

	c.JSON(http.StatusOK, existing)
}

func GetCourseResults(c *gin.Context, repo *repository.ActivityResultRepository) {
	courseIDStr := c.Param("id")
	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
		return
	}

	results, err := repo.ListByCourse(courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func GetUserResults(c *gin.Context, repo *repository.ActivityResultRepository) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	results, err := repo.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, results)
}
