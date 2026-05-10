// internal/handler/squad_membership.go
package handler

import (
	"net/http"

	"github.com/AxilTH/scout-backend/services/squad/internal/model"
	"github.com/AxilTH/scout-backend/services/squad/internal/validator"
	"github.com/gin-gonic/gin"
)

// AddMembershipRequest представляет запрос на добавление участника в отряд
type AddMembershipRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
	RoleID int64 `json:"role_id" binding:"required"`
}

// UpdateMembershipRequest представляет запрос на обновление членства
type UpdateMembershipRequest struct {
	RoleID   int64  `json:"role_id" binding:"required"`
	IsActive bool   `json:"is_active"`
	LeftAt   *int64 `json:"left_at,omitempty"` // timestamp
}

// GetSquadMembers возвращает всех участников отряда с пагинацией.
// GET /members
func (h *Handler) GetSquadMembers(c *gin.Context) {
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

	memberships, err := h.squadMembershipRepo.GetBySquadID(c.Request.Context(), squadID, limit, offset)
	if err != nil {
		RespondInternalError(c, "Failed to get squad members: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, memberships)
}

// GetActiveSquadMembers возвращает только активных участников отряда с пагинацией.
// GET /members/active
func (h *Handler) GetActiveSquadMembers(c *gin.Context) {
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

	memberships, err := h.squadMembershipRepo.GetActiveBySquadID(c.Request.Context(), squadID, limit, offset)
	if err != nil {
		RespondInternalError(c, "Failed to get active squad members: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, memberships)
}

// AddMembership добавляет участника в отряд.
// POST /memberships
func (h *Handler) AddMembership(c *gin.Context) {
	// Извлекаем squad_id из контекста (для изоляции данных)
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		RespondInternalError(c, "Failed to get squad information: "+err.Error())
		return
	}

	// 1. Парсим тело запроса в Request (только для десериализации)
	var req AddMembershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, "Invalid JSON format")
		return
	}

	// 2. Преобразуем в Input для валидации
	input := validator.AddMembershipInput{
		UserID:  req.UserID,
		SquadID: squadID, // берем из контекста
		RoleID:  req.RoleID,
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

	// Проверяем, что роль существует
	_, err = h.roleRepo.GetByID(c.Request.Context(), req.RoleID)
	if err != nil {
		RespondValidationError(c, "Role not found")
		return
	}

	// Проверяем, что у пользователя еще нет активного членства в этом отряде
	existingMembership, err := h.squadMembershipRepo.GetActiveByUserID(c.Request.Context(), req.UserID)
	if err == nil && existingMembership != nil && existingMembership.SquadID == squadID {
		RespondError(c, http.StatusConflict, "User already has active membership in this squad")
		return
	}

	// Создаем новое членство
	membership := &model.SquadMembership{
		UserID:    req.UserID,
		SquadID:   squadID,
		RoleID:    req.RoleID,
		IsActive:  true,
		JoinedAt:  GetCurrentTime(),
		CreatedAt: GetCurrentTime(),
		UpdatedAt: GetCurrentTime(),
	}

	if err := h.squadMembershipRepo.Create(c.Request.Context(), membership); err != nil {
		RespondInternalError(c, "Failed to add membership: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusCreated, membership)
}

// GetUserRole возвращает роль пользователя в отряде.
// GET /members/:id/role
func (h *Handler) GetUserRole(c *gin.Context) {
	membershipID, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid membership ID")
		return
	}

	// Извлекаем squad_id из контекста (для изоляции данных)
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		RespondInternalError(c, "Failed to get squad information: "+err.Error())
		return
	}

	// Получаем членство по ID с проверкой squad_id
	membership, err := h.squadMembershipRepo.GetByIDAndSquadID(c.Request.Context(), membershipID, squadID)
	if err != nil {
		RespondNotFound(c, "Membership not found")
		return
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"user_id":   membership.UserID,
		"squad_id":  membership.SquadID,
		"role_id":   membership.RoleID,
		"is_active": membership.IsActive,
	})
}

// UpdateMembership обновляет членство пользователя в отряде.
// PUT /members/:id
func (h *Handler) UpdateMembership(c *gin.Context) {
	membershipID, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid membership ID")
		return
	}

	// Извлекаем squad_id из контекста (для изоляции данных)
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		RespondInternalError(c, "Failed to get squad information: "+err.Error())
		return
	}

	// 1. Парсим тело запроса в Request (только для десериализации)
	var req UpdateMembershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, "Invalid JSON format")
		return
	}

	// 2. Преобразуем в Input для валидации
	input := validator.UpdateMembershipInput{
		RoleID:   req.RoleID,
		IsActive: req.IsActive,
		LeftAt:   req.LeftAt,
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

	// Проверяем, что роль существует
	_, err = h.roleRepo.GetByID(c.Request.Context(), req.RoleID)
	if err != nil {
		RespondValidationError(c, "Role not found")
		return
	}

	// Получаем членство по ID с проверкой squad_id
	membership, err := h.squadMembershipRepo.GetByIDAndSquadID(c.Request.Context(), membershipID, squadID)
	if err != nil {
		RespondNotFound(c, "Membership not found")
		return
	}

	// Обновляем членство
	membership.RoleID = req.RoleID
	membership.IsActive = req.IsActive
	membership.UpdatedAt = GetCurrentTime()

	// Если деактивируем членство, устанавливаем дату ухода
	if !req.IsActive && membership.LeftAt == nil {
		leftAt := GetCurrentTime()
		membership.LeftAt = &leftAt
	}

	if err := h.squadMembershipRepo.Update(c.Request.Context(), membership); err != nil {
		RespondInternalError(c, "Failed to update membership: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, membership)
}

// RemoveMembership удаляет членство пользователя из отряда.
// DELETE /members/:id
func (h *Handler) RemoveMembership(c *gin.Context) {
	membershipID, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid membership ID")
		return
	}

	// Извлекаем squad_id из контекста (для изоляции данных)
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		RespondInternalError(c, "Failed to get squad information: "+err.Error())
		return
	}

	// Получаем членство по ID с проверкой squad_id
	membership, err := h.squadMembershipRepo.GetByIDAndSquadID(c.Request.Context(), membershipID, squadID)
	if err != nil {
		RespondNotFound(c, "Membership not found")
		return
	}

	if err := h.squadMembershipRepo.Delete(c.Request.Context(), membership.ID); err != nil {
		RespondInternalError(c, "Failed to remove membership: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"message": "Membership removed successfully",
	})
}

