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
				CourseID:       1,
				ActivityTypeID: 1,
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
				ActivityTypeID: 1,
				Title:          "No CourseID",
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
			},
			expectError: true,
			errorMsg:    "CourseID is required",
		},
		{
			name: "invalid course_id zero",
			input: CreateActivityInput{
				CourseID:       0,
				ActivityTypeID: 1,
				Title:          "Invalid CourseID",
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
			},
			expectError: true,
			errorMsg:    "CourseID is required",
		},
		{
			name: "invalid activity_type_id zero",
			input: CreateActivityInput{
				CourseID:       1,
				ActivityTypeID: 0,
				Title:          "Invalid ActivityTypeID",
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
			},
			expectError: true,
			errorMsg:    "ActivityTypeID is required",
		},
		{
			name: "empty title",
			input: CreateActivityInput{
				CourseID:       1,
				ActivityTypeID: 1,
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
				CourseID:       1,
				ActivityTypeID: 1,
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
				CourseID:       1,
				ActivityTypeID: 1,
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
