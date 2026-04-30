// internal/testutil/testutil.go
package testutil

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/AxilTH/scout-backend/services/education/internal/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// TestConfig содержит конфигурацию для тестов
type TestConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// DefaultTestConfig возвращает конфигурацию по умолчанию для тестов
// Использует переменные окружения или дефолтные значения
func DefaultTestConfig() *TestConfig {
	return &TestConfig{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "education"),
	}
}

// getEnv возвращает значение переменной окружения или дефолт
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SetupTestDB создает тестовую базу данных и возвращает подключение
// Если база данных недоступна, пропускает тесты с t.Skip()
func SetupTestDB(t *testing.T, config *TestConfig) *sqlx.DB {
	if config == nil {
		config = DefaultTestConfig()
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		t.Skipf("Skipping test: database not available: %v", err)
		return nil
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		t.Skipf("Skipping test: cannot ping database: %v", err)
		return nil
	}

	return db
}

// TeardownTestDB закрывает подключение к тестовой базе данных
func TeardownTestDB(t *testing.T, db *sqlx.DB) {
	if err := db.Close(); err != nil {
		t.Logf("Failed to close test database: %v", err)
	}
}

// CreateTestCourse создает тестовый курс с указанным squad_id
// Примечание: squad_id должен существовать в базе данных (внешний ключ)
func CreateTestCourse(t *testing.T, db *sqlx.DB, squadID int64) *model.Course {
	course := &model.Course{
		Title:       fmt.Sprintf("Test Course %d", time.Now().UnixNano()),
		Description: strPtr("Test course description"),
		StartsAt:    time.Now().UTC(),
		EndsAt:      time.Now().UTC().Add(30 * 24 * time.Hour),
		Year:        time.Now().UTC().Year(),
		IsActive:    true,
		SquadID:     squadID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	query := `
		INSERT INTO courses (title, description, starts_at, ends_at, year, is_active, created_at, updated_at, squad_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	err := db.QueryRow(
		query,
		course.Title,
		course.Description,
		course.StartsAt,
		course.EndsAt,
		course.Year,
		course.IsActive,
		course.CreatedAt,
		course.UpdatedAt,
		course.SquadID,
	).Scan(&course.ID)

	if err != nil {
		t.Fatalf("Failed to create test course: %v", err)
	}

	return course
}

// CreateTestActivityType создает тестовый тип активности
// squadID может быть nil для глобальных типов активностей
func CreateTestActivityType(t *testing.T, db *sqlx.DB, squadID *int64) *model.ActivityType {
	activityType := &model.ActivityType{
		Title:       fmt.Sprintf("Test Activity Type %d", time.Now().UnixNano()),
		Description: strPtr("Test activity type description"),
		SquadID:     squadID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	query := `
		INSERT INTO activity_types (title, description, created_at, updated_at, squad_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := db.QueryRow(
		query,
		activityType.Title,
		activityType.Description,
		activityType.CreatedAt,
		activityType.UpdatedAt,
		activityType.SquadID,
	).Scan(&activityType.ID)

	if err != nil {
		t.Fatalf("Failed to create test activity type: %v", err)
	}

	return activityType
}

// CreateTestActivity создает тестовую активность
func CreateTestActivity(t *testing.T, db *sqlx.DB, courseID, activityTypeID int64) *model.Activity {
	activity := &model.Activity{
		Title:          fmt.Sprintf("Test Activity %d", time.Now().UnixNano()),
		Description:    strPtr("Test activity description"),
		StartsAt:       time.Now().UTC(),
		EndsAt:         time.Now().UTC().Add(2 * time.Hour),
		MaxScore:       100,
		IsPublished:    true,
		CourseID:       courseID,
		ActivityTypeID: activityTypeID,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	query := `
		INSERT INTO activities (title, description, starts_at, ends_at, max_score, is_published, created_at, updated_at, course_id, activity_type_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	err := db.QueryRow(
		query,
		activity.Title,
		activity.Description,
		activity.StartsAt,
		activity.EndsAt,
		activity.MaxScore,
		activity.IsPublished,
		activity.CreatedAt,
		activity.UpdatedAt,
		activity.CourseID,
		activity.ActivityTypeID,
	).Scan(&activity.ID)

	if err != nil {
		t.Fatalf("Failed to create test activity: %v", err)
	}

	return activity
}

// CleanupTestData удаляет тестовые данные
// Примечание: не удаляет данные из таблицы squads, так как она управляется другим сервисом
func CleanupTestData(t *testing.T, db *sqlx.DB) {
	// Удаляем в правильном порядке из-за внешних ключей
	_, _ = db.Exec("DELETE FROM activities")
	_, _ = db.Exec("DELETE FROM courses")
	_, _ = db.Exec("DELETE FROM activity_types WHERE squad_id IS NOT NULL")
}

// GenerateTestJWT генерирует тестовый JWT токен
// В реальном проекте это должно вызывать Auth Service
func GenerateTestJWT(userID string, squadID int64) string {
	// Заглушка - в реальном проекте здесь должен быть вызов Auth Service
	// или генерация токена с использованием того же секрета, что и в middleware
	return fmt.Sprintf("test.jwt.token.%s.%d", userID, squadID)
}

// strPtr возвращает указатель на строку
func strPtr(s string) *string {
	return &s
}

// AssertNoError проверяет, что ошибки нет
func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

// AssertError проверяет, что ошибка есть
func AssertError(t *testing.T, err error) {
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// AssertEqual проверяет равенство значений
func AssertEqual[T comparable](t *testing.T, expected, actual T) {
	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

// AssertNotNil проверяет, что значение не nil
func AssertNotNil(t *testing.T, value interface{}) {
	if value == nil {
		t.Fatal("Expected value to be non-nil")
	}
}

// AssertNil проверяет, что значение nil
func AssertNil(t *testing.T, value interface{}) {
	if value != nil {
		t.Errorf("Expected value to be nil, got %v", value)
	}
}