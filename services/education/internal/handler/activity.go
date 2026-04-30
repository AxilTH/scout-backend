// internal/handler/activity.go
package handler

import (
	"log"
	"net/http"

	"github.com/AxilTH/scout-backend/services/education/internal/model"
	"github.com/AxilTH/scout-backend/services/education/internal/repository"
	"github.com/AxilTH/scout-backend/services/education/internal/validator"
	"github.com/gin-gonic/gin"
)

// GetActivities возвращает все активности указанного курса отряда пользователя
func GetActivities(c *gin.Context, repo *repository.ActivityRepository) {
	// Извлекаем squad_id из контекста
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		log.Printf("Error getting squad_id from context: %v", err)
		respondWithInternalError(c, "Failed to get squad information")
		return
	}

	// Парсим course_id из пути
	courseID, err := parseIDParam(c, "id")
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Запрос в репозиторий с проверкой принадлежности курса к отряду
	activities, err := repo.GetAll(courseID, squadID)
	if err != nil {
		log.Printf("Database error in GetActivities: %v", err)
		respondWithInternalError(c, "Failed to fetch activities")
		return
	}

	respondWithSuccess(c, http.StatusOK, activities)
}

// GetActivity возвращает активность по ID, только если она принадлежит курсу и курс принадлежит отряду
func GetActivity(c *gin.Context, repo *repository.ActivityRepository) {
	// Извлекаем squad_id из контекста
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		log.Printf("Error getting squad_id from context: %v", err)
		respondWithInternalError(c, "Failed to get squad information")
		return
	}

	// Парсим id активности из пути
	activityID, err := parseIDParam(c, "activity_id")
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Парсим course_id из пути
	courseID, err := parseIDParam(c, "id")
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Запрос в репозиторий с проверкой принадлежности к курсу и отряду
	activity, err := repo.GetByID(activityID, courseID, squadID)
	if err != nil {
		log.Printf("Database error in GetActivity: %v", err)
		respondWithInternalError(c, "Failed to fetch activity")
		return
	}

	// Активность не найдена ИЛИ не принадлежит курсу/отряду
	if activity == nil {
		respondWithNotFound(c)
		return
	}

	respondWithSuccess(c, http.StatusOK, activity)
}

// CreateActivity создаёт новую активность в курсе
func CreateActivity(c *gin.Context, activityRepo *repository.ActivityRepository, courseRepo *repository.CourseRepository) {
	// 1. Извлекаем squad_id из контекста
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		log.Printf("Error getting squad_id from context: %v", err)
		respondWithInternalError(c, "Failed to get squad information")
		return
	}

	// 2. Парсим course_id из пути
	courseID, err := parseIDParam(c, "id")
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Парсим тело запроса в Request (только для десериализации)
	var req validator.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// 4. Преобразуем в Input для валидации
	input := validator.CreateActivityInput{
		Title:          req.Title,
		Description:    req.Description,
		StartsAt:       req.StartsAt,
		EndsAt:         req.EndsAt,
		MaxScore:       req.MaxScore,
		IsPublished:    req.IsPublished,
		CourseID:       courseID,
		ActivityTypeID: req.ActivityTypeID,
	}

	// 4. Структурная валидация (теги validate)
	structErrs := validator.ValidateStruct(&input)

	// 5. Бизнес-валидация (StartsAt < EndsAt и другие правила)
	bizErrs := input.Validate()

	// 6. Объединяем и проверяем ошибки
	allErrs := append(structErrs, bizErrs...)
	if len(allErrs) > 0 {
		respondWithValidationError(c, allErrs.ToHumanReadable())
		return
	}

	// 6. Проверка, принадлежит ли курс этому отряду?
	// Запрашиваем курс из БД с проверкой squad_id
	exists, err := courseRepo.ExistsForSquad(courseID, squadID)
	if err != nil {
		log.Printf("Database error checking course ownership: %v", err)
		respondWithInternalError(c, "Failed to create activity")
		return
	}
	if !exists {
		respondWithForbidden(c, "Course not found or access denied")
		return
	}

	// 7. Создаём модель активности
	// CourseID и ActivityTypeID беруте из запроса пользователя
	// CreatedAt и UpdatedAt устанавливаются в репозитории
	activity := &model.Activity{
		Title:          input.Title,
		Description:    input.Description,
		StartsAt:       input.StartsAt,
		EndsAt:         input.EndsAt,
		MaxScore:       input.MaxScore,
		IsPublished:    input.IsPublished,
		CourseID:       input.CourseID,
		ActivityTypeID: input.ActivityTypeID,
		// CreatedAt/UpdatedAt будут установлены в репозитории
	}

	// 8. Сохраняем в БД
	if err := activityRepo.Create(activity); err != nil {
		log.Printf("Failed to create activity: %v", err)
		respondWithInternalError(c, "Failed to create activity")
		return
	}

	respondWithCreated(c, activity)
}

