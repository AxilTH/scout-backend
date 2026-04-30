// internal/validator/activity_test.go
package validator

import (
	"testing"
	"time"
)

// =============================================================================
// Тесты структурной валидации (теги validate:"...")
// =============================================================================

func TestCreateActivityInput_StructuralValidation(t *testing.T) {
	tests := []struct {
		name         string
		input        CreateActivityInput
		expectErrors []string // ожидаемые подстроки в сообщениях об ошибках
	}{
		{
			name: "valid input - all fields correct",
			input: CreateActivityInput{
				Title:          "Valid Activity Title",
				Description:    stringPtr("Optional description"),
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       100,
				IsPublished:    true,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: nil,
		},
		{
			name: "valid input - without description",
			input: CreateActivityInput{
				Title:          "Activity No Desc",
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       50,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: nil,
		},
		{
			name: "missing title - required",
			input: CreateActivityInput{
				Description:    stringPtr("No title"),
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       50,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: []string{"title is required"},
		},
		{
			name: "empty title - fails min=1",
			input: CreateActivityInput{
				Title:          "",
				Description:    stringPtr("Empty title"),
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       50,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: []string{"title is required"},
		},
		{
			name: "title too long - exceeds max=255",
			input: CreateActivityInput{
				Title:          stringRepeat("A", 256),
				Description:    stringPtr("Long title"),
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       50,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: []string{"title is out of range"},
		},
		{
			name: "negative max_score - fails min=0",
			input: CreateActivityInput{
				Title:          "Negative Score",
				Description:    stringPtr("Bad score"),
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       -5,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: []string{"max score"},
		},
		{
			name: "missing starts_at - required",
			input: CreateActivityInput{
				Title:       "No StartsAt",
				Description: stringPtr("Missing start time"),
				// StartsAt omitted (zero value)
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       50,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: []string{"starts at is required"},
		},
		{
			name: "missing ends_at - required",
			input: CreateActivityInput{
				Title:       "No EndsAt",
				Description: stringPtr("Missing end time"),
				StartsAt:    time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				// EndsAt omitted (zero value)
				MaxScore:       50,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: []string{"ends at is required"},
		},
		{
			name: "missing course_id - required",
			input: CreateActivityInput{
				Title:       "No CourseID",
				Description: stringPtr("Missing course"),
				StartsAt:    time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:      time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:    50,
				// CourseID omitted (zero value)
				ActivityTypeID: 1,
			},
			expectErrors: []string{"course ID is required"},
		},
		{
			name: "course_id zero - fails min=1",
			input: CreateActivityInput{
				Title:          "Zero CourseID",
				Description:    stringPtr("Zero course"),
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       50,
				CourseID:       0,
				ActivityTypeID: 1,
			},
			expectErrors: []string{"course ID"},
		},
		{
			name: "missing activity_type_id - required",
			input: CreateActivityInput{
				Title:       "No ActivityTypeID",
				Description: stringPtr("Missing type"),
				StartsAt:    time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:      time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:    50,
				CourseID:    1,
				// ActivityTypeID omitted (zero value)
			},
			expectErrors: []string{"activity type ID is required"},
		},
		{
			name: "activity_type_id zero - fails min=1",
			input: CreateActivityInput{
				Title:          "Zero ActivityTypeID",
				Description:    stringPtr("Zero type"),
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       50,
				CourseID:       1,
				ActivityTypeID: 0,
			},
			expectErrors: []string{"activity type ID"},
		},
		{
			name: "multiple structural errors at once",
			input: CreateActivityInput{
				Title:          "",          // fails required + min=1
				StartsAt:       time.Time{}, // zero value
				MaxScore:       -10,         // fails min=0
				CourseID:       0,           // fails required + min=1
				ActivityTypeID: 0,           // fails required + min=1
			},
			expectErrors: []string{"title", "starts at", "max score", "course ID", "activity type ID"},
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
// Тесты бизнес-валидации (метод Validate())
// =============================================================================

func TestCreateActivityInput_BusinessValidation(t *testing.T) {
	tests := []struct {
		name         string
		input        CreateActivityInput
		expectErrors []string
	}{
		{
			name: "valid business data - all checks pass",
			input: CreateActivityInput{
				Title:          "Valid Business Activity",
				Description:    stringPtr("All checks OK"),
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       100,
				IsPublished:    true,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: nil,
		},
		{
			name: "starts_at not in UTC",
			input: CreateActivityInput{
				Title:          "Wrong TZ Starts",
				Description:    stringPtr("MSK timezone"),
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.FixedZone("MSK", 3*3600)),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       50,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: []string{"starts at must be in UTC"},
		},
		{
			name: "ends_at not in UTC",
			input: CreateActivityInput{
				Title:          "Wrong TZ Ends",
				Description:    stringPtr("PST timezone"),
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.FixedZone("PST", -8*3600)),
				MaxScore:       50,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: []string{"ends at must be in UTC"},
		},
		{
			name: "both dates not in UTC - both errors",
			input: CreateActivityInput{
				Title:          "Both Wrong TZ",
				Description:    stringPtr("MSK and PST"),
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.FixedZone("MSK", 3*3600)),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.FixedZone("PST", -8*3600)),
				MaxScore:       50,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: []string{"starts at must be in UTC", "ends at must be in UTC"},
		},
		{
			name: "ends_at equals starts_at - must be strictly after",
			input: CreateActivityInput{
				Title:          "Equal Dates",
				Description:    stringPtr("Same moment"),
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				MaxScore:       50,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: []string{"ends at must be after starts at"},
		},
		{
			name: "ends_at before starts_at",
			input: CreateActivityInput{
				Title:          "Reversed Dates",
				Description:    stringPtr("Ends before starts"),
				StartsAt:       time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				MaxScore:       50,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: []string{"ends at must be after starts at"},
		},
		{
			name: "multiple business errors at once",
			input: CreateActivityInput{
				Title:          "Multiple Errors",
				Description:    stringPtr("All wrong"),
				StartsAt:       time.Date(2026, 3, 10, 12, 0, 0, 0, time.FixedZone("MSK", 3*3600)), // not UTC
				EndsAt:         time.Date(2026, 3, 10, 8, 0, 0, 0, time.UTC),                       // before starts (in UTC terms)
				MaxScore:       50,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: []string{
				"starts at must be in UTC",
				"ends at must be after starts at",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.input.Validate()

			if len(tt.expectErrors) == 0 {
				if errs != nil {
					t.Errorf("expected no business errors, got: %v", errs.ToHumanReadable())
				}
				return
			}

			if errs == nil {
				t.Errorf("expected business errors %v, got none", tt.expectErrors)
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

func TestValidationErrors_ToHumanReadable_Activity(t *testing.T) {
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

func TestToHumanField_Activity(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"CourseID", "course ID"},
		{"ActivityTypeID", "activity type ID"},
		{"Title", "title"},
		{"StartsAt", "starts at"},
		{"EndsAt", "ends at"},
		{"MaxScore", "max score"},
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
// Интеграционные тесты: полная валидация (структурная + бизнес)
// =============================================================================

func TestCreateActivityInput_FullValidation(t *testing.T) {
	tests := []struct {
		name         string
		input        CreateActivityInput
		expectErrors int // ожидаемое количество ошибок после объединения
	}{
		{
			name: "fully valid input",
			input: CreateActivityInput{
				Title:          "Complete Valid Activity",
				Description:    stringPtr("Full validation OK"),
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       100,
				IsPublished:    true,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: 0,
		},
		{
			name: "fully valid without description",
			input: CreateActivityInput{
				Title:          "Activity No Description",
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       50,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: 0,
		},
		{
			name: "structural + business errors combined",
			input: CreateActivityInput{
				Title:          "",                                                                 // structural: required
				StartsAt:       time.Date(2026, 3, 10, 12, 0, 0, 0, time.FixedZone("MSK", 3*3600)), // business: not UTC (12:00 MSK = 09:00 UTC)
				EndsAt:         time.Date(2026, 3, 10, 8, 0, 0, 0, time.UTC),                       // business: before starts (08:00 UTC < 09:00 UTC) ✓
				MaxScore:       -5,                                                                 // structural: min=0
				CourseID:       0,                                                                  // structural: min=1
				ActivityTypeID: 0,                                                                  // structural: min=1
			},
			expectErrors: 6, // title, max_score, course_id, activity_type_id, starts_at TZ, ends_at order
		},
		{
			name: "title max length boundary - 255 chars",
			input: CreateActivityInput{
				Title:          stringRepeat("A", 255),
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       50,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: 0,
		},
		{
			name: "title exceeds max - 256 chars",
			input: CreateActivityInput{
				Title:          stringRepeat("A", 256),
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       50,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: 1,
		},
		{
			name: "max_score zero is valid",
			input: CreateActivityInput{
				Title:          "Zero Score Activity",
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       0,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: 0,
		},
		{
			name: "is_published false is valid",
			input: CreateActivityInput{
				Title:          "Unpublished Activity",
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       50,
				IsPublished:    false,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			structErrs := ValidateStruct(&tt.input)
			bizErrs := tt.input.Validate()
			allErrs := append(structErrs, bizErrs...)

			if len(allErrs) != tt.expectErrors {
				t.Errorf("expected %d total errors, got %d: %v",
					tt.expectErrors, len(allErrs), allErrs.ToHumanReadable())
			}
		})
	}
}

// =============================================================================
// Граничные случаи
// =============================================================================

func TestCreateActivityInput_EdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		input        CreateActivityInput
		expectErrors int
		description  string
	}{
		{
			name: "activity duration 1 second",
			input: CreateActivityInput{
				Title:          "1 Second Activity",
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 10, 0, 1, 0, time.UTC),
				MaxScore:       1,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: 0,
			description:  "Минимальная длительность активности",
		},
		{
			name: "activity spanning multiple days",
			input: CreateActivityInput{
				Title:          "Long Activity",
				StartsAt:       time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 31, 23, 59, 59, 0, time.UTC),
				MaxScore:       1000,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: 0,
			description:  "Активность на весь месяц",
		},
		{
			name: "very high max_score",
			input: CreateActivityInput{
				Title:          "High Score Activity",
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       999999,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: 0,
			description:  "Большой максимальный балл",
		},
		{
			name: "long description",
			input: CreateActivityInput{
				Title:          "Activity with Long Description",
				Description:    stringPtr(stringRepeat("Lorem ipsum ", 100)),
				StartsAt:       time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC),
				EndsAt:         time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
				MaxScore:       50,
				CourseID:       1,
				ActivityTypeID: 1,
			},
			expectErrors: 0,
			description:  "Длинное описание (нет валидации длины)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			structErrs := ValidateStruct(&tt.input)
			bizErrs := tt.input.Validate()
			allErrs := append(structErrs, bizErrs...)

			if len(allErrs) != tt.expectErrors {
				t.Errorf("[%s] expected %d errors, got %d: %v",
					tt.description, tt.expectErrors, len(allErrs), allErrs.ToHumanReadable())
			}
		})
	}
}
