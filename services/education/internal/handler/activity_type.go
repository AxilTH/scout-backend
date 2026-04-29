// internal/handler/activity_type.go
package handler

import (
	"log"
	"net/http"

	"github.com/AxilTH/scout-backend/services/education/internal/model"
	"github.com/AxilTH/scout-backend/services/education/internal/repository"
	"github.com/AxilTH/scout-backend/services/education/internal/validator"
	"github.com/gin-gonic/gin"
)

// GetActivityTypes возвращает все типы активностей, доступные указанному отряду
func GetActivityTypes(c *gin.Context, repo *repository.ActivityTypeRepository) {
	// Извлекаем squad_id из контекста
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		log.Printf("Error getting squad_id from context: %v", err)
		respondWithInternalError(c, "Failed to get squad information")
		return
	}

	activityTypes, err := repo.GetAll(squadID)
	if err != nil {
		log.Printf("Database error in GetActivityTypes: %v", err)
		respondWithInternalError(c, "Failed to fetch activity types")
		return
	}

	respondWithSuccess(c, http.StatusOK, activityTypes)
}

// GetActivityType возвращает тип активности по ID, только если он доступен отряду пользователя
func GetActivityType(c *gin.Context, repo *repository.ActivityTypeRepository) {
	// Извлекаем squad_id из контекста
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		log.Printf("Error getting squad_id from context: %v", err)
		respondWithInternalError(c, "Failed to get squad information")
		return
	}

	// Парсим id типа активности из пути
	id, err := parseIDParam(c, "id")
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Запрос в репозиторий с проверкой принадлежности к отряду
	activityType, err := repo.GetByID(id, squadID)
	if err != nil {
		log.Printf("Database error in GetActivityType: %v", err)
		respondWithInternalError(c, "Failed to fetch activity type")
		return
	}

	// Активность не найдена ИЛИ не принадлежит отряду
	if activityType == nil {
		respondWithNotFound(c)
		return
	}

	respondWithSuccess(c, http.StatusOK, activityType)
}

// CreateActivityType создаёт новый тип активности (локальный для отряда)
func CreateActivityType(c *gin.Context, repo *repository.ActivityTypeRepository) {
	// 1. Извлекаем squad_id из контекста
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		log.Printf("Error getting squad_id from context: %v", err)
		respondWithInternalError(c, "Failed to get squad information")
		return
	}

	// 2. Парсим тело запроса в Request (только для десериализации)
	var req validator.CreateActivityTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// 3. Преобразуем в Input для валидации
	input := validator.CreateActivityTypeInput{
		Title:       req.Title,
		Description: req.Description,
		SquadID:     squadID, // берем из контекста
	}

	// 4. Структурная валидация (теги validate)
	structErrs := validator.ValidateStruct(&input)

	// 5. Проверяем ошибки
	if len(structErrs) > 0 {
		respondWithValidationError(c, structErrs.ToHumanReadable())
		return
	}

	// 6. Создаём модель типа активности
	// SquadID устанавливается здесь
	// CreatedAt и UpdatedAt устанавливаются в репозитории
	activityType := &model.ActivityType{
		Title:       input.Title,
		Description: input.Description,
		SquadID:     &input.SquadID, // берем из контекста
		// CreatedAt/UpdatedAt будут установлены в репозитории
	}

	// 7. Сохраняем в БД
	if err := repo.Create(activityType); err != nil {
		log.Printf("Failed to create activity type: %v", err)
		respondWithInternalError(c, "Failed to create activity type")
		return
	}

	respondWithCreated(c, activityType)
}

// UpdateActivityType обновляет локальный тип активности
func UpdateActivityType(c *gin.Context, repo *repository.ActivityTypeRepository) {
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		log.Printf("Error getting squad_id from context: %v", err)
		respondWithInternalError(c, "Failed to get squad information")
		return
	}

	id, err := parseIDParam(c, "id")
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	var req validator.CreateActivityTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	input := validator.CreateActivityTypeInput{
		Title:       req.Title,
		Description: req.Description,
		SquadID:     squadID,
	}

	structErrs := validator.ValidateStruct(&input)
	if len(structErrs) > 0 {
		respondWithValidationError(c, structErrs.ToHumanReadable())
		return
	}

	// Проверяем, что тип существует и принадлежит отряду
	existing, err := repo.GetByID(id, squadID)
	if err != nil {
		log.Printf("Database error in UpdateActivityType (get): %v", err)
		respondWithInternalError(c, "Failed to update activity type")
		return
	}
	if existing == nil {
		respondWithNotFound(c)
		return
	}
	// Проверка что это локальный тип (не глобальный)
	if existing.SquadID == nil {
		respondWithForbidden(c, "Cannot modify global activity types")
		return
	}

	// Обновляем модель
	updated := &model.ActivityType{
		ID:          id,
		Title:       input.Title,
		Description: input.Description,
		SquadID:     &input.SquadID, // не даем изменить принадлежность
	}

	if err := repo.Update(updated); err != nil {
		log.Printf("Failed to update activity type: %v", err)
		respondWithInternalError(c, "Failed to update activity type")
		return
	}

	respondWithSuccess(c, http.StatusOK, updated)
}

// DeleteActivityType удаляет тип активности
func DeleteActivityType(c *gin.Context, repo *repository.ActivityTypeRepository) {
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		log.Printf("Error getting squad_id from context: %v", err)
		respondWithInternalError(c, "Failed to get squad information")
		return
	}

	id, err := parseIDParam(c, "id")
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Проверяем, что тип существует и принадлежит отряду
	existing, err := repo.GetByID(id, squadID)
	if err != nil {
		log.Printf("Database error in DeleteActivityType (get): %v", err)
		respondWithInternalError(c, "Failed to delete activity type")
		return
	}
	if existing == nil {
		respondWithNotFound(c)
		return
	}

	// Проверка что это локальный тип (не глобальный)
	if existing.SquadID == nil {
		respondWithForbidden(c, "Cannot delete global activity type")
		return
	}

	if err := repo.Delete(id, squadID); err != nil {
		log.Printf("Failed to delete activity type: %v", err)
		respondWithInternalError(c, "Failed to delete activity type")
		return
	}

	respondWithNoContent(c)
}
