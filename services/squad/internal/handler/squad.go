// internal/handler/squad.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

// CreateSquad создает новый отряд
// POST /squads
func (h *Handler) CreateSquad(c *gin.Context) {
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

// GetSquad возвращает отряд по ID
// GET /squads/:id
func (h *Handler) GetSquad(c *gin.Context) {
	id, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid squad ID")
		return
	}

	squad, err := h.squadRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		RespondNotFound(c, "Squad not found")
		return
	}

	RespondSuccess(c, http.StatusOK, squad)
}

// GetSquads возвращает список отрядов с пагинацией
// GET /squads
func (h *Handler) GetSquads(c *gin.Context) {
	limit, offset := GetLimitOffset(c)

	squads, err := h.squadRepo.List(c.Request.Context(), limit, offset)
	if err != nil {
		RespondInternalError(c, "Failed to get squads: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, squads)
}

// GetSquadsByRegion возвращает список отрядов по региону с пагинацией
// GET /regions/:region_id/squads
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

	squads, err := h.squadRepo.GetByRegionID(c.Request.Context(), regionID, limit, offset)
	if err != nil {
		RespondInternalError(c, "Failed to get squads by region: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, squads)
}

// UpdateSquad обновляет существующий отряд
// PUT /squads/:id
func (h *Handler) UpdateSquad(c *gin.Context) {
	id, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid squad ID")
		return
	}

	var req UpdateSquadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, "Invalid request body: "+err.Error())
		return
	}

	// Проверяем, что отряд существует
	squad, err := h.squadRepo.GetByID(c.Request.Context(), id)
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
		if err == nil && existingSquad != nil && existingSquad.ID != id {
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

// DeleteSquad удаляет отряд по ID
// DELETE /squads/:id
func (h *Handler) DeleteSquad(c *gin.Context) {
	id, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid squad ID")
		return
	}

	// Проверяем, что отряд существует
	_, err = h.squadRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		RespondNotFound(c, "Squad not found")
		return
	}

	if err := h.squadRepo.Delete(c.Request.Context(), id); err != nil {
		RespondInternalError(c, "Failed to delete squad: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"message": "Squad deleted successfully",
	})
}