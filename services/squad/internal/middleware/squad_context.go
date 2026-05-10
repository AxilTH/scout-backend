// internal/middleware/squad_context.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SquadContextRequired middleware — проверяет, что squad_id присутствует в контексте.
// Должен использоваться после AuthRequired.
// squad_id извлекается из JWT токена и помещается в контекст AuthRequired middleware.
func SquadContextRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := GetSquadID(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Squad ID not found in context",
			})
			return
		}

		// squad_id уже установлен в контексте AuthRequired middleware
		// Этот middleware просто проверяет его наличие
		c.Next()
	}
}