// internal/middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Context keys для передачи данных пользователя между middleware и handler
const (
	CtxKeyUserID   = "user_id"
	CtxKeyUserRole = "user_role"
)

// Claims представляет JWT claims, ожидаемые от будущего auth сервиса
type Claims struct {
	jwt.RegisteredClaims
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
}

// AuthRequired создает middleware, проверяющий валидность JWT токена.
// Извлекает токен из заголовка Authorization: Bearer <token>,
// парсит claims (user_id, role) и помещает их в gin.Context.
func AuthRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authorization header required",
			})
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid or expired token",
			})
			return
		}

		c.Set(CtxKeyUserID, claims.UserID)
		c.Set(CtxKeyUserRole, claims.Role)

		c.Next()
	}
}

// AdminOnly middleware — пропускает только пользователей с ролью "admin".
// Должен использоваться после AuthRequired.
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(CtxKeyUserRole)
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Admin access required",
			})
			return
		}
		c.Next()
	}
}

// GetUserID извлекает user_id из контекста
func GetUserID(c *gin.Context) (int64, bool) {
	val, ok := c.Get(CtxKeyUserID)
	if !ok {
		return 0, false
	}
	id, ok := val.(int64)
	return id, ok
}

// GetUserRole извлекает роль пользователя из контекста
func GetUserRole(c *gin.Context) (string, bool) {
	val, ok := c.Get(CtxKeyUserRole)
	if !ok {
		return "", false
	}
	role, ok := val.(string)
	return role, ok
}

// extractToken извлекает Bearer token из заголовка Authorization
func extractToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return parts[1]
}