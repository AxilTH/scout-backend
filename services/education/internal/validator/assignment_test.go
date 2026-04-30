package validator

// import (
// 	"testing"
// 	"time"
// )

// func int64Ptr(i int64) *int64 {
// 	return &i
// }

// func TestCreateAssignmentInput_StructuralValidation(t *testing.T) {
// 	v := GetValidator()

// 	tests := []struct {
// 		name        string
// 		input       CreateAssignmentInput
// 		expectError bool
// 		errorMsg    string
// 	}{
// 		{
// 			name: "valid input",
// 			input: CreateAssignmentInput{
// 				ActivityID:       int64Ptr(1),
// 				AssignmentTypeID: 1,
// 				Title:            "Педагогическое эссе",
// 				Deadline:         time.Now().Add(7 * 24 * time.Hour),
// 				MaxScore:         20,
// 			},
// 			expectError: false,
// 		},
// 		{
// 			name: "valid input without activity_id",
// 			input: CreateAssignmentInput{
// 				AssignmentTypeID: 2,
// 				Title:            "Итоговый проект",
// 				Deadline:         time.Now().Add(30 * 24 * time.Hour),
// 				MaxScore:         100,
// 			},
// 			expectError: false,
// 		},
// 		{
// 			name: "missing assignment_type_id",
// 			input: CreateAssignmentInput{
// 				ActivityID: int64Ptr(1),
// 				Title:      "Без типа",
// 				Deadline:   time.Now().Add(7 * 24 * time.Hour),
// 				MaxScore:   20,
// 			},
// 			expectError: true,
// 			errorMsg:    "AssignmentTypeID is required",
// 		},
// 		{
// 			name: "empty title",
// 			input: CreateAssignmentInput{
// 				ActivityID:       int64Ptr(1),
// 				AssignmentTypeID: 1,
// 				Title:            "",
// 				Deadline:         time.Now().Add(7 * 24 * time.Hour),
// 				MaxScore:         20,
// 			},
// 			expectError: true,
// 			errorMsg:    "Title must not be empty",
// 		},
// 		{
// 			name: "negative max_score",
// 			input: CreateAssignmentInput{
// 				ActivityID:       int64Ptr(1),
// 				AssignmentTypeID: 1,
// 				Title:            "Отрицательный балл",
// 				Deadline:         time.Now().Add(7 * 24 * time.Hour),
// 				MaxScore:         -5,
// 			},
// 			expectError: true,
// 			errorMsg:    "MaxScore is too small",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := v.Struct(&tt.input)
// 			if tt.expectError {
// 				if err == nil {
// 					t.Errorf("expected validation error, got none")
// 				} else {
// 					formatted := FormatError(err)
// 					if formatted != tt.errorMsg {
// 						t.Errorf("expected error %q, got %q", tt.errorMsg, formatted)
// 					}
// 				}
// 			} else {
// 				if err != nil {
// 					t.Errorf("unexpected validation error: %v", err)
// 				}
// 			}
// 		})
// 	}
// }

// func TestCreateAssignmentInput_BusinessValidation(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		deadline    time.Time
// 		expectError bool
// 	}{
// 		{
// 			name:        "deadline in future",
// 			deadline:    time.Now().Add(7 * 24 * time.Hour),
// 			expectError: false,
// 		},
// 		{
// 			name:        "deadline in past",
// 			deadline:    time.Now().Add(-24 * time.Hour),
// 			expectError: true,
// 		},
// 		{
// 			name:        "deadline is now",
// 			deadline:    time.Now(),
// 			expectError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			input := CreateAssignmentInput{
// 				ActivityID:       int64Ptr(1),
// 				AssignmentTypeID: 1,
// 				Title:            "Эссе",
// 				Deadline:         tt.deadline,
// 				MaxScore:         20,
// 			}
// 			err := input.Validate()
// 			if (err != nil) != tt.expectError {
// 				t.Errorf("Validate() error = %v, expectError = %v", err, tt.expectError)
// 			}
// 		})
// 	}
// }
