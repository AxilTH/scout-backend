// internal/handler/position.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/AxilTH/scout-backend/services/squad/internal/model"
)

// CreatePositionRequest представляет запрос на создание должности
type CreatePositionRequest struct {
	Title string `json:"title" binding:"required"`
}

// UpdatePositionRequest представляет запрос на обновление должности
type UpdatePositionRequest struct {
	Title string `json:"title" binding:"required"`
}

// CreatePosition создает новую должность
// POST /positions
func (h *Handler) CreatePosition(c *gin.Context) {
	var req CreatePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, "Invalid request body: "+err.Error())
		return
	}

	// Проверяем, что должность с таким названием еще не существует
	existingPosition, err := h.positionRepo.GetByTitle(c.Request.Context(), req.Title)
	if err == nil && existingPosition != nil {
		RespondError(c, http.StatusConflict, "Position with this title already exists")
		return
	}

	// Создаем новую должность
	position := &model.Position{
		Title:     req.Title,
		CreatedAt: GetCurrentTime(),
		UpdatedAt: GetCurrentTime(),
	}

	if err := h.positionRepo.Create(c.Request.Context(), position); err != nil {
		RespondInternalError(c, "Failed to create position: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusCreated, position)
}

// GetPosition возвращает должность по ID
// GET /positions/:id
func (h *Handler) GetPosition(c *gin.Context) {
	id, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid position ID")
		return
	}

	position, err := h.positionRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		RespondNotFound(c, "Position not found")
		return
	}

	RespondSuccess(c, http.StatusOK, position)
}

// GetPositions возвращает список должностей с пагинацией
// GET /positions
func (h *Handler) GetPositions(c *gin.Context) {
	limit, offset := GetLimitOffset(c)

	positions, err := h.positionRepo.List(c.Request.Context(), limit, offset)
	if err != nil {
		RespondInternalError(c, "Failed to get positions: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, positions)
}

// UpdatePosition обновляет существующую должность
// PUT /positions/:id
func (h *Handler) UpdatePosition(c *gin.Context) {
	id, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid position ID")
		return
	}

	var req UpdatePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, "Invalid request body: "+err.Error())
		return
	}

	// Проверяем, что должность существует
	position, err := h.positionRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		RespondNotFound(c, "Position not found")
		return
	}

	// Проверяем, что должность с таким названием еще не существует (если название изменилось)
	if req.Title != position.Title {
		existingPosition, err := h.positionRepo.GetByTitle(c.Request.Context(), req.Title)
		if err == nil && existingPosition != nil && existingPosition.ID != id {
			RespondError(c, http.StatusConflict, "Position with this title already exists")
			return
		}
	}

	// Обновляем должность
	position.Title = req.Title
	position.UpdatedAt = GetCurrentTime()

	if err := h.positionRepo.Update(c.Request.Context(), position); err != nil {
		RespondInternalError(c, "Failed to update position: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, position)
}

// DeletePosition удаляет должность по ID
// DELETE /positions/:id
func (h *Handler) DeletePosition(c *gin.Context) {
	id, err := GetIDFromPath(c)
	if err != nil {
		RespondValidationError(c, "Invalid position ID")
		return
	}

	// Проверяем, что должность существует
	_, err = h.positionRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		RespondNotFound(c, "Position not found")
		return
	}

	if err := h.positionRepo.Delete(c.Request.Context(), id); err != nil {
		RespondInternalError(c, "Failed to delete position: "+err.Error())
		return
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"message": "Position deleted successfully",
	})
}