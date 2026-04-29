package validator

// import (
// 	"fmt"
// 	"time"
// )

// type CreateAssignmentRequest struct {
// 	ActivityID       *int64    `json:"activity_id"`
// 	AssignmentTypeID int64     `json:"assignment_type_id"`
// 	Title            string    `json:"title"`
// 	Description      *string   `json:"description"`
// 	Deadline         time.Time `json:"deadline"`
// 	MaxScore         int       `json:"max_score"`
// 	IsPublished      bool      `json:"is_published"`
// }

// type CreateAssignmentInput struct {
// 	ActivityID       *int64    `validate:"omitempty,min=1"`
// 	AssignmentTypeID int64     `validate:"required,min=1"`
// 	Title            string    `validate:"required,min=1,max=255"`
// 	Description      *string   `validate:"omitempty"`
// 	Deadline         time.Time `validate:"required"`
// 	MaxScore         int       `validate:"required,min=0"`
// 	IsPublished      bool      `validate:"-"`
// }

// func (c *CreateAssignmentInput) Validate() error {
// 	if !c.Deadline.After(time.Now()) {
// 		return fmt.Errorf("deadline must be in the future")
// 	}

// 	return nil
// }
