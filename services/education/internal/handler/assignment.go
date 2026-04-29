package handler

// import (
// 	"log"
// 	"net/http"
// 	"strconv"

// 	"github.com/AxilTH/scout-backend/services/education/internal/model"
// 	"github.com/AxilTH/scout-backend/services/education/internal/repository"
// 	"github.com/AxilTH/scout-backend/services/education/internal/validator"
// 	"github.com/gin-gonic/gin"
// )

// // GetAssignments обрабатывает получение списка заданий
// func GetAssignments(c *gin.Context, repo *repository.AssignmentRepository) {
// 	activityIDStr := c.Query("activity_id")

// 	if activityIDStr != "" {
// 		activityID, err := strconv.ParseInt(activityIDStr, 10, 64)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity_id format"})
// 			return
// 		}
// 		// Получаем задания для конкретной активности
// 		assignments, err := repo.ListByActivityID(activityID)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
// 			return
// 		}
// 		c.JSON(http.StatusOK, assignments)
// 	} else {
// 		// Получаем все задания
// 		assignments, err := repo.ListAll()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
// 			return
// 		}
// 		c.JSON(http.StatusOK, assignments)
// 	}
// }

// // GetAssignment обрабатывает получение задания по ID
// func GetAssignment(c *gin.Context, repo *repository.AssignmentRepository) {
// 	idStr := c.Param("id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID format"})
// 		return
// 	}

// 	assignment, err := repo.GetByID(id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
// 		return
// 	}
// 	if assignment == nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, assignment)
// }

// // CreateAssignment обрабатывает создание задания
// func CreateAssignment(c *gin.Context, repo *repository.AssignmentRepository) {
// 	userID, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication middleware not configured"})
// 		return
// 	}
// 	log.Printf("User %s is creating an assignment", userID)

// 	var req validator.CreateAssignmentRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
// 		return
// 	}

// 	// Преобразуем в структуру для валидации
// 	input := validator.CreateAssignmentInput{
// 		ActivityID:       req.ActivityID,
// 		AssignmentTypeID: req.AssignmentTypeID,
// 		Title:            req.Title,
// 		Description:      req.Description,
// 		Deadline:         req.Deadline,
// 		MaxScore:         req.MaxScore,
// 		IsPublished:      req.IsPublished,
// 	}

// 	// Структурная валидация
// 	if err := validate.Struct(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": validator.FormatError(err)})
// 		return
// 	}

// 	// Бизнес-валидация
// 	if err := input.Validate(); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Создаём задание
// 	assignment := &model.Assignment{
// 		ActivityID:       input.ActivityID,
// 		AssignmentTypeID: input.AssignmentTypeID,
// 		Title:            input.Title,
// 		Description:      input.Description,
// 		Deadline:         input.Deadline,
// 		MaxScore:         input.MaxScore,
// 		IsPublished:      input.IsPublished,
// 	}

// 	// Сохраняем в БД
// 	if err := repo.Create(assignment); err != nil {
// 		log.Printf("Failed to create assignment: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create assignment"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, assignment)
// }

// // UpdateAssignment обрабатывает обновление задания
// func UpdateAssignment(c *gin.Context, repo *repository.AssignmentRepository) {
// 	idStr := c.Param("id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID format"})
// 		return
// 	}

// 	// Проверяем существование задания
// 	existing, err := repo.GetByID(id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
// 		return
// 	}
// 	if existing == nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
// 		return
// 	}

// 	// Парсим запрос
// 	var req validator.CreateAssignmentRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
// 		return
// 	}

// 	// Преобразуем в структуру для валидации
// 	input := validator.CreateAssignmentInput{
// 		ActivityID:       req.ActivityID,
// 		AssignmentTypeID: req.AssignmentTypeID,
// 		Title:            req.Title,
// 		Description:      req.Description,
// 		Deadline:         req.Deadline,
// 		MaxScore:         req.MaxScore,
// 		IsPublished:      req.IsPublished,
// 	}

// 	// Валидация
// 	if err := validate.Struct(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": validator.FormatError(err)})
// 		return
// 	}
// 	if err := input.Validate(); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Обновляем задание
// 	assignment := &model.Assignment{
// 		ID:               id,
// 		ActivityID:       input.ActivityID,
// 		AssignmentTypeID: input.AssignmentTypeID,
// 		Title:            input.Title,
// 		Description:      input.Description,
// 		Deadline:         input.Deadline,
// 		MaxScore:         input.MaxScore,
// 		IsPublished:      input.IsPublished,
// 	}

// 	// Сохраняем в БД
// 	if err := repo.Update(assignment); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update assignment"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, assignment)
// }

// // DeleteAssignment обрабатывает удаление задания
// func DeleteAssignment(c *gin.Context, repo *repository.AssignmentRepository) {
// 	idStr := c.Param("id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID format"})
// 		return
// 	}

// 	// Проверяем существование задания
// 	existing, err := repo.GetByID(id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
// 		return
// 	}
// 	if existing == nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
// 		return
// 	}

// 	// Удаляем задание
// 	if err := repo.Delete(id); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete assignment"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Assignment deleted"})
// }
