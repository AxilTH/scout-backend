// internal/handler/handler.go
package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/AxilTH/scout-backend/services/squad/internal/repository"
)

// Handler содержит зависимости для всех handlers
type Handler struct {
	regionRepo         repository.RegionRepository
	roleRepo           repository.RoleRepository
	positionRepo       repository.PositionRepository
	squadRepo          repository.SquadRepository
	squadMembershipRepo repository.SquadMembershipRepository
	squadLeadershipRepo repository.SquadLeadershipRepository
}

// NewHandler создает новый Handler
func NewHandler(
	regionRepo repository.RegionRepository,
	roleRepo repository.RoleRepository,
	positionRepo repository.PositionRepository,
	squadRepo repository.SquadRepository,
	squadMembershipRepo repository.SquadMembershipRepository,
	squadLeadershipRepo repository.SquadLeadershipRepository,
) *Handler {
	return &Handler{
		regionRepo:         regionRepo,
		roleRepo:           roleRepo,
		positionRepo:       positionRepo,
		squadRepo:          squadRepo,
		squadMembershipRepo: squadMembershipRepo,
		squadLeadershipRepo: squadLeadershipRepo,
	}
}
// Helper функции для работы с запросами

// GetIDFromPath извлекает ID из URL-параметра :id
func GetIDFromPath(c *gin.Context) (int64, error) {
	return GetIDFromParam(c, "id")
}

// GetIDFromParam извлекает ID из произвольного URL-параметра
func GetIDFromParam(c *gin.Context, param string) (int64, error) {
	idStr := c.Param(param)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetLimitOffset извлекает параметры пагинации из query
func GetLimitOffset(c *gin.Context) (limit, offset int) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ = strconv.Atoi(limitStr)
	offset, _ = strconv.Atoi(offsetStr)

	// Ограничиваем лимит
	if limit > 100 {
		limit = 100
	}
	if limit <= 0 {
		limit = 10
	}

	return limit, offset
}

// RespondSuccess отправляет успешный ответ
func RespondSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
	})
}

// RespondError отправляет ответ с ошибкой
func RespondError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   message,
	})
}

// RespondValidationError отправляет ответ с ошибкой валидации
func RespondValidationError(c *gin.Context, message string) {
	RespondError(c, http.StatusBadRequest, message)
}

// RespondNotFound отправляет ответ "не найдено"
func RespondNotFound(c *gin.Context, message string) {
	RespondError(c, http.StatusNotFound, message)
}

// RespondInternalError отправляет ответ с внутренней ошибкой
func RespondInternalError(c *gin.Context, message string) {
	RespondError(c, http.StatusInternalServerError, message)
}

// GetCurrentTime возвращает текущее время в UTC
func GetCurrentTime() time.Time {
	return time.Now().UTC()
}