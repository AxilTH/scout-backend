// internal/middleware/logging.go
package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingConfig содержит конфигурацию для middleware логирования
type LoggingConfig struct {
	// SkipPaths - пути, которые не нужно логировать
	SkipPaths []string
	// EnableDebug - включить debug-логирование
	EnableDebug bool
}

// DefaultLoggingConfig возвращает конфигурацию по умолчанию
func DefaultLoggingConfig() *LoggingConfig {
	return &LoggingConfig{
		SkipPaths: []string{
			"/health",
			"/metrics",
		},
		EnableDebug: false,
	}
}

// RequestLoggerMiddleware создает middleware для структурированного логирования запросов
func RequestLoggerMiddleware(config *LoggingConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultLoggingConfig()
	}

	// Создаем map для быстрой проверки путей
	skipPathsMap := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipPathsMap[path] = true
	}

	return func(c *gin.Context) {
		// Пропускаем логирование для указанных путей
		if skipPathsMap[c.Request.URL.Path] {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Извлекаем информацию из контекста (если есть)
		userID, _ := c.Get("user_id")
		squadID, _ := c.Get("squad_id")

		// Обрабатываем запрос
		c.Next()

		// Формируем лог после обработки запроса
		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()

		// Структурированный лог в формате, удобном для парсинга в Loki/Prometheus
		log.Printf("[HTTP] method=%s path=%s status=%d latency_ms=%d client_ip=%s user_id=%v squad_id=%v",
			method,
			path,
			status,
			latency.Milliseconds(),
			clientIP,
			userID,
			squadID,
		)

		// Debug-логирование для отладки (только если включено)
		if config.EnableDebug {
			log.Printf("[DEBUG] Headers: %v", c.Request.Header)
		}
	}
}