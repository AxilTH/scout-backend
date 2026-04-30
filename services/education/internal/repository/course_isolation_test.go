// internal/repository/course_isolation_test.go
package repository_test

import (
	"testing"

	"github.com/AxilTH/scout-backend/services/education/internal/model"
	"github.com/AxilTH/scout-backend/services/education/internal/repository"
	"github.com/AxilTH/scout-backend/services/education/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCourseIsolation тестирует изоляцию данных между отрядами
// Примечание: Использует фиктивные squad_id, которые должны существовать в базе данных
func TestCourseIsolation(t *testing.T) {
	// Настройка тестовой БД
	db := testutil.SetupTestDB(t, testutil.DefaultTestConfig())
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestData(t, db)

	// Используем фиктивные squad_id для тестирования изоляции
	// В реальном проекте эти squad_id должны существовать в таблице squads
	squadIDA := int64(1)
	squadIDB := int64(2)

	// Создаем курс в отряде A
	courseA := testutil.CreateTestCourse(t, db, squadIDA)

	// Создаем репозиторий
	repo := repository.NewCourseRepository(db)

	// ТЕСТ 1: Отряд A видит свой курс
	t.Run("SquadA can see own course", func(t *testing.T) {
		courses, err := repo.GetAll(squadIDA)
		require.NoError(t, err)
		assert.Len(t, courses, 1)
		assert.Equal(t, courseA.ID, courses[0].ID)
	})

	// ТЕСТ 2: Отряд B не видит курс отряда A
	t.Run("SquadB cannot see SquadA course", func(t *testing.T) {
		courses, err := repo.GetAll(squadIDB)
		require.NoError(t, err)
		assert.Empty(t, courses, "Squad B should not see Squad A's courses")
	})

	// ТЕСТ 3: Отряд A может получить свой курс по ID
	t.Run("SquadA can get own course by ID", func(t *testing.T) {
		course, err := repo.GetByID(courseA.ID, squadIDA)
		require.NoError(t, err)
		require.NotNil(t, course)
		assert.Equal(t, courseA.ID, course.ID)
		assert.Equal(t, courseA.Title, course.Title)
	})

	// ТЕСТ 4: Отряд B не может получить курс отряда A по ID
	t.Run("SquadB cannot get SquadA course by ID", func(t *testing.T) {
		course, err := repo.GetByID(courseA.ID, squadIDB)
		require.NoError(t, err)
		assert.Nil(t, course, "Squad B should not be able to get Squad A's course")
	})

	// ТЕСТ 5: Отряд B не может обновить курс отряда A
	t.Run("SquadB cannot update SquadA course", func(t *testing.T) {
		updatedCourse := &model.Course{
			ID:          courseA.ID,
			Title:       "Hacked Title",
			Description: courseA.Description,
			StartsAt:    courseA.StartsAt,
			EndsAt:      courseA.EndsAt,
			Year:        courseA.Year,
			IsActive:    courseA.IsActive,
			SquadID:     squadIDB, // Пытаемся изменить принадлежность
		}

		err := repo.Update(updatedCourse)
		require.NoError(t, err)

		// Проверяем, что курс не изменился
		originalCourse, err := repo.GetByID(courseA.ID, squadIDA)
		require.NoError(t, err)
		require.NotNil(t, originalCourse)
		assert.Equal(t, courseA.Title, originalCourse.Title, "Course title should not be changed")
		assert.Equal(t, squadIDA, originalCourse.SquadID, "Course should still belong to Squad A")
	})

	// ТЕСТ 6: Отряд B не может удалить курс отряда A
	t.Run("SquadB cannot delete SquadA course", func(t *testing.T) {
		err := repo.Delete(courseA.ID, squadIDB)
		require.NoError(t, err)

		// Проверяем, что курс все еще существует для отряда A
		course, err := repo.GetByID(courseA.ID, squadIDA)
		require.NoError(t, err)
		assert.NotNil(t, course, "Course should still exist for Squad A")
	})

	// ТЕСТ 7: Отряд A может удалить свой курс
	t.Run("SquadA can delete own course", func(t *testing.T) {
		err := repo.Delete(courseA.ID, squadIDA)
		require.NoError(t, err)

		// Проверяем, что курс удален
		course, err := repo.GetByID(courseA.ID, squadIDA)
		require.NoError(t, err)
		assert.Nil(t, course, "Course should be deleted")
	})

	// ТЕСТ 8: ExistsForSquad корректно проверяет принадлежность
	t.Run("ExistsForSquad correctly checks ownership", func(t *testing.T) {
		// Создаем новый курс для теста
		courseA2 := testutil.CreateTestCourse(t, db, squadIDA)

		// Проверяем для отряда A
		exists, err := repo.ExistsForSquad(courseA2.ID, squadIDA)
		require.NoError(t, err)
		assert.True(t, exists, "Course should exist for Squad A")

		// Проверяем для отряда B
		exists, err = repo.ExistsForSquad(courseA2.ID, squadIDB)
		require.NoError(t, err)
		assert.False(t, exists, "Course should not exist for Squad B")
	})
}

// TestActivityIsolation тестирует изоляцию активностей между отрядами
func TestActivityIsolation(t *testing.T) {
	// Настройка тестовой БД
	db := testutil.SetupTestDB(t, testutil.DefaultTestConfig())
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestData(t, db)

	// Используем фиктивные squad_id
	squadIDA := int64(1)
	squadIDB := int64(2)

	// Создаем курс для отряда A
	courseA := testutil.CreateTestCourse(t, db, squadIDA)

	// Создаем тип активности для отряда A
	activityTypeA := testutil.CreateTestActivityType(t, db, &squadIDA)

	// Создаем активность в курсе отряда A
	activityA := testutil.CreateTestActivity(t, db, courseA.ID, activityTypeA.ID)

	// Создаем репозиторий
	repo := repository.NewActivityRepository(db)

	// ТЕСТ 1: Отряд A видит активности своего курса
	t.Run("SquadA can see own course activities", func(t *testing.T) {
		activities, err := repo.GetAll(courseA.ID, squadIDA)
		require.NoError(t, err)
		assert.Len(t, activities, 1)
		assert.Equal(t, activityA.ID, activities[0].ID)
	})

	// ТЕСТ 2: Отряд B не видит активности курса отряда A
	t.Run("SquadB cannot see SquadA course activities", func(t *testing.T) {
		activities, err := repo.GetAll(courseA.ID, squadIDB)
		require.NoError(t, err)
		assert.Empty(t, activities, "Squad B should not see Squad A's activities")
	})

	// ТЕСТ 3: Отряд A может получить свою активность по ID
	t.Run("SquadA can get own activity by ID", func(t *testing.T) {
		activity, err := repo.GetByID(activityA.ID, courseA.ID, squadIDA)
		require.NoError(t, err)
		require.NotNil(t, activity)
		assert.Equal(t, activityA.ID, activity.ID)
		assert.Equal(t, activityA.Title, activity.Title)
	})

	// ТЕСТ 4: Отряд B не может получить активность отряда A по ID
	t.Run("SquadB cannot get SquadA activity by ID", func(t *testing.T) {
		activity, err := repo.GetByID(activityA.ID, courseA.ID, squadIDB)
		require.NoError(t, err)
		assert.Nil(t, activity, "Squad B should not be able to get Squad A's activity")
	})

	// ТЕСТ 5: Отряд B не может удалить активность отряда A
	t.Run("SquadB cannot delete SquadA activity", func(t *testing.T) {
		err := repo.Delete(activityA.ID, courseA.ID, squadIDB)
		require.NoError(t, err)

		// Проверяем, что активность все еще существует для отряда A
		activity, err := repo.GetByID(activityA.ID, courseA.ID, squadIDA)
		require.NoError(t, err)
		assert.NotNil(t, activity, "Activity should still exist for Squad A")
	})

	// ТЕСТ 6: Отряд A может удалить свою активность
	t.Run("SquadA can delete own activity", func(t *testing.T) {
		err := repo.Delete(activityA.ID, courseA.ID, squadIDA)
		require.NoError(t, err)

		// Проверяем, что активность удалена
		activity, err := repo.GetByID(activityA.ID, courseA.ID, squadIDA)
		require.NoError(t, err)
		assert.Nil(t, activity, "Activity should be deleted")
	})
}