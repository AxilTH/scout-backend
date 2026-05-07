// internal/handler/squad_leadership.go
package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AxilTH/scout-backend/services/squad/internal/model"
	"github.com/AxilTH/scout-backend/services/squad/internal/validator"
	"github.com/gin-gonic/gin"
)

// GetSquadLeadership возвращает командный состав отряда с пагинацией
// GET /squads/:id/leadership
func (h *Handler) GetSquadLeadership(c *gin.Context) {
	squadID, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid squad ID")
		return
	}

	// Проверяем, что отряд существует
	_, err = h.squadRepo.GetByID(c.Request.Context(), squadID)
	if err != nil {
		RespondNotFound(c, "Squad not found")
		return
	}

	limit, offset := GetLimitOffset(c)

	leadership, err := h.squadLeadershipRepo.GetBySquadID(c.Request.Context(), squadID, limit, offset)
	if err != nil {
		RespondInternalError(c, "Failed to get squad leadership: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, leadership)
}

// GetUserLeadership возвращает должности пользователя с пагинацией
// GET /users/:user_id/leadership
func (h *Handler) GetUserLeadership(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		RespondValidationError(c, "Invalid user ID")
		return
	}

	limit, offset := GetLimitOffset(c)

	leadership, err := h.squadLeadershipRepo.GetByUserID(c.Request.Context(), userID, limit, offset)
	if err != nil {
		RespondInternalError(c, "Failed to get user leadership: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, leadership)
}

// AssignPosition назначает должность пользователю в отряде
// POST /squads/:id/leadership
func (h *Handler) AssignPosition(c *gin.Context) {
	squadID, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid squad ID")
		return
	}

	// 1. Парсим тело запроса в Request (только для десериализации)
	var req validator.AssignPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, "Invalid JSON format")
		return
	}

	// 2. Преобразуем в Input для валидации
	input := validator.AssignPositionInput{
		UserID:     req.UserID,
		SquadID:    squadID, // берем из пути
		PositionID: req.PositionID,
	}

	// 3. Структурная валидация (теги validate)
	structErrs := validator.ValidateStruct(&input)

	// 4. Бизнес-валидация (логика предметной области)
	bizErrs := input.Validate()

	// 5. Объединяем и проверяем ошибки
	allErrs := append(structErrs, bizErrs...)
	if len(allErrs) > 0 {
		RespondValidationErrors(c, allErrs.ToHumanReadable())
		return
	}

	// Проверяем, что squadID из пути совпадает с squadID из тела запроса
	if req.SquadID != squadID {
		RespondValidationError(c, "Squad ID mismatch")
		return
	}

	// Проверяем, что отряд существует
	_, err = h.squadRepo.GetByID(c.Request.Context(), squadID)
	if err != nil {
		RespondNotFound(c, "Squad not found")
		return
	}

	// Проверяем, что должность существует
	_, err = h.positionRepo.GetByID(c.Request.Context(), req.PositionID)
	if err != nil {
		RespondValidationError(c, "Position not found")
		return
	}

	// Проверяем, что пользователь является членом отряда
	memberships, err := h.squadMembershipRepo.GetBySquadID(c.Request.Context(), squadID, 100, 0)
	if err != nil {
		RespondInternalError(c, "Failed to check membership: "+err.Error())
		return
	}

	isMember := false
	var userRoleID int64
	for _, membership := range memberships {
		if membership.UserID == req.UserID && membership.IsActive {
			isMember = true
			userRoleID = membership.RoleID
			break
		}
	}

	if !isMember {
		RespondValidationError(c, "User is not an active member of this squad")
		return
	}

	// Бизнес-правило: все члены командного состава должны быть fighter
	fighterRole, err := h.roleRepo.GetByTitle(c.Request.Context(), "fighter")
	if err != nil {
		RespondInternalError(c, "Failed to get fighter role: "+err.Error())
		return
	}

	if userRoleID != fighterRole.ID {
		RespondValidationError(c, "Only fighters can be assigned to leadership positions")
		return
	}

	// Проверяем, что у пользователя еще нет активной должности в этом отряде
	existingLeadership, err := h.squadLeadershipRepo.GetBySquadID(c.Request.Context(), squadID, 100, 0)
	if err == nil {
		for _, leadership := range existingLeadership {
			if leadership.UserID == req.UserID && leadership.DismissedAt == nil {
				RespondError(c, http.StatusConflict, "User already has an active position in this squad")
				return
			}
		}
	}

	// Создаем новую должность
	leadership := &model.SquadLeadership{
		UserID:      req.UserID,
		SquadID:     req.SquadID,
		PositionID:  req.PositionID,
		AppointedAt: GetCurrentTime(),
		CreatedAt:   GetCurrentTime(),
		UpdatedAt:   GetCurrentTime(),
	}

	if err := h.squadLeadershipRepo.Create(c.Request.Context(), leadership); err != nil {
		RespondInternalError(c, "Failed to assign position: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusCreated, leadership)
}

// UpdateLeadershipPosition обновляет должность пользователя в отряде
// PUT /squads/:id/leadership/:user_id
func (h *Handler) UpdateLeadershipPosition(c *gin.Context) {
	squadID, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid squad ID")
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		RespondValidationError(c, "Invalid user ID")
		return
	}

	// 1. Парсим тело запроса в Request (только для десериализации)
	var req validator.UpdateLeadershipPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, "Invalid JSON format")
		return
	}

	// 2. Преобразуем в Input для валидации
	input := validator.UpdateLeadershipPositionInput{
		PositionID:  req.PositionID,
		DismissedAt: req.DismissedAt,
	}

	// 3. Структурная валидация (теги validate)
	structErrs := validator.ValidateStruct(&input)

	// 4. Бизнес-валидация (логика предметной области)
	bizErrs := input.Validate()

	// 5. Объединяем и проверяем ошибки
	allErrs := append(structErrs, bizErrs...)
	if len(allErrs) > 0 {
		RespondValidationErrors(c, allErrs.ToHumanReadable())
		return
	}

	// Проверяем, что отряд существует
	_, err = h.squadRepo.GetByID(c.Request.Context(), squadID)
	if err != nil {
		RespondNotFound(c, "Squad not found")
		return
	}

	// Проверяем, что должность существует
	_, err = h.positionRepo.GetByID(c.Request.Context(), req.PositionID)
	if err != nil {
		RespondValidationError(c, "Position not found")
		return
	}

	// Получаем текущую должность пользователя
	leadershipList, err := h.squadLeadershipRepo.GetBySquadID(c.Request.Context(), squadID, 100, 0)
	if err != nil {
		RespondInternalError(c, "Failed to get leadership: "+err.Error())
		return
	}

	// Ищем и обновляем должность пользователя
	for _, leadership := range leadershipList {
		if leadership.UserID == userID && leadership.DismissedAt == nil {
			leadership.PositionID = req.PositionID
			leadership.UpdatedAt = GetCurrentTime()

			// Если указана дата увольнения, устанавливаем её
			if req.DismissedAt != nil {
				dismissedAt := time.Unix(*req.DismissedAt, 0).UTC()
				leadership.DismissedAt = &dismissedAt
			}

			if err := h.squadLeadershipRepo.Update(c.Request.Context(), leadership); err != nil {
				RespondInternalError(c, "Failed to update position: "+err.Error())
				return
			}

			RespondSuccess(c, http.StatusOK, leadership)
			return
		}
	}

	RespondNotFound(c, "Active leadership position not found")
}

