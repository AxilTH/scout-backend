// internal/handler/health.go
package handler

import "github.com/gin-gonic/gin"

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}