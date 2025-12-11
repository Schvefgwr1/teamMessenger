//go:build integration
// +build integration

package integration

import (
	"testing"

	"taskService/internal/controllers"
	"taskService/internal/models"
	"taskService/internal/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTaskStatusController_Create_Integration тестирует создание статуса задачи с реальной БД
func TestTaskStatusController_Create_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - настройка реальной БД
	db := setupTestDB(t)
	repo := repositories.NewTaskStatusRepository(db)
	controller := controllers.NewTaskStatusController(repo)

	statusName := "test_in_progress"

	// Act - выполнение реального сценария
	status, err := controller.Create(statusName)

	// Assert - проверка реального результата
	require.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, statusName, status.Name)
	assert.NotZero(t, status.ID)

	// Проверяем, что статус реально сохранен в БД
	var savedStatus models.TaskStatus
	err = db.Where("name = ?", statusName).First(&savedStatus).Error
	require.NoError(t, err)
	assert.Equal(t, status.ID, savedStatus.ID)
	assert.Equal(t, statusName, savedStatus.Name)
}

// TestTaskStatusController_Create_Integration_Duplicate тестирует обработку ошибки при создании дубликата статуса
func TestTaskStatusController_Create_Integration_Duplicate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	repo := repositories.NewTaskStatusRepository(db)
	controller := controllers.NewTaskStatusController(repo)

	statusName := "test_duplicate_status"

	// Создаем первый статус
	_, err := controller.Create(statusName)
	require.NoError(t, err)

	// Act - пытаемся создать дубликат
	status, err := controller.Create(statusName)

	// Assert
	require.Error(t, err)
	assert.Nil(t, status)
	assert.Contains(t, err.Error(), "already exists")
}

// TestTaskStatusController_GetByID_Integration тестирует получение статуса по ID с реальной БД
func TestTaskStatusController_GetByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	repo := repositories.NewTaskStatusRepository(db)
	controller := controllers.NewTaskStatusController(repo)

	// Создаем статус для теста
	statusName := "test_get_by_id"
	createdStatus, err := controller.Create(statusName)
	require.NoError(t, err)

	// Act
	status, err := controller.GetByID(createdStatus.ID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, createdStatus.ID, status.ID)
	assert.Equal(t, statusName, status.Name)
}

// TestTaskStatusController_GetByID_Integration_NotFound тестирует обработку ошибки, когда статус не найден
func TestTaskStatusController_GetByID_Integration_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	repo := repositories.NewTaskStatusRepository(db)
	controller := controllers.NewTaskStatusController(repo)

	nonExistentID := 99999

	// Act
	status, err := controller.GetByID(nonExistentID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, status)
	assert.Contains(t, err.Error(), "task status")
}

// TestTaskStatusController_GetAll_Integration тестирует получение всех статусов с реальной БД
func TestTaskStatusController_GetAll_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	repo := repositories.NewTaskStatusRepository(db)
	controller := controllers.NewTaskStatusController(repo)

	// Создаем несколько статусов для теста
	statusNames := []string{"test_status_1", "test_status_2", "test_status_3"}
	for _, name := range statusNames {
		_, err := controller.Create(name)
		require.NoError(t, err)
	}

	// Act
	statuses, err := controller.GetAll()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, statuses)
	assert.GreaterOrEqual(t, len(statuses), len(statusNames))

	// Проверяем, что созданные статусы присутствуют
	statusMap := make(map[string]bool)
	for _, status := range statuses {
		statusMap[status.Name] = true
	}

	for _, name := range statusNames {
		assert.True(t, statusMap[name], "Status %s should be present", name)
	}

	// Проверяем, что предустановленные статусы тоже присутствуют
	assert.True(t, statusMap["created"], "Predefined status 'created' should be present")
	assert.True(t, statusMap["canseled"], "Predefined status 'canseled' should be present")
}

// TestTaskStatusController_DeleteByID_Integration тестирует удаление статуса с реальной БД
func TestTaskStatusController_DeleteByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	repo := repositories.NewTaskStatusRepository(db)
	controller := controllers.NewTaskStatusController(repo)

	// Создаем статус для удаления
	statusName := "test_delete_status"
	createdStatus, err := controller.Create(statusName)
	require.NoError(t, err)

	// Проверяем, что статус существует
	_, err = controller.GetByID(createdStatus.ID)
	require.NoError(t, err)

	// Act
	err = controller.DeleteByID(createdStatus.ID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что статус реально удален из БД
	var deletedStatus models.TaskStatus
	err = db.First(&deletedStatus, createdStatus.ID).Error
	assert.Error(t, err, "Status should be deleted from database")
}

// TestTaskStatusController_DeleteByID_Integration_NotFound тестирует обработку ошибки при удалении несуществующего статуса
func TestTaskStatusController_DeleteByID_Integration_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	repo := repositories.NewTaskStatusRepository(db)
	controller := controllers.NewTaskStatusController(repo)

	nonExistentID := 99999

	// Act
	err := controller.DeleteByID(nonExistentID)

	// Assert
	// GORM не возвращает ошибку при удалении несуществующей записи
	// Это нормальное поведение, поэтому тест просто проверяет, что нет паники
	assert.NoError(t, err)
}

// TestTaskStatusController_FullFlow_Integration тестирует полный цикл работы со статусами
func TestTaskStatusController_FullFlow_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	repo := repositories.NewTaskStatusRepository(db)
	controller := controllers.NewTaskStatusController(repo)

	statusName := "test_full_flow"

	// Act & Assert - полный цикл операций

	// 1. Создание
	createdStatus, err := controller.Create(statusName)
	require.NoError(t, err)
	assert.NotNil(t, createdStatus)
	assert.Equal(t, statusName, createdStatus.Name)

	// 2. Получение по ID
	retrievedStatus, err := controller.GetByID(createdStatus.ID)
	require.NoError(t, err)
	assert.Equal(t, createdStatus.ID, retrievedStatus.ID)
	assert.Equal(t, statusName, retrievedStatus.Name)

	// 3. Получение всех статусов (должен включать созданный)
	allStatuses, err := controller.GetAll()
	require.NoError(t, err)
	found := false
	for _, status := range allStatuses {
		if status.ID == createdStatus.ID {
			found = true
			assert.Equal(t, statusName, status.Name)
			break
		}
	}
	assert.True(t, found, "Created status should be in GetAll result")

	// 4. Удаление
	err = controller.DeleteByID(createdStatus.ID)
	require.NoError(t, err)

	// 5. Проверка, что статус удален
	_, err = controller.GetByID(createdStatus.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "task status")
}
