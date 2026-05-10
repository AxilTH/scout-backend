// internal/handler/squad.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/AxilTH/scout-backend/services/squad/internal/middleware"
	"github.com/AxilTH/scout-backend/services/squad/internal/model"
)

// CreateSquadRequest представляет запрос на создание отряда
type CreateSquadRequest struct {
	Title    string `json:"title" binding:"required"`
	RegionID int64  `json:"region_id" binding:"required"`
}

// UpdateSquadRequest представляет запрос на обновление отряда
type UpdateSquadRequest struct {
	Title    string `json:"title" binding:"required"`
	RegionID int64  `json:"region_id" binding:"required"`
}

// CreateSquad создает новый отряд.
// Только admin может создавать отряды.
// POST /squads
func (h *Handler) CreateSquad(c *gin.Context) {
	// Проверяем, что пользователь admin
	role, _ := middleware.GetUserRole(c)
	if role != "admin" {
		RespondError(c, http.StatusForbidden, "Admin access required")
		return
	}

	var req CreateSquadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, "Invalid request body: "+err.Error())
		return
	}

	// Проверяем, что регион существует
	_, err := h.regionRepo.GetByID(c.Request.Context(), req.RegionID)
	if err != nil {
		RespondValidationError(c, "Region not found")
		return
	}

	// Проверяем, что отряд с таким названием еще не существует
	existingSquad, err := h.squadRepo.GetByTitle(c.Request.Context(), req.Title)
	if err == nil && existingSquad != nil {
		RespondError(c, http.StatusConflict, "Squad with this title already exists")
		return
	}

	// Создаем новый отряд
	squad := &model.Squad{
		Title:     req.Title,
		RegionID:  req.RegionID,
		CreatedAt: GetCurrentTime(),
		UpdatedAt: GetCurrentTime(),
	}

	if err := h.squadRepo.Create(c.Request.Context(), squad); err != nil {
		RespondInternalError(c, "Failed to create squad: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusCreated, squad)
}

// GetCurrentSquad возвращает текущий отряд пользователя.
// Fighter видит только свой текущий отряд, admin — все отряды.
// GET /squads/current
func (h *Handler) GetCurrentSquad(c *gin.Context) {
	// Извлекаем squad_id из контекста (для изоляции данных)
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		RespondInternalError(c, "Failed to get squad information: "+err.Error())
		return
	}

	// Проверяем, что отряд существует
	squad, err := h.squadRepo.GetByID(c.Request.Context(), squadID)
	if err != nil {
		RespondNotFound(c, "Squad not found")
		return
	}

	RespondSuccess(c, http.StatusOK, squad)
}

// GetSquads возвращает список отрядов с пагинацией.
// Админ получает все отряды, fighter — только свой текущий отряд.
// GET /squads
func (h *Handler) GetSquads(c *gin.Context) {
	limit, offset := GetLimitOffset(c)

	// Fighter видит только свой текущий отряд — squad_id из контекста
	role, _ := middleware.GetUserRole(c)
	if role == "fighter" {
		squadID, err := getSquadIDFromContext(c)
		if err != nil {
			RespondInternalError(c, "Failed to get squad information: "+err.Error())
			return
		}

		squad, err := h.squadRepo.GetByID(c.Request.Context(), squadID)
		if err != nil {
			RespondNotFound(c, "Squad not found")
			return
		}

		RespondSuccess(c, http.StatusOK, squad)
		return
	}

	// Admin получает все отряды
	squads, err := h.squadRepo.List(c.Request.Context(), limit, offset)
	if err != nil {
		RespondInternalError(c, "Failed to get squads: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, squads)
}

// GetSquadsByRegion возвращает список отрядов по региону с пагинацией.
// Админ получает все отряды региона, fighter — только свои отряды в этом регионе.
// GET /regions/:id/squads
func (h *Handler) GetSquadsByRegion(c *gin.Context) {
	regionID, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid region ID")
		return
	}

	// Проверяем, что регион существует
	_, err = h.regionRepo.GetByID(c.Request.Context(), regionID)
	if err != nil {
		RespondNotFound(c, "Region not found")
		return
	}

		limit, offset := GetLimitOffset(c)

	// Fighter видит только свои отряды в регионе — squad_id из контекста
	role, _ := middleware.GetUserRole(c)
	if role == "fighter" {
		squadID, err := getSquadIDFromContext(c)
		if err != nil {
			RespondInternalError(c, "Failed to get squad information: "+err.Error())
			return
		}

		squad, err := h.squadRepo.GetByID(c.Request.Context(), squadID)
		if err != nil {
			RespondNotFound(c, "Squad not found")
			return
		}

	RespondSuccess(c, http.StatusOK, squad)
		return
	}

	// Admin получает все отряды региона
	squads, err := h.squadRepo.GetByRegionID(c.Request.Context(), regionID, limit, offset)
	if err != nil {
		RespondInternalError(c, "Failed to get squads by region: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, squads)
}

// UpdateCurrentSquad обновляет текущий отряд пользователя.
// Только commander может обновлять свой отряд.
// PUT /squads/current
func (h *Handler) UpdateCurrentSquad(c *gin.Context) {
	// Извлекаем squad_id из контекста (для изоляции данных)
	squadID, err := getSquadIDFromContext(c)
	if err != nil {
		RespondInternalError(c, "Failed to get squad information: "+err.Error())
		return
	}

	var req UpdateSquadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, "Invalid request body: "+err.Error())
		return
	}

	// Проверяем, что отряд существует
	squad, err := h.squadRepo.GetByID(c.Request.Context(), squadID)
	if err != nil {
		RespondNotFound(c, "Squad not found")
		return
	}

	// Проверяем, что регион существует
	_, err = h.regionRepo.GetByID(c.Request.Context(), req.RegionID)
	if err != nil {
		RespondValidationError(c, "Region not found")
		return
	}

	// Проверяем, что отряд с таким названием еще не существует (если название изменилось)
	if req.Title != squad.Title {
		existingSquad, err := h.squadRepo.GetByTitle(c.Request.Context(), req.Title)
		if err == nil && existingSquad != nil && existingSquad.ID != squadID {
			RespondError(c, http.StatusConflict, "Squad with this title already exists")
			return
		}
	}

	// Обновляем отряд
	squad.Title = req.Title
	squad.RegionID = req.RegionID
	squad.UpdatedAt = GetCurrentTime()

	if err := h.squadRepo.Update(c.Request.Context(), squad); err != nil {
		RespondInternalError(c, "Failed to update squad: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, squad)
}

// DeleteCurrentSquad удаляет текущий отряд пользователя.
// Только admin может удалять отряды.
// DELETE /squads/current
func (h *Handler) DeleteCurrentSquad(c *gin.Context) {
	// Проверяем, что пользователь admin
	role, _ := middleware.GetUserRole(c)
	if role != "admin" {
		RespondError(c, http.StatusForbidden, "Admin access required")
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

	if err := h.squadRepo.Delete(c.Request.Context(), squadID); err != nil {
		RespondInternalError(c, "Failed to delete squad: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"message": "Squad deleted successfully",
	})
}