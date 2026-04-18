package validator

type AssignMentorRequest struct {
	MentorUserID int64 `json:"mentor_user_id"`
}

type AssignMentorInput struct {
	MentorUserID int64 `validate:"required,min=1"`
}