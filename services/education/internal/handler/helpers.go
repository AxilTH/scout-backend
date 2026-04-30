// internal/handler/helpers.go
package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// getSquadIDFromContext извлекает squad_id из контекста Gin
// Возвращает ошибку, если squad_id отсутствует или имеет неверный тип
func getSquadIDFromContext(c *gin.Context) (int64, error) {
	squadID, ok := c.Get("squad_id")
	if !ok {
		return 0, errors.New("squad_id not found in context")
	}

	squadIDInt, ok := squadID.(int64)
	if !ok {
		return 0, errors.New("squad_id has invalid type")
	}

	return squadIDInt, nil
}

// getUserIDFromContext извлекает user_id из контекста Gin
// Возвращает ошибку, если user_id отсутствует или имеет неверный тип
func getUserIDFromContext(c *gin.Context) (string, error) {
	userID, ok := c.Get("user_id")
	if !ok {
		return "", errors.New("user_id not found in context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", errors.New("user_id has invalid type")
	}

	return userIDStr, nil
}

// parseIDParam парсит параметр ID из URL пути
// paramName - имя параметра (например, "id", "course_id", "activity_id")
func parseIDParam(c *gin.Context, paramName string) (int64, error) {
	idStr := c.Param(paramName)
	if idStr == "" {
		return 0, fmt.Errorf("%s parameter is required", paramName)
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format", paramName)
	}

	return id, nil
}

// respondWithError отправляет ошибку в формате JSON
func respondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

// respondWithValidationError отправляет ошибки валидации в формате JSON
func respondWithValidationError(c *gin.Context, errors []string) {
	c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
}

// respondWithSuccess отправляет успешный ответ в формате JSON
func respondWithSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

// respondWithCreated отправляет ответ 201 Created
func respondWithCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// respondWithNoContent отправляет ответ 204 No Content
func respondWithNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// respondWithNotFound отправляет ответ 404 Not Found
func respondWithNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
}

// respondWithForbidden отправляет ответ 403 Forbidden
func respondWithForbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Access denied"
	}
	c.JSON(http.StatusForbidden, gin.H{"error": message})
}

// respondWithInternalError отправляет ответ 500 Internal Server Error
func respondWithInternalError(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": message})
}