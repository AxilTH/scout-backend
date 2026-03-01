package validator

import (
	"testing"
	"time"
)

func TestCreateCourseInput_StructuralValidation(t *testing.T) {
	v := GetValidator()

	tests := []struct {
		name        string
		input       CreateCourseInput
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid input",
			input: CreateCourseInput{
				SquadID:  "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				Year:     2026,
				Title:    "Valid Course",
				StartsAt: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:   time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
			},
			expectError: false,
		},
		{
			name: "missing squad_id",
			input: CreateCourseInput{
				Year:     2026,
				Title:    "No SquadID",
				StartsAt: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:   time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
			},
			expectError: true,
			errorMsg:    "SquadID is required",
		},
		{
			name: "invalid uuid format",
			input: CreateCourseInput{
				SquadID:  "not-a-uuid",
				Year:     2026,
				Title:    "Invalid UUID",
				StartsAt: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:   time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
			},
			expectError: true,
			errorMsg:    "Invalid UUID format for SquadID",
		},
		{
			name: "empty title",
			input: CreateCourseInput{
				SquadID:  "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				Year:     2026,
				Title:    "",
				StartsAt: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:   time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
			},
			expectError: true,
			errorMsg:    "Title must not be empty",
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

func TestCreateCourseInput_BusinessValidation(t *testing.T) {
	tests := []struct {
		name        string
		startsAt    time.Time
		endsAt      time.Time
		expectError bool
	}{
		{
			name:        "ends_at after starts_at",
			startsAt:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
			endsAt:      time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "ends_at equals starts_at",
			startsAt:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
			endsAt:      time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
			expectError: true,
		},
		{
			name:        "ends_at before starts_at",
			startsAt:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
			endsAt:      time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := CreateCourseInput{
				SquadID:  "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				Year:     2026,
				Title:    "Test Course",
				StartsAt: tt.startsAt,
				EndsAt:   tt.endsAt,
			}
			err := input.Validate()
			if (err != nil) != tt.expectError {
				t.Errorf("Validate() error = %v, expectError = %v", err, tt.expectError)
			}
		})
	}
}