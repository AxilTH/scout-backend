// internal/handler/course.go
package handler

import (
	"log"
	"net/http"

	"github.com/AxilTH/scout-backend/services/education/internal/model"
	"github.com/AxilTH/scout-backend/services/education/internal/repository"
	"github.com/AxilTH/scout-backend/services/education/internal/validator"
	"github.com/gin-gonic/gin"
)

// GetCourses возвращает все курсы, принадлежащие указанному отряду
func GetCourses(c *gin.Context, repo *repository.CourseRepository) {
	// Извлекаем squad_id из контекста с использованием helper-функции
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		log.Printf("Error getting squad_id from context: %v", err)
		respondWithInternalError(c, "Failed to get squad information")
		return
	}

	courses, err := repo.GetAll(squadID)
	if err != nil {
		log.Printf("Database error in GetCourses: %v", err)
		respondWithInternalError(c, "Failed to fetch courses")
		return
	}

	respondWithSuccess(c, http.StatusOK, courses)
}

// GetCourse возвращает курс по ID, только если он принадлежит отряду пользователя
func GetCourse(c *gin.Context, repo *repository.CourseRepository) {
	// Извлекаем squad_id из контекста
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		log.Printf("Error getting squad_id from context: %v", err)
		respondWithInternalError(c, "Failed to get squad information")
		return
	}

	// Парсим id курса из пути с использованием helper-функции
	id, err := parseIDParam(c, "id")
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Запрос в репозиторий с проверкой принадлежности к отряду
	course, err := repo.GetByID(id, squadID)
	if err != nil {
		log.Printf("Database error in GetCourse: %v", err)
		respondWithInternalError(c, "Failed to fetch course")
		return
	}

	// Курс не найден ИЛИ не принадлежит отряду
	if course == nil {
		respondWithNotFound(c)
		return
	}

	respondWithSuccess(c, http.StatusOK, course)
}

// CreateCourse создает новый курс для отряда, к которому принадлежит пользователь
func CreateCourse(c *gin.Context, repo *repository.CourseRepository) {
	// 1. Извлекаем squad_id из контекста
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		log.Printf("Error getting squad_id from context: %v", err)
		respondWithInternalError(c, "Failed to get squad information")
		return
	}

	// 2. Парсим тело запроса в Request (только для десериализации)
	var req validator.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// 3. Преобразуем в Input для валидации
	input := validator.CreateCourseInput{
		Title:       req.Title,
		Description: req.Description,
		StartsAt:    req.StartsAt,
		EndsAt:      req.EndsAt,
		Year:        req.Year,
		SquadID:     squadID, // берем из контекста
	}

	// 4. Структурная валидация (теги validate)
	structErrs := validator.ValidateStruct(&input)

	// 5. Бизнес-валидация (логика предметной области)
	bizErrs := input.Validate()

	// 6. Объединяем и проверяем ошибки
	allErrs := append(structErrs, bizErrs...)
	if len(allErrs) > 0 {
		respondWithValidationError(c, allErrs.ToHumanReadable())
		return
	}

	// 7. Создаём модель курса
	course := &model.Course{
		Title:       input.Title,
		Description: input.Description,
		StartsAt:    input.StartsAt,
		EndsAt:      input.EndsAt,
		Year:        input.Year,
		IsActive:    true, // по умолчанию новый курс активен
		SquadID:     input.SquadID,
		// CreatedAt/UpdatedAt будут установлены в репозитории
	}

	// 8. Сохраняем в БД
	if err := repo.Create(course); err != nil {
		log.Printf("Failed to create course: %v", err)
		respondWithInternalError(c, "Failed to create course")
		return
	}

	respondWithCreated(c, course)
}

// UpdateCourse обновляет существующий курс (только если он принадлежит отряду)
func UpdateCourse(c *gin.Context, repo *repository.CourseRepository) {
	// Извлекаем squad_id из контекста
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		log.Printf("Error getting squad_id from context: %v", err)
		respondWithInternalError(c, "Failed to get squad information")
		return
	}

	// Парсим id курса из пути
	id, err := parseIDParam(c, "id")
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Парсим тело запроса
	var req validator.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Преобразуем в Input для валидации
	input := validator.CreateCourseInput{
		Title:       req.Title,
		Description: req.Description,
		StartsAt:    req.StartsAt,
		EndsAt:      req.EndsAt,
		Year:        req.Year,
		SquadID:     squadID,
	}

	// Структурная валидация
	structErrs := validator.ValidateStruct(&input)
	bizErrs := input.Validate()
	allErrs := append(structErrs, bizErrs...)
	if len(allErrs) > 0 {
		respondWithValidationError(c, allErrs.ToHumanReadable())
		return
	}

	// Проверяем, что курс существует и принадлежит отряду
	existing, err := repo.GetByID(id, squadID)
	if err != nil {
		log.Printf("Database error in UpdateCourse (get): %v", err)
		respondWithInternalError(c, "Failed to update course")
		return
	}
	if existing == nil {
		respondWithNotFound(c)
		return
	}

	// Обновляем модель
	updated := &model.Course{
		ID:          id,
		Title:       input.Title,
		Description: input.Description,
		StartsAt:    input.StartsAt,
		EndsAt:      input.EndsAt,
		Year:        input.Year,
		IsActive:    existing.IsActive, // не меняем статус здесь (отдельный эндпоинт)
		SquadID:     input.SquadID,     // не даем изменить принадлежность
	}

	if err := repo.Update(updated); err != nil {
		log.Printf("Failed to update course: %v", err)
		respondWithInternalError(c, "Failed to update course")
		return
	}

	respondWithSuccess(c, http.StatusOK, updated)
}

// DeleteCourse удаляет курс (только если он принадлежит отряду)
func DeleteCourse(c *gin.Context, repo *repository.CourseRepository) {
	// Извлекаем squad_id из контекста
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		log.Printf("Error getting squad_id from context: %v", err)
		respondWithInternalError(c, "Failed to get squad information")
		return
	}

	// Парсим id курса из пути
	id, err := parseIDParam(c, "id")
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := repo.Delete(id, squadID); err != nil {
		log.Printf("Failed to delete course: %v", err)
		respondWithInternalError(c, "Failed to delete course")
		return
	}

	respondWithNoContent(c)
}