// GetMyMembership возвращает членство текущего пользователя в отряде.
// GET /members/me
func (h *Handler) GetMyMembership(c *gin.Context) {
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

	// Получаем членство пользователя в этом отряде
	memberships, err := h.squadMembershipRepo.GetBySquadID(c.Request.Context(), squadID, 100, 0)
	if err != nil {
		RespondInternalError(c, "Failed to get membership: "+err.Error())
		return
	}

	// Ищем членство текущего пользователя
	for _, membership := range memberships {
		if membership.UserID == userID {
			RespondSuccess(c, http.StatusOK, membership)
			return
		}
	}

	RespondNotFound(c, "Membership not found")
}

// GetUserSquads возвращает все отряды текущего пользователя с пагинацией.
// GET /users/me/squads
func (h *Handler) GetUserSquads(c *gin.Context) {
	// Извлекаем user_id из контекста
	userID, err := getUserIDFromContext(c)
	if err != nil {
		RespondInternalError(c, "Failed to get user information: "+err.Error())
		return
	}

	limit, offset := GetLimitOffset(c)

	memberships, err := h.squadMembershipRepo.GetByUserID(c.Request.Context(), userID, limit, offset)
	if err != nil {
		RespondInternalError(c, "Failed to get user squads: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, memberships)
}

// GetUserActiveSquad возвращает активный отряд текущего пользователя.
// GET /users/me/squads/active
func (h *Handler) GetUserActiveSquad(c *gin.Context) {
	// Извлекаем user_id из контекста
	userID, err := getUserIDFromContext(c)
	if err != nil {
		RespondInternalError(c, "Failed to get user information: "+err.Error())
		return
	}

	membership, err := h.squadMembershipRepo.GetActiveByUserID(c.Request.Context(), userID)
	if err != nil {
		RespondNotFound(c, "Active squad not found")
		return
	}

	RespondSuccess(c, http.StatusOK, membership)
}
