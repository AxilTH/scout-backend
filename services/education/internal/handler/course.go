package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/AxilTH/scout-backend/services/education/internal/model"
	"github.com/AxilTH/scout-backend/services/education/internal/repository"
	"github.com/AxilTH/scout-backend/services/education/internal/validator"
	"github.com/gin-gonic/gin"
)

func GetCourses(c *gin.Context, repo *repository.CourseRepository) {
	squadIDStr := c.Query("squad_id")
	if squadIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "squad_id query parameter is required"})
		return
	}

	squadID, err := strconv.ParseInt(squadIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid squad_id format"})
		return
	}

	courses, err := repo.ListBySquad(squadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusOK, courses)
}

func GetCourse(c *gin.Context, repo *repository.CourseRepository) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
		return
	}

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
	// Получаем user_id из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication middleware not configured"})
		return
	}
	log.Printf("User %s is creating a course", userID)

	// Парсим JSON в Request
	var req validator.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Преобразуем в Input
	input := validator.CreateCourseInput{
		SquadID:     req.SquadID,
		Year:        req.Year,
		Title:       req.Title,
		Description: req.Description,
		StartsAt:    req.StartsAt,
		EndsAt:      req.EndsAt,
	}

	// Структурная валидация
	if err := validate.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.FormatError(err)})
		return
	}

	// Бизнес-валидация
	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создаем модель и сохраняем в БД
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
