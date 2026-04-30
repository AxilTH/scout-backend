package repository

// import (
// 	"database/sql"
// 	"time"

// 	"github.com/AxilTH/scout-backend/services/education/internal/model"
// 	"github.com/jmoiron/sqlx"
// )

// type AssignmentRepository struct {
// 	db *sqlx.DB
// }

// func NewAssignmentRepository(db *sqlx.DB) *AssignmentRepository {
// 	return &AssignmentRepository{db: db}
// }

// func (r *AssignmentRepository) Create(assignment *model.Assignment) error {
// 	assignment.CreatedAt = time.Now()
// 	assignment.UpdatedAt = time.Now()

// 	query := `
// 		INSERT INTO assignments
// 		(activity_id, assignment_type_id, title, description, deadline, max_score, is_published, created_at, updated_at)
// 		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
// 		RETURNING id
// 	`

// 	err := r.db.QueryRow(
// 		query,
// 		assignment.ActivityID,
// 		assignment.AssignmentTypeID,
// 		assignment.Title,
// 		assignment.Description,
// 		assignment.Deadline,
// 		assignment.MaxScore,
// 		assignment.IsPublished,
// 		assignment.CreatedAt,
// 		assignment.UpdatedAt,
// 	).Scan(&assignment.ID)

// 	return err
// }

// func (r *AssignmentRepository) GetByID(id int64) (*model.Assignment, error) {
// 	var assignment model.Assignment
// 	err := r.db.Get(&assignment, "SELECT * FROM assignments WHERE id = $1", id)
// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	}

// 	return &assignment, err
// }

// func (r *AssignmentRepository) ListByActivityID(activityId int64) ([]model.Assignment, error) {
// 	var assignments []model.Assignment
// 	err := r.db.Select(&assignments, "SELECT * FROM assignments WHERE activity_id = $1 ORDER BY deadline ASC", activityId)
// 	return assignments, err
// }

// func (r *AssignmentRepository) ListAll() ([]model.Assignment, error) {
// 	var assignments []model.Assignment
// 	err := r.db.Select(&assignments, "SELECT * FROM assignments ORDER BY deadline ASC")
// 	return assignments, err
// }

// func (r *AssignmentRepository) Update(assignment *model.Assignment) error {
// 	assignment.UpdatedAt = time.Now()

// 	query := `
// 		UPDATE assignments
// 		SET title = $1, description = $2, deadline = $3, max_score = $4,
// 			assignment_type_id = $5, is_published = $6, updated_at = $7
// 		WHERE id = $8
// 	`

// 	_, err := r.db.Exec(
// 		query,
// 		assignment.Title,
// 		assignment.Description,
// 		assignment.Deadline,
// 		assignment.MaxScore,
// 		assignment.AssignmentTypeID,
// 		assignment.IsPublished,
// 		assignment.UpdatedAt,
// 		assignment.ID,
// 	)

// 	return err
// }

// func (r *AssignmentRepository) Delete(id int64) error {
// 	_, err := r.db.Exec("DELETE FROM assignments WHERE id = $1", id)
// 	return err
// }
