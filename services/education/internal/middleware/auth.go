// internal/middleware/auth.go
package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/AxilTH/scout-backend/services/education/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware возвращает Gin-middleware для проверки JWT-токенов.
// Секрет должен совпадать с тем, что используется при генерации токенов.
func AuthMiddleware(secret string) gin.HandlerFunc {
	// Валидация секрета при инициализации (fail-fast)
	if secret == "" {
		panic("AuthMiddleware: JWT_SECRET cannot be empty")
	}

	return func(c *gin.Context) {
		// 1. Извлекаем заголовок Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		// 2. Парсим формат "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}
		tokenString := parts[1]

		// 3. Парсим и верифицируем токен
		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			// Проверяем алгоритм подписи (защита от alg=none и подмены)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})

		// 4. Обрабатываем ошибки валидации
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// 5. Извлекаем claims и сохраняем в контекст
		c.Set("user_id", claims.UserID)

		// Squad ID обязателен - если его нет, токен невалиден
		if claims.SquadID <= 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Squad ID must be present in token"})
			return
		}

		c.Set("squad_id", claims.SquadID)
		c.Next()
	}
}