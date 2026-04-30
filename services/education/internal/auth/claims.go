// internal/auth/claims.go
package auth

import "github.com/golang-jwt/jwt/v5"

// Claims представляет custom claims для JWT-токенов платформы.
type Claims struct {
	UserID  string `json:"user_id"`  // Идентификатор пользователя (строка для совместимости с JWT)
	SquadID int64  `json:"squad_id,omitempty"` // ID отряда (опционально, для тестов)
	jwt.RegisteredClaims
}