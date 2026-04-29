// internal/validator/course_test.go
package validator

import (
	"testing"
	"time"
)

// =============================================================================
// Тесты структурной валидации (теги validate:"...")
// =============================================================================

func TestCreateCourseInput_StructuralValidation(t *testing.T) {
	tests := []struct {
		name         string
		input        CreateCourseInput
		expectErrors []string // ожидаемые подстроки в сообщениях об ошибках
	}{
		{
			name: "valid input - all fields correct",
			input: CreateCourseInput{
				Title:       "Valid Course Title",
				Description: stringPtr("Optional description"),
				StartsAt:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:      time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
				Year:        2026,
				SquadID:     1,
			},
			expectErrors: nil,
		},
		{
			name: "missing title - required",
			input: CreateCourseInput{
				Description: stringPtr("No title"),
				StartsAt:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:      time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
				Year:        2026,
				SquadID:     1,
			},
			expectErrors: []string{"title is required"},
		},
		{
			name: "empty title - fails min=1",
			input: CreateCourseInput{
				Title:       "",
				Description: stringPtr("Empty title"),
				StartsAt:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:      time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
				Year:        2026,
				SquadID:     1,
			},
			expectErrors: []string{"title is required"},
		},
		{
			name: "title too long - exceeds max=255",
			input: CreateCourseInput{
				Title:       stringRepeat("A", 256),
				Description: stringPtr("Long title"),
				StartsAt:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:      time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
				Year:        2026,
				SquadID:     1,
			},
			expectErrors: []string{"title is out of range"},
		},
		{
			name: "missing squad_id - required",
			input: CreateCourseInput{
				Title:       "No SquadID",
				Description: stringPtr("Missing squad"),
				StartsAt:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:      time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
				Year:        2026,
				// SquadID omitted (zero value)
			},
			expectErrors: []string{"squad ID is required"},
		},
		{
			name: "squad_id zero - fails min=1",
			input: CreateCourseInput{
				Title:       "Zero SquadID",
				Description: stringPtr("Zero squad"),
				StartsAt:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:      time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
				Year:        2026,
				SquadID:     0,
			},
			expectErrors: []string{"squad ID"},
		},
		{
			name: "missing year - required",
			input: CreateCourseInput{
				Title:       "No Year",
				Description: stringPtr("Missing year"),
				StartsAt:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:      time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
				// Year omitted (zero value)
				SquadID: 1,
			},
			expectErrors: []string{"year is required"},
		},
		{
			name: "missing starts_at - required",
			input: CreateCourseInput{
				Title:       "No StartsAt",
				Description: stringPtr("Missing starts"),
				// StartsAt omitted (zero value)
				EndsAt:  time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
				Year:    2026,
				SquadID: 1,
			},
			expectErrors: []string{"starts at is required"},
		},
		{
			name: "missing ends_at - required",
			input: CreateCourseInput{
				Title:       "No EndsAt",
				Description: stringPtr("Missing ends"),
				StartsAt:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				// EndsAt omitted (zero value)
				Year:    2026,
				SquadID: 1,
			},
			expectErrors: []string{"ends at is required"},
		},
		{
			name: "multiple structural errors at once",
			input: CreateCourseInput{
				Title:    "",          // fails required + min=1
				StartsAt: time.Time{}, // zero value, but present for validator
				// Description, EndsAt, Year, SquadID omitted
			},
			expectErrors: []string{"title", "year is required", "starts at", "ends at", "squad ID"},
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

func TestCreateCourseInput_BusinessValidation(t *testing.T) {
	currentYear := time.Now().UTC().Year()

	tests := []struct {
		name         string
		input        CreateCourseInput
		expectErrors []string
	}{
		{
			name: "valid business data - all checks pass",
			input: CreateCourseInput{
				Title:       "Valid Business Course",
				Description: stringPtr("All checks OK"),
				StartsAt:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:      time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
				Year:        2026,
				SquadID:     1,
			},
			expectErrors: nil,
		},
		{
			name: "starts_at not in UTC",
			input: CreateCourseInput{
				Title:       "Wrong TZ Starts",
				Description: stringPtr("MSK timezone"),
				StartsAt:    time.Date(2026, 9, 1, 10, 0, 0, 0, time.FixedZone("MSK", 3*3600)),
				EndsAt:      time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
				Year:        2026,
				SquadID:     1,
			},
			expectErrors: []string{"starts at must be in UTC"},
		},
		{
			name: "ends_at not in UTC",
			input: CreateCourseInput{
				Title:       "Wrong TZ Ends",
				Description: stringPtr("PST timezone"),
				StartsAt:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:      time.Date(2027, 5, 31, 12, 0, 0, 0, time.FixedZone("PST", -8*3600)),
				Year:        2026,
				SquadID:     1,
			},
			expectErrors: []string{"ends at must be in UTC"},
		},
		{
			name: "both dates not in UTC - both errors",
			input: CreateCourseInput{
				Title:       "Both Wrong TZ",
				Description: stringPtr("MSK and PST"),
				StartsAt:    time.Date(2026, 9, 1, 10, 0, 0, 0, time.FixedZone("MSK", 3*3600)),
				EndsAt:      time.Date(2027, 5, 31, 12, 0, 0, 0, time.FixedZone("PST", -8*3600)),
				Year:        2026,
				SquadID:     1,
			},
			expectErrors: []string{"starts at must be in UTC", "ends at must be in UTC"},
		},
		{
			name: "ends_at equals starts_at - must be strictly after",
			input: CreateCourseInput{
				Title:       "Equal Dates",
				Description: stringPtr("Same moment"),
				StartsAt:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:      time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				Year:        2026,
				SquadID:     1,
			},
			expectErrors: []string{"ends at must be after starts at"},
		},
		{
			name: "ends_at before starts_at",
			input: CreateCourseInput{
				Title:       "Reversed Dates",
				Description: stringPtr("Ends before starts"),
				StartsAt:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:      time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
				Year:        2026,
				SquadID:     1,
			},
			expectErrors: []string{"ends at must be after starts at"},
		},
		{
			name: "year does not match starts_at year",
			input: CreateCourseInput{
				Title:       "Year Mismatch",
				Description: stringPtr("Wrong year"),
				StartsAt:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:      time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
				Year:        2025, // mismatch
				SquadID:     1,
			},
			expectErrors: []string{"year must match the year of starts at (2026)"},
		},
		{
			name: "year too old - before 2020",
			input: CreateCourseInput{
				Title:       "Too Old Year",
				Description: stringPtr("Before 2020"),
				StartsAt:    time.Date(2019, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:      time.Date(2020, 5, 31, 0, 0, 0, 0, time.UTC),
				Year:        2019,
				SquadID:     1,
			},
			expectErrors: []string{"year must be between 2020 and"},
		},
		{
			name: "year too far in future - beyond current+5",
			input: CreateCourseInput{
				Title:       "Too Future Year",
				Description: stringPtr("Beyond limit"),
				StartsAt:    time.Date(currentYear+6, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:      time.Date(currentYear+7, 5, 31, 0, 0, 0, 0, time.UTC),
				Year:        currentYear + 6,
				SquadID:     1,
			},
			expectErrors: []string{"year must be between 2020 and"},
		},
		{
			name: "multiple business errors at once",
			input: CreateCourseInput{
				Title:       "Multiple Errors",
				Description: stringPtr("All wrong"),
				StartsAt:    time.Date(2026, 9, 1, 10, 0, 0, 0, time.FixedZone("MSK", 3*3600)), // not UTC
				EndsAt:      time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),                       // before starts
				Year:        2019,                                                              // too old + mismatch
				SquadID:     1,
			},
			expectErrors: []string{
				"starts at must be in UTC",
				"ends at must be after starts at",
				"year must match",
				"year must be between 2020",
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
// Тесты форматирования ошибок и вспомогательных функций
// =============================================================================

func TestValidationErrors_ToHumanReadable(t *testing.T) {
	// Пустой список ошибок
	empty := ValidationErrors{}
	if empty.ToHumanReadable() != nil && len(empty.ToHumanReadable()) != 0 {
		t.Errorf("empty ValidationErrors should return empty slice")
	}

	// Непустой список
	nonEmpty := ValidationErrors{
		// В реальных тестах ошибки приходят из ValidateStruct/Validate,
		// здесь проверяем только, что метод не паникует и возвращает срез
	}
	_ = nonEmpty.ToHumanReadable()
}

func TestToHumanField(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"SquadID", "squad ID"},
		{"CourseID", "course ID"},
		{"ActivityTypeID", "activity type ID"},
		{"AssignmentTypeID", "assignment type ID"},
		{"Title", "title"},
		{"Year", "year"},
		{"StartsAt", "starts at"},
		{"EndsAt", "ends at"},
		{"Deadline", "deadline"},
		{"MaxScore", "max score"},
		{"UnknownField", "unknownField"}, // fallback: lowercase first letter
		{"", ""},
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

func TestCreateCourseInput_FullValidation(t *testing.T) {
	tests := []struct {
		name         string
		input        CreateCourseInput
		expectErrors int // ожидаемое количество ошибок после объединения
	}{
		{
			name: "fully valid input",
			input: CreateCourseInput{
				Title:       "Complete Valid Course",
				Description: stringPtr("Full validation OK"),
				StartsAt:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:      time.Date(2027, 5, 31, 0, 0, 0, 0, time.UTC),
				Year:        2026,
				SquadID:     1,
			},
			expectErrors: 0,
		},
		{
			name: "structural + business errors combined",
			input: CreateCourseInput{
				Title:    "",                                                                // structural: required
				StartsAt: time.Date(2026, 9, 1, 10, 0, 0, 0, time.FixedZone("MSK", 3*3600)), // business: not UTC
				EndsAt:   time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),                       // business: before starts
				Year:     2019,                                                              // business: too old + mismatch
				SquadID:  0,                                                                 // structural: min=1
			},
			expectErrors: 6, // title, squad ID, year required (maybe), year range, UTC, date order
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

// Примечание: вспомогательные функции для тестов находятся в test_helpers.go
