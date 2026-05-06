// internal/handler/squad_membership.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/AxilTH/scout-backend/services/squad/internal/model"
	"github.com/gin-gonic/gin"
)

// AddMembershipRequest представляет запрос на добавление участника в отряд
type AddMembershipRequest struct {
	UserID  int64 `json:"user_id" binding:"required"`
	SquadID int64 `json:"squad_id" binding:"required"`
	RoleID  int64 `json:"role_id" binding:"required"`
}

// UpdateMembershipRequest представляет запрос на обновление членства
type UpdateMembershipRequest struct {
	RoleID   int64  `json:"role_id" binding:"required"`
	IsActive bool   `json:"is_active"`
	LeftAt   *int64 `json:"left_at,omitempty"` // timestamp
}

// GetSquadMembers возвращает всех участников отряда с пагинацией
// GET /squads/:id/members
func (h *Handler) GetSquadMembers(c *gin.Context) {
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

	memberships, err := h.squadMembershipRepo.GetBySquadID(c.Request.Context(), squadID, limit, offset)
	if err != nil {
		RespondInternalError(c, "Failed to get squad members: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, memberships)
}

// GetActiveSquadMembers возвращает только активных участников отряда с пагинацией
// GET /squads/:id/members/active
func (h *Handler) GetActiveSquadMembers(c *gin.Context) {
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

	memberships, err := h.squadMembershipRepo.GetActiveBySquadID(c.Request.Context(), squadID, limit, offset)
	if err != nil {
		RespondInternalError(c, "Failed to get active squad members: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, memberships)
}

// AddMembership добавляет участника в отряд
// POST /squads/:id/memberships
func (h *Handler) AddMembership(c *gin.Context) {
	squadID, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid squad ID")
		return
	}

	var req AddMembershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, "Invalid request body: "+err.Error())
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
		SquadID:   req.SquadID,
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

// GetUserRole возвращает роль пользователя в отряде
// GET /squads/:id/members/:user_id/role
func (h *Handler) GetUserRole(c *gin.Context) {
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

	// Получаем членство пользователя
	memberships, err := h.squadMembershipRepo.GetBySquadID(c.Request.Context(), squadID, 100, 0)
	if err != nil {
		RespondInternalError(c, "Failed to get membership: "+err.Error())
		return
	}

	// Ищем членство пользователя
	for _, membership := range memberships {
		if membership.UserID == userID {
			RespondSuccess(c, http.StatusOK, gin.H{
				"user_id":   membership.UserID,
				"squad_id":  membership.SquadID,
				"role_id":   membership.RoleID,
				"is_active": membership.IsActive,
			})
			return
		}
	}

	RespondNotFound(c, "Membership not found")
}

// UpdateMembership обновляет членство пользователя в отряде
// PUT /squads/:id/members/:user_id
func (h *Handler) UpdateMembership(c *gin.Context) {
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

	var req UpdateMembershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, "Invalid request body: "+err.Error())
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

	// Получаем членство пользователя
	memberships, err := h.squadMembershipRepo.GetBySquadID(c.Request.Context(), squadID, 100, 0)
	if err != nil {
		RespondInternalError(c, "Failed to get membership: "+err.Error())
		return
	}

	// Ищем и обновляем членство пользователя
	for _, membership := range memberships {
		if membership.UserID == userID {
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
			return
		}
	}

	RespondNotFound(c, "Membership not found")
}

// RemoveMembership удаляет членство пользователя из отряда
// DELETE /squads/:id/members/:user_id
func (h *Handler) RemoveMembership(c *gin.Context) {
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

	// Получаем членство пользователя
	memberships, err := h.squadMembershipRepo.GetBySquadID(c.Request.Context(), squadID, 100, 0)
	if err != nil {
		RespondInternalError(c, "Failed to get membership: "+err.Error())
		return
	}

	// Ищем и удаляем членство пользователя
	for _, membership := range memberships {
		if membership.UserID == userID {
			if err := h.squadMembershipRepo.Delete(c.Request.Context(), membership.ID); err != nil {
				RespondInternalError(c, "Failed to remove membership: "+err.Error())
				return
			}

			RespondSuccess(c, http.StatusOK, gin.H{
				"message": "Membership removed successfully",
			})
			return
		}
	}

	RespondNotFound(c, "Membership not found")
}

// GetUserSquads возвращает все отряды пользователя с пагинацией
// GET /users/:user_id/squads
func (h *Handler) GetUserSquads(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		RespondValidationError(c, "Invalid user ID")
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

// GetUserActiveSquad возвращает активный отряд пользователя
// GET /users/:user_id/squads/active
func (h *Handler) GetUserActiveSquad(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		RespondValidationError(c, "Invalid user ID")
		return
	}

	membership, err := h.squadMembershipRepo.GetActiveByUserID(c.Request.Context(), userID)
	if err != nil {
		RespondNotFound(c, "Active squad not found")
		return
	}

	RespondSuccess(c, http.StatusOK, membership)
}
