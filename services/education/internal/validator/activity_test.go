package validator

import (
	"testing"
	"time"
)

func intPtr(i int) *int {
	return &i
}

func TestCreateActivityInput_StructuralValidation(t *testing.T) {
	v := GetValidator()

	tests := []struct {
		name        string
		input       CreateActivityInput
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid input",
			input: CreateActivityInput{
				CourseID:       "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				ActivityTypeID: "a1b2c3d4-e5f6-4890-b1c2-d3e4f5a6b7c8",
				Title:          "Valid Activity",
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       intPtr(10),
				IsPublished:    true,
			},
			expectError: false,
		},
		{
			name: "missing course_id",
			input: CreateActivityInput{
				ActivityTypeID: "a1b2c3d4-e5f6-4890-b1c2-d3e4f5a6b7c8",
				Title:          "No CourseID",
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
			},
			expectError: true,
			errorMsg:    "CourseID is required",
		},
		{
			name: "invalid course_id uuid",
			input: CreateActivityInput{
				CourseID:       "not-a-uuid",
				ActivityTypeID: "a1b2c3d4-e5f6-4890-b1c2-d3e4f5a6b7c8",
				Title:          "Invalid UUID",
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
			},
			expectError: true,
			errorMsg:    "Invalid UUID format for CourseID",
		},
		{
			name: "invalid activity_type_id uuid",
			input: CreateActivityInput{
				CourseID:       "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				ActivityTypeID: "not-a-uuid",
				Title:          "Invalid UUID",
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
			},
			expectError: true,
			errorMsg:    "Invalid UUID format for ActivityTypeID",
		},
		{
			name: "empty title",
			input: CreateActivityInput{
				CourseID:       "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				ActivityTypeID: "a1b2c3d4-e5f6-4890-b1c2-d3e4f5a6b7c8",
				Title:          "",
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
			},
			expectError: true,
			errorMsg:    "Title must not be empty",
		},
		{
			name: "negative max_score",
			input: CreateActivityInput{
				CourseID:       "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				ActivityTypeID: "a1b2c3d4-e5f6-4890-b1c2-d3e4f5a6b7c8",
				Title:          "Negative Score",
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       intPtr(-5),
			},
			expectError: true,
			errorMsg:    "MaxScore is too small",
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

func TestCreateActivityInput_BusinessValidation(t *testing.T) {
	tests := []struct {
		name        string
		startsAt    time.Time
		endsAt      time.Time
		expectError bool
	}{
		{
			name:        "ends_at after starts_at",
			startsAt:    time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
			endsAt:      time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "ends_at equals starts_at",
			startsAt:    time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
			endsAt:      time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
			expectError: true,
		},
		{
			name:        "ends_at before starts_at",
			startsAt:    time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
			endsAt:      time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := CreateActivityInput{
				CourseID:       "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				ActivityTypeID: "a1b2c3d4-e5f6-4890-b1c2-d3e4f5a6b7c8",
				Title:          "Test Activity",
				StartsAt:       tt.startsAt,
				EndsAt:         tt.endsAt,
			}
			err := input.Validate()
			if (err != nil) != tt.expectError {
				t.Errorf("Validate() error = %v, expectError = %v", err, tt.expectError)
			}
		})
	}
}