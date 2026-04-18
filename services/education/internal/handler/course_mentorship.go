package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/AxilTH/scout-backend/services/education/internal/model"
	"github.com/AxilTH/scout-backend/services/education/internal/repository"
	"github.com/AxilTH/scout-backend/services/education/internal/validator"
	"github.com/gin-gonic/gin"
)

func AssignMentor(c *gin.Context, repo *repository.CourseMentorshipRepository) {
	courseIDStr := c.Param("id")
	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
		return
	}

	var req validator.AssignMentorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	input := validator.AssignMentorInput{
		MentorUserID: req.MentorUserID,
	}

	if err := validate.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.FormatError(err)})
		return
	}

	// Проверяем, не назначен ли уже этот куратор
	exists, err := repo.Exists(courseID, req.MentorUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mentor already assigned to this course"})
		return
	}

	mentorship := &model.CourseMentorship{
		CourseID: courseID,
		MentorUserID: req.MentorUserID,
	}

	if err := repo.Create(mentorship); err != nil {
		log.Printf("Failed to assign mentor: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign mentor"})
		return
	}

	c.JSON(http.StatusCreated, mentorship)
}

func GetCourseMentors(c *gin.Context, repo *repository.CourseMentorshipRepository) {
	courseIDStr := c.Param("id")
	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
		return
	}

	mentorships, err := repo.ListByCourse(courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return 
	}

	c.JSON(http.StatusOK, mentorships)
}

func RemoveMentor(c *gin.Context, repo *repository.CourseMentorshipRepository) {
	courseIDStr := c.Param("id")
	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
		return
	}

	mentorUserIDStr := c.Param("user_id")
	mentorUserID, err := strconv.ParseInt(mentorUserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mentor user ID format"})
		return
	}

	if err := repo.Delete(courseID, mentorUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove mentor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mentor removed"})
}
