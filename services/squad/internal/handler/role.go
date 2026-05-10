// internal/handler/role.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/AxilTH/scout-backend/services/squad/internal/model"
)

// CreateRoleRequest представляет запрос на создание роли
type CreateRoleRequest struct {
	Title string `json:"title" binding:"required"`
}

// UpdateRoleRequest представляет запрос на обновление роли
type UpdateRoleRequest struct {
	Title string `json:"title" binding:"required"`
}

// CreateRole создает новую роль
// POST /roles
func (h *Handler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, "Invalid request body: "+err.Error())
		return
	}

	// Проверяем, что роль с таким названием еще не существует
	existingRole, err := h.roleRepo.GetByTitle(c.Request.Context(), req.Title)
	if err == nil && existingRole != nil {
		RespondError(c, http.StatusConflict, "Role with this title already exists")
		return
	}

	// Создаем новую роль
	role := &model.Role{
		Title:     req.Title,
		CreatedAt: GetCurrentTime(),
		UpdatedAt: GetCurrentTime(),
	}

	if err := h.roleRepo.Create(c.Request.Context(), role); err != nil {
		RespondInternalError(c, "Failed to create role: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusCreated, role)
}

// GetRole возвращает роль по ID
// GET /roles/:id
func (h *Handler) GetRole(c *gin.Context) {
	id, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid role ID")
		return
	}

	role, err := h.roleRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		RespondNotFound(c, "Role not found")
		return
	}

	RespondSuccess(c, http.StatusOK, role)
}

// GetRoles возвращает список ролей с пагинацией
// GET /roles
func (h *Handler) GetRoles(c *gin.Context) {
	limit, offset := GetLimitOffset(c)

	roles, err := h.roleRepo.List(c.Request.Context(), limit, offset)
	if err != nil {
		RespondInternalError(c, "Failed to get roles: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, roles)
}

// UpdateRole обновляет существующую роль
// PUT /roles/:id
func (h *Handler) UpdateRole(c *gin.Context) {
	id, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid role ID")
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, "Invalid request body: "+err.Error())
		return
	}

	// Проверяем, что роль существует
	role, err := h.roleRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		RespondNotFound(c, "Role not found")
		return
	}

	// Проверяем, что роль с таким названием еще не существует (если название изменилось)
	if req.Title != role.Title {
		existingRole, err := h.roleRepo.GetByTitle(c.Request.Context(), req.Title)
		if err == nil && existingRole != nil && existingRole.ID != id {
			RespondError(c, http.StatusConflict, "Role with this title already exists")
			return
		}
	}

	// Обновляем роль
	role.Title = req.Title
	role.UpdatedAt = GetCurrentTime()

	if err := h.roleRepo.Update(c.Request.Context(), role); err != nil {
		RespondInternalError(c, "Failed to update role: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, role)
}

// DeleteRole удаляет роль по ID
// DELETE /roles/:id
func (h *Handler) DeleteRole(c *gin.Context) {
	id, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid role ID")
		return
	}

	// Проверяем, что роль существует
	_, err = h.roleRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		RespondNotFound(c, "Role not found")
		return
	}

	if err := h.roleRepo.Delete(c.Request.Context(), id); err != nil {
		RespondInternalError(c, "Failed to delete role: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"message": "Role deleted successfully",
	})
}