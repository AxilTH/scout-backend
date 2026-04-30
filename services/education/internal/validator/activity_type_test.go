// internal/validator/activity_type_test.go
package validator

import (
	"testing"
)

// =============================================================================
// Тесты структурной валидации (теги validate:"...")
// =============================================================================

func TestCreateActivityTypeInput_StructuralValidation(t *testing.T) {
	tests := []struct {
		name         string
		input        CreateActivityTypeInput
		expectErrors []string // ожидаемые подстроки в сообщениях об ошибках
	}{
		{
			name: "valid input - all fields correct",
			input: CreateActivityTypeInput{
				Title:       "Valid Activity Type",
				Description: stringPtr("Optional description"),
				SquadID:     1,
			},
			expectErrors: nil,
		},
		{
			name: "valid input - without description",
			input: CreateActivityTypeInput{
				Title:   "Lecture",
				SquadID: 1,
			},
			expectErrors: nil,
		},
		{
			name: "missing title - required",
			input: CreateActivityTypeInput{
				Description: stringPtr("No title"),
				SquadID:     1,
			},
			expectErrors: []string{"title is required"},
		},
		{
			name: "empty title - fails min=1",
			input: CreateActivityTypeInput{
				Title:       "",
				Description: stringPtr("Empty title"),
				SquadID:     1,
			},
			expectErrors: []string{"title is required"},
		},
		{
			name: "title too long - exceeds max=255",
			input: CreateActivityTypeInput{
				Title:       stringRepeat("A", 256),
				Description: stringPtr("Long title"),
				SquadID:     1,
			},
			expectErrors: []string{"title is out of range"},
		},
		{
			name: "missing squad_id - required",
			input: CreateActivityTypeInput{
				Title:       "No SquadID",
				Description: stringPtr("Missing squad"),
				// SquadID omitted (zero value)
			},
			expectErrors: []string{"squad ID is required"},
		},
		{
			name: "squad_id zero - fails min=1",
			input: CreateActivityTypeInput{
				Title:       "Zero SquadID",
				Description: stringPtr("Zero squad"),
				SquadID:     0,
			},
			expectErrors: []string{"squad ID"},
		},
		{
			name: "squad_id negative - fails min=1",
			input: CreateActivityTypeInput{
				Title:       "Negative SquadID",
				Description: stringPtr("Negative squad"),
				SquadID:     -1,
			},
			expectErrors: []string{"squad ID"},
		},
		{
			name: "multiple structural errors at once",
			input: CreateActivityTypeInput{
				Title:   "", // fails required + min=1
				SquadID: 0,  // fails required + min=1
				// Description omitted (optional, no error)
			},
			expectErrors: []string{"title", "squad ID"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateStruct(&tt.input)

			if len(tt.expectErrors) == 0 {
				if errs != nil {
					t.Errorf("expected no structural errors, got: %v", errs.ToHumanReadable())
				}
				return
			}

			if errs == nil {
				t.Errorf("expected structural errors %v, got none", tt.expectErrors)
				return
			}

			human := errs.ToHumanReadable()
			for _, expected := range tt.expectErrors {
				found := false
				for _, actual := range human {
					if containsIgnoreCase(actual, expected) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error containing %q, got: %v", expected, human)
				}
			}
		})
	}
}

// =============================================================================
// Тесты форматирования ошибок
// =============================================================================

func TestValidationErrors_ToHumanReadable_ActivityType(t *testing.T) {
	// Пустой список ошибок
	empty := ValidationErrors{}
	if empty.ToHumanReadable() != nil && len(empty.ToHumanReadable()) != 0 {
		t.Errorf("empty ValidationErrors should return empty slice")
	}

	// Непустой список
	nonEmpty := ValidationErrors{
		// В реальных тестах ошибки приходят из ValidateStruct
	}
	_ = nonEmpty.ToHumanReadable()
}

func TestToHumanField_ActivityType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"SquadID", "squad ID"},
		{"Title", "title"},
		{"Description", "description"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toHumanField(tt.input)
			if got != tt.expected {
				t.Errorf("toHumanField(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// =============================================================================
// Интеграционные тесты: полная валидация
// =============================================================================

func TestCreateActivityTypeInput_FullValidation(t *testing.T) {
	tests := []struct {
		name         string
		input        CreateActivityTypeInput
		expectErrors int // ожидаемое количество ошибок
	}{
		{
			name: "fully valid input",
			input: CreateActivityTypeInput{
				Title:       "Complete Valid Activity Type",
				Description: stringPtr("Full validation OK"),
				SquadID:     1,
			},
			expectErrors: 0,
		},
		{
			name: "fully valid without description",
			input: CreateActivityTypeInput{
				Title:   "Lecture",
				SquadID: 1,
			},
			expectErrors: 0,
		},
		{
			name: "multiple structural errors combined",
			input: CreateActivityTypeInput{
				Title:   "", // structural: required
				SquadID: 0,  // structural: required + min=1
			},
			expectErrors: 2, // title and squad ID
		},
		{
			name: "title max length boundary - 255 chars",
			input: CreateActivityTypeInput{
				Title:   stringRepeat("A", 255),
				SquadID: 1,
			},
			expectErrors: 0,
		},
		{
			name: "title exceeds max - 256 chars",
			input: CreateActivityTypeInput{
				Title:   stringRepeat("A", 256),
				SquadID: 1,
			},
			expectErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateStruct(&tt.input)

			if len(errs) != tt.expectErrors {
				t.Errorf("expected %d total errors, got %d: %v",
					tt.expectErrors, len(errs), errs.ToHumanReadable())
			}
		})
	}
}
