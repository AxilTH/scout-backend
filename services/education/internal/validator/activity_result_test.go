package validator

import (
	"testing"
)

func TestSubmitActivityResultInput_Validation(t *testing.T) {
	v := GetValidator()

	tests := []struct {
		name        string
		input       SubmitActivityResultInput
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid input",
			input: SubmitActivityResultInput{
				UserID:     1,
				ActivityID: 1,
			},
			expectError: false,
		},
		{
			name: "missing user_id",
			input: SubmitActivityResultInput{
				ActivityID: 1,
			},
			expectError: true,
			errorMsg:    "UserID is required",
		},
		{
			name: "missing activity_id",
			input: SubmitActivityResultInput{
				UserID: 1,
			},
			expectError: true,
			errorMsg:    "ActivityID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(&tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected validation error, got none")
				} else {
					formatted := FormatError(err)
					if formatted != tt.errorMsg {
						t.Errorf("expected error %q, got %q", tt.errorMsg, formatted)
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected validation error: %v", err)
				}
			}
		})
	}
}

func TestReviewActivityResultInput_Validation(t *testing.T) {
	v := GetValidator()

	tests := []struct {
		name        string
		input       ReviewActivityResultInput
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid review - reviewed",
			input: ReviewActivityResultInput{
				Score:    intPtr(18),
				Feedback: strPtr("Good work!"),
				Status:   "reviewed",
			},
			expectError: false,
		},
		{
			name: "valid review - rejected",
			input: ReviewActivityResultInput{
				Feedback: strPtr("Needs improvement"),
				Status:   "rejected",
			},
			expectError: false,
		},
		{
			name: "missing status",
			input: ReviewActivityResultInput{
				Score:    intPtr(18),
				Feedback: strPtr("Good work!"),
			},
			expectError: true,
			errorMsg:    "Status is required",
		},
		{
			name: "invalid status",
			input: ReviewActivityResultInput{
				Score:  intPtr(18),
				Status: "invalid",
			},
			expectError: true,
			errorMsg:    "Invalid value for Status",
		},
		{
			name: "negative score",
			input: ReviewActivityResultInput{
				Score:  intPtr(-5),
				Status: "reviewed",
			},
			expectError: true,
			errorMsg:    "Score is too small",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(&tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected validation error, got none")
				} else {
					formatted := FormatError(err)
					if formatted != tt.errorMsg {
						t.Errorf("expected error %q, got %q", tt.errorMsg, formatted)
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected validation error: %v", err)
				}
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
