// internal/handler/region.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/AxilTH/scout-backend/services/squad/internal/model"
)

// CreateRegionRequest представляет запрос на создание региона
type CreateRegionRequest struct {
	Title string `json:"title" binding:"required"`
}

// UpdateRegionRequest представляет запрос на обновление региона
type UpdateRegionRequest struct {
	Title string `json:"title" binding:"required"`
}

// CreateRegion создает новый регион
// POST /regions
func (h *Handler) CreateRegion(c *gin.Context) {
	var req CreateRegionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, "Invalid request body: "+err.Error())
		return
	}

	// Проверяем, что регион с таким названием еще не существует
	existingRegion, err := h.regionRepo.GetByTitle(c.Request.Context(), req.Title)
	if err == nil && existingRegion != nil {
		RespondError(c, http.StatusConflict, "Region with this title already exists")
		return
	}

	// Создаем новый регион
	region := &model.Region{
		Title:     req.Title,
		CreatedAt: GetCurrentTime(),
		UpdatedAt: GetCurrentTime(),
	}

	if err := h.regionRepo.Create(c.Request.Context(), region); err != nil {
		RespondInternalError(c, "Failed to create region: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusCreated, region)
}

// GetRegion возвращает регион по ID
// GET /regions/:id
func (h *Handler) GetRegion(c *gin.Context) {
	id, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid region ID")
		return
	}

	region, err := h.regionRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		RespondNotFound(c, "Region not found")
		return
	}

	RespondSuccess(c, http.StatusOK, region)
}

// GetRegions возвращает список регионов с пагинацией
// GET /regions
func (h *Handler) GetRegions(c *gin.Context) {
	limit, offset := GetLimitOffset(c)

	regions, err := h.regionRepo.List(c.Request.Context(), limit, offset)
	if err != nil {
		RespondInternalError(c, "Failed to get regions: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, regions)
}

// UpdateRegion обновляет существующий регион
// PUT /regions/:id
func (h *Handler) UpdateRegion(c *gin.Context) {
	id, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid region ID")
		return
	}

	var req UpdateRegionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, "Invalid request body: "+err.Error())
		return
	}

	// Проверяем, что регион существует
	region, err := h.regionRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		RespondNotFound(c, "Region not found")
		return
	}

	// Проверяем, что регион с таким названием еще не существует (если название изменилось)
	if req.Title != region.Title {
		existingRegion, err := h.regionRepo.GetByTitle(c.Request.Context(), req.Title)
		if err == nil && existingRegion != nil && existingRegion.ID != id {
			RespondError(c, http.StatusConflict, "Region with this title already exists")
			return
		}
	}

	// Обновляем регион
	region.Title = req.Title
	region.UpdatedAt = GetCurrentTime()

	if err := h.regionRepo.Update(c.Request.Context(), region); err != nil {
		RespondInternalError(c, "Failed to update region: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, region)
}

// DeleteRegion удаляет регион по ID
// DELETE /regions/:id
func (h *Handler) DeleteRegion(c *gin.Context) {
	id, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid region ID")
		return
	}

	// Проверяем, что регион существует
	_, err = h.regionRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		RespondNotFound(c, "Region not found")
		return
	}

	if err := h.regionRepo.Delete(c.Request.Context(), id); err != nil {
		RespondInternalError(c, "Failed to delete region: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"message": "Region deleted successfully",
	})
}