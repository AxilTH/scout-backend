package validator

import (
	"testing"
)

func TestAssignMentorInput_Validation(t *testing.T) {
	v := GetValidator()

	tests := []struct {
		name        string
		input       AssignMentorInput
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid input",
			input: AssignMentorInput{
				MentorUserID: 1,
			},
			expectError: false,
		},
		{
			name: "missing mentor_user_id",
			input: AssignMentorInput{
				MentorUserID: 0,
			},
			expectError: true,
			errorMsg: "MentorUserID is required",
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
