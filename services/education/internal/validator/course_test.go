package validator

import (
	"testing"
	"time"
)

func TestCreateCourseInput_StructuralValidation(t *testing.T) {
	v := NewValidator()

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
			errorMsg:    "Invalid UUID format for squad_id",
		},
		{
			name: "year too small",
			input: CreateCourseInput{
				SquadID:  "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				Year:     1999,
				Title:    "Old Year",
				StartsAt: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:   time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
			},
			expectError: true,
			errorMsg:    "Year out of range",
		},
		{
			name: "year too large",
			input: CreateCourseInput{
				SquadID:  "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				Year:     2101,
				Title:    "Future Year",
				StartsAt: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:   time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
			},
			expectError: true,
			errorMsg:    "Year out of range",
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
		{
			name: "title too long",
			input: CreateCourseInput{
				SquadID:  "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				Year:     2026,
				Title:    string(make([]byte, 256)), // 256 символов
				StartsAt: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:   time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
			},
			expectError: true,
			errorMsg:    "Title is too large",
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

func TestValidateSquadID(t *testing.T) {
	tests := []struct {
		name        string
		squadID     string
		expectError bool
	}{
		{
			name:        "valid uuid",
			squadID:     "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			expectError: false,
		},
		{
			name:        "invalid uuid",
			squadID:     "not-a-uuid",
			expectError: true,
		},
		{
			name:        "empty string",
			squadID:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSquadID(tt.squadID)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateSquadID(%q) error = %v, expectError = %v", tt.squadID, err, tt.expectError)
			}
		})
	}
}

func TestFormatError(t *testing.T) {
	v := NewValidator()

	input := CreateCourseInput{SquadID: "not-a-uuid"}
	err := v.Struct(&input)
	if err == nil {
		t.Fatal("expected validation error")
	}

	formatted := FormatError(err)
	expected := "Invalid UUID format for squad_id"
	if formatted != expected {
		t.Errorf("expected %q, got %q", expected, formatted)
	}
}
