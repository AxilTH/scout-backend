// internal/middleware/leadership.go
package middleware

import (
	"net/http"
	"strconv"

	"github.com/AxilTH/scout-backend/services/squad/internal/repository"
	"github.com/gin-gonic/gin"
)

// CommanderOnly middleware — пропускает только пользователей с должностью "commander" в указанном отряде.
// Должен использоваться после AuthRequired.
// squad_id извлекается из параметра пути :id.
func CommanderOnly(
	squadLeadershipRepo repository.SquadLeadershipRepository,
	positionRepo repository.PositionRepository,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем user_id из контекста
		userID, ok := GetUserID(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "User ID not found in context",
			})
			return
		}

		// Извлекаем squad_id из пути
		squadIDStr := c.Param("id")
		squadID, err := strconv.ParseInt(squadIDStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid squad ID",
			})
			return
		}

		// Получаем должность "commander"
		commanderPosition, err := positionRepo.GetByTitle(c.Request.Context(), "commander")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to get commander position",
			})
			return
		}

		// Получаем все должности пользователя в этом отряде
		leaderships, err := squadLeadershipRepo.GetBySquadID(c.Request.Context(), squadID, 100, 0)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to get leadership",
			})
			return
		}

		// Проверяем, имеет ли пользователь должность командира
		isCommander := false
		for _, leadership := range leaderships {
			if leadership.UserID == userID &&
				leadership.PositionID == commanderPosition.ID &&
				leadership.DismissedAt == nil {
				isCommander = true
				break
			}
		}

		if !isCommander {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Commander access required",
			})
			return
		}

		c.Next()
	}
}