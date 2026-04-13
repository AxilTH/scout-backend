package validator

type SubmitActivityResultRequest struct {
	UserID     int64 `json:"user_id"`
	ActivityID int64 `json:"activity_id"`
}

type SubmitActivityResultInput struct {
	UserID     int64 `validate:"required,min=1"`
	ActivityID int64 `validate:"required,min=1"`
}

type ReviewActivityResultRequest struct {
	Score    *int    `json:"score"`
	Feedback *string `json:"feedback"`
	Status   string  `json:"status"`
}

type ReviewActivityResultInput struct {
	Score    *int    `validate:"omitempty,min=0"`
	Feedback *string `validate:"omitempty"`
	Status   string  `validate:"required,oneof=reviewed rejected"`
}
