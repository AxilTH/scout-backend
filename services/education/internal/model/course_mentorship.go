package model

import "time"

type CourseMentorship struct {
	ID           int64     `db:"id" json:"id"`
	CourseID     int64     `db:"course_id" json:"course_id"`
	MentorUserID int64     `db:"mentor_user_id" json:"mentor_user_id"`
	AssignedAt   time.Time `db:"assigned_at" json:"assigned_at"`
}