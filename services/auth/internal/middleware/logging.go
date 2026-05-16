// internal/middleware/logging.go
package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type LoggingConfig struct {
	// TODO: добавить конфигурацию для логирования
}

func DefaultLoggingConfig() LoggingConfig {
	return LoggingConfig{}
}

func RequestLoggerMiddleware(config LoggingConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		log.Printf("[%s] %s %d %v", method, path, statusCode, latency)
	}
}