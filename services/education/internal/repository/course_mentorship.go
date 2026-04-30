package repository

// import (
// 	"database/sql"
// 	"time"

// 	"github.com/AxilTH/scout-backend/services/education/internal/model"
// 	"github.com/jmoiron/sqlx"
// )

// type CourseMentorshipRepository struct {
// 	db *sqlx.DB
// }

// func NewCourseMentorshipRepository(db *sqlx.DB) *CourseMentorshipRepository {
// 	return &CourseMentorshipRepository{db: db}
// }

// func (r *CourseMentorshipRepository) Create(mentorship *model.CourseMentorship) error {
// 	mentorship.AssignedAt = time.Now()

// 	query := `
// 		INSERT INTO course_mentorships (course_id, mentor_user_id, assigned_at)
// 		VALUES ($1, $2, $3)
// 		RETURNING id
// 	`

// 	err := r.db.QueryRow(
// 		query,
// 		mentorship.CourseID,
// 		mentorship.MentorUserID,
// 		mentorship.AssignedAt,
// 	).Scan(&mentorship.ID)

// 	return err
// }

// func (r *CourseMentorshipRepository) ListByCourse(courseID int64) ([]model.CourseMentorship, error) {
// 	var mentorships []model.CourseMentorship
// 	err := r.db.Select(&mentorships, "SELECT * FROM course_mentorships WHERE course_id = $1 ORDER BY assigned_at DESC", courseID)
// 	return mentorships, err
// }

// func (r *CourseMentorshipRepository) Delete(courseID, mentorUserID int64) error {
// 	_, err := r.db.Exec("DELETE FROM course_mentorships WHERE course_id = $1 AND mentor_user_id = $2", courseID, mentorUserID)
// 	return err
// }

// func (r *CourseMentorshipRepository) Exists(courseID, mentorUserID int64) (bool, error) {
// 	var exists bool
// 	err := r.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM course_mentorships WHERE course_id = $1 AND mentor_user_id = $2)", courseID, mentorUserID)
// 	if err == sql.ErrNoRows {
// 		return false, nil
// 	}
// 	return exists, err
// }