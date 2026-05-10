// internal/handler/squad_leadership.go
package handler

import (
	"net/http"
	"time"

	"github.com/AxilTH/scout-backend/services/squad/internal/model"
	"github.com/AxilTH/scout-backend/services/squad/internal/validator"
	"github.com/gin-gonic/gin"
)

// GetSquadLeadership возвращает командный состав отряда с пагинацией.
// GET /leadership
func (h *Handler) GetSquadLeadership(c *gin.Context) {
	// Извлекаем squad_id из контекста (для изоляции данных)
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		RespondInternalError(c, "Failed to get squad information: "+err.Error())
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

// GetUserLeadership возвращает должности пользователя с пагинацией.
// GET /users/:user_id/leadership
// Доступно только если текущий пользователь и запрашиваемый пользователь находятся в одном отряде.
func (h *Handler) GetUserLeadership(c *gin.Context) {
	// Извлекаем user_id из пути
	targetUserID, err := parseIDParam(c, "user_id")
	if err != nil {
		RespondValidationError(c, "Invalid user ID")
		return
	}

	// Извлекаем squad_id текущего пользователя из контекста
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		RespondInternalError(c, "Failed to get squad information: "+err.Error())
		return
	}

	// Проверяем, что запрашиваемый пользователь является членом того же отряда
	memberships, err := h.squadMembershipRepo.GetBySquadID(c.Request.Context(), squadID, 100, 0)
	if err != nil {
		RespondInternalError(c, "Failed to check membership: "+err.Error())
		return
	}

	isInSameSquad := false
	for _, membership := range memberships {
		if membership.UserID == targetUserID && membership.IsActive {
			isInSameSquad = true
			break
		}
	}

	if !isInSameSquad {
		RespondError(c, http.StatusForbidden, "You can only view leadership positions of users in your squad")
		return
	}

	limit, offset := GetLimitOffset(c)

	leadership, err := h.squadLeadershipRepo.GetByUserID(c.Request.Context(), targetUserID, limit, offset)
	if err != nil {
		RespondInternalError(c, "Failed to get user leadership: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, leadership)
}

// GetMyLeadership возвращает должности текущего пользователя в отряде.
// GET /leadership/me
func (h *Handler) GetMyLeadership(c *gin.Context) {
	// Извлекаем user_id из контекста
	userID, err := getUserIDFromContext(c)
	if err != nil {
		RespondInternalError(c, "Failed to get user information: "+err.Error())
		return
	}

	// Извлекаем squad_id из контекста
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		RespondInternalError(c, "Failed to get squad information: "+err.Error())
		return
	}

	// Получаем все должности пользователя в этом отряде
	leaderships, err := h.squadLeadershipRepo.GetBySquadID(c.Request.Context(), squadID, 100, 0)
	if err != nil {
		RespondInternalError(c, "Failed to get leadership: "+err.Error())
		return
	}

	// Ищем должности текущего пользователя
	var myLeaderships []*model.SquadLeadership
	for _, leadership := range leaderships {
		if leadership.UserID == userID {
			myLeaderships = append(myLeaderships, leadership)
		}
	}

	if len(myLeaderships) == 0 {
		RespondNotFound(c, "Leadership positions not found")
		return
	}

	RespondSuccess(c, http.StatusOK, myLeaderships)
}

// AssignPosition назначает должность пользователю в отряде.
// POST /leadership
func (h *Handler) AssignPosition(c *gin.Context) {
	// Извлекаем squad_id из контекста (для изоляции данных)
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		RespondInternalError(c, "Failed to get squad information: "+err.Error())
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
		SquadID:    squadID, // берем из контекста
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
		SquadID:     squadID,
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

// UpdateLeadershipPosition обновляет должность пользователя в отряде.
// PUT /leadership/:id
func (h *Handler) UpdateLeadershipPosition(c *gin.Context) {
	leadershipID, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid leadership ID")
		return
	}

	// Извлекаем squad_id из контекста (для изоляции данных)
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		RespondInternalError(c, "Failed to get squad information: "+err.Error())
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

	// Получаем текущую должность по ID с проверкой squad_id
	leadership, err := h.squadLeadershipRepo.GetByIDAndSquadID(c.Request.Context(), leadershipID, squadID)
	if err != nil {
		RespondNotFound(c, "Active leadership position not found")
		return
	}

	// Обновляем должность
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
}

// DismissPosition снимает пользователя с должности в отряде.
// DELETE /leadership/:id
func (h *Handler) DismissPosition(c *gin.Context) {
	leadershipID, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid leadership ID")
		return
	}

	// Извлекаем squad_id из контекста (для изоляции данных)
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		RespondInternalError(c, "Failed to get squad information: "+err.Error())
		return
	}

	// Проверяем, что отряд существует
	_, err = h.squadRepo.GetByID(c.Request.Context(), squadID)
	if err != nil {
		RespondNotFound(c, "Squad not found")
		return
	}

	// Получаем текущую должность по ID с проверкой squad_id
	leadership, err := h.squadLeadershipRepo.GetByIDAndSquadID(c.Request.Context(), leadershipID, squadID)
	if err != nil {
		RespondNotFound(c, "Active leadership position not found")
		return
	}

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
}
