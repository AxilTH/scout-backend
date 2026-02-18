package handler

import "github.com/gin-gonic/gin"

func GetCourses(c *gin.Context) {
	c.JSON(200, []gin.H{
		{"id": "1", "title": "Введение в педагогику", "year": 2026},
	})
}

func GetCourse(c *gin.Context) {
	courseID := c.Param("id")
	c.JSON(200, gin.H{
		"id":   courseID,
		"title": "Курс #" + courseID,
		"year": 2026,
	})
}

func CreateCourse(c *gin.Context) {
	c.JSON(201, gin.H{
		"message": "Курс создан",
		"id":      "new-course-id",
	})
}