// DismissPosition снимает пользователя с должности в отряде
// DELETE /squads/:id/leadership/:user_id
func (h *Handler) DismissPosition(c *gin.Context) {
	squadID, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid squad ID")
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		RespondValidationError(c, "Invalid user ID")
		return
	}

	// Проверяем, что отряд существует
	_, err = h.squadRepo.GetByID(c.Request.Context(), squadID)
	if err != nil {
		RespondNotFound(c, "Squad not found")
		return
	}

	// Получаем текущую должность пользователя
	leadershipList, err := h.squadLeadershipRepo.GetBySquadID(c.Request.Context(), squadID, 100, 0)
	if err != nil {
		RespondInternalError(c, "Failed to get leadership: "+err.Error())
		return
	}

	// Ищем и снимаем пользователя с должности
	for _, leadership := range leadershipList {
		if leadership.UserID == userID && leadership.DismissedAt == nil {
			// Устанавливаем дату увольнения
			dismissedAt := GetCurrentTime()
			leadership.DismissedAt = &dismissedAt
			leadership.UpdatedAt = GetCurrentTime()

			if err := h.squadLeadershipRepo.Update(c.Request.Context(), leadership); err != nil {
				RespondInternalError(c, "Failed to dismiss position: "+err.Error())
				return
			}

			RespondSuccess(c, http.StatusOK, gin.H{
				"message": "Position dismissed successfully",
			})
			return
		}
	}

	RespondNotFound(c, "Active leadership position not found")
}
