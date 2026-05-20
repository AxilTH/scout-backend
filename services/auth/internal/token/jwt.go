// internal/token/jwt.go
package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims представляет собой набор утверждений (claims) для JWT токена.
type Claims struct {
	UserID       int64    `json:"user_id"`
	SquadIDs     []int64  `json:"squad_ids"`
	CurrentSquad int64    `json:"current_squad_id"`
	jwt.RegisteredClaims                    // стандартные поля: iss, iat, exp, sub, aud и т.д.
}

// GenerateJWT создаёт signed JWT токен с указанными claims и секретом.
// ttl — время жизни токена.
func GenerateJWT(userID int64, squadIDs []int64, currentSquadID int64, secret string, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID:       userID,
		SquadIDs:     squadIDs,
		CurrentSquad: currentSquadID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "scout-backend",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseJWT проверяет токен, возвращает claims если токен валиден.
func ParseJWT(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// Проверяем, что метод подписи соответствует ожидаемому
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

// RefreshJWT создает новый токен на основе старого, увеличивая время жизни на extraTTL.
// Если старый токен просрочен, возвращает ошибку.
func RefreshJWT(tokenString, secret string, extraTTL time.Duration) (string, error) {
	claims, err := ParseJWT(tokenString, secret)
	if err != nil {
		return "", err
	}
	// Сбрасываем время выдачи и истечения
	now := time.Now().UTC()
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(extraTTL))

	// Генерируем новый токен с теми же claims
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return newToken.SignedString([]byte(secret))
}