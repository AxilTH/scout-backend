// internal/middleware/logging.go
package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type LogginConfig struct {
	EnableRequestLogging bool
	EnableResponseLogging bool
}

func DefaultLoggingConfig() LogginConfig {
	return LogginConfig{
		EnableRequestLogging: true,
		EnableResponseLogging: true,
	}
}

func RequestLoggerMiddleware(config LogginConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.EnableRequestLogging {
			log.Printf("REQUEST: %s %s", c.Request.Method, c.Request.URL.Path)
		}

		start := time.Now()
		c.Next()

		if config.EnableResponseLogging {
			log.Printf("RESPONSE: %s %s - Status: %d - Duration: %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start))
		}
	}
}