// UpdateActivity обновляет активность
func UpdateActivity(c *gin.Context, activityRepo *repository.ActivityRepository) {
	// 1. Извлекаем squad_id из контекста
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		log.Printf("Error getting squad_id from context: %v", err)
		respondWithInternalError(c, "Failed to get squad information")
		return
	}

	// 2. Парсим id активности из пути
	activityID, err := parseIDParam(c, "activity_id")
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Парсим course_id из пути
	courseID, err := parseIDParam(c, "id")
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 4. Парсим тело запроса
	var req validator.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// 5. Преобразуем в Input для валидации
	input := validator.CreateActivityInput{
		Title:          req.Title,
		Description:    req.Description,
		StartsAt:       req.StartsAt,
		EndsAt:         req.EndsAt,
		MaxScore:       req.MaxScore,
		IsPublished:    req.IsPublished,
		CourseID:       courseID,
		ActivityTypeID: req.ActivityTypeID,
	}

	// 6. Структурная валидация
	structErrs := validator.ValidateStruct(&input)
	if len(structErrs) > 0 {
		respondWithValidationError(c, structErrs.ToHumanReadable())
		return
	}

	// 7. Бизнес-валидация
	bizErrs := input.Validate()
	if len(bizErrs) > 0 {
		respondWithValidationError(c, bizErrs.ToHumanReadable())
		return
	}

	// 8. Проверяем, что активность существует и принадлежит курсу и отряду
	existing, err := activityRepo.GetByID(activityID, courseID, squadID)
	if err != nil {
		log.Printf("Database error in UpdateActivity (get): %v", err)
		respondWithInternalError(c, "Failed to update activity")
		return
	}
	if existing == nil {
		respondWithNotFound(c)
		return
	}

	// 9. Убедимся что не пытаемся изменить привязку к курсу
	// (активность всегда привязана к одному курсу)
	if input.CourseID != existing.CourseID {
		respondWithError(c, http.StatusBadRequest, "Cannot change course_id for existing activity")
		return
	}

	// 10. Обновляем модель
	updated := &model.Activity{
		ID:             activityID,
		Title:          input.Title,
		Description:    input.Description,
		StartsAt:       input.StartsAt,
		EndsAt:         input.EndsAt,
		MaxScore:       input.MaxScore,
		IsPublished:    input.IsPublished,
		CourseID:       input.CourseID,
		ActivityTypeID: input.ActivityTypeID,
	}

	// 11. Сохраняем в БД
	if err := activityRepo.Update(updated, squadID); err != nil {
		log.Printf("Failed to update activity: %v", err)
		respondWithInternalError(c, "Failed to update activity")
		return
	}

	respondWithSuccess(c, http.StatusOK, updated)
}

// DeleteActivity удаляет активность
func DeleteActivity(c *gin.Context, activityRepo *repository.ActivityRepository) {
	// 1. Извлекаем squad_id из контекста
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		log.Printf("Error getting squad_id from context: %v", err)
		respondWithInternalError(c, "Failed to get squad information")
		return
	}

	// 2. Парсим id активности из пути
	activityID, err := parseIDParam(c, "activity_id")
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Парсим course_id из пути
	courseID, err := parseIDParam(c, "id")
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 4. Проверяем, что активность существует и принадлежит курсу и отряду
	existing, err := activityRepo.GetByID(activityID, courseID, squadID)
	if err != nil {
		log.Printf("Database error in DeleteActivity (get): %v", err)
		respondWithInternalError(c, "Failed to delete activity")
		return
	}
	if existing == nil {
		respondWithNotFound(c)
		return
	}

	// 5. Удаляем активность из БД
	if err := activityRepo.Delete(activityID, courseID, squadID); err != nil {
		log.Printf("Failed to delete activity: %v", err)
		respondWithInternalError(c, "Failed to delete activity")
		return
	}

	respondWithNoContent(c)
}
