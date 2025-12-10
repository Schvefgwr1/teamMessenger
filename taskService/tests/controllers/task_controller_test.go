package controllers

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"taskService/internal/controllers"
	"taskService/internal/custom_errors"
	"taskService/internal/handlers/dto"
	"taskService/internal/models"
)

// Тесты для TaskController.Create

func TestTaskController_Create_Success(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	creatorID := uuid.New()
	executorID := uuid.New()
	status := createTestTaskStatus()
	creator := createTestUserResponse()
	executor := createTestUserResponseWithEmail("executor@example.com")

	taskDTO := &dto.CreateTaskDTO{
		Title:       "Test Task",
		Description: stringPtr("Test Description"),
		CreatorID:   creatorID,
		ExecutorID:  executorID,
		FileIDs:     []int{1, 2},
	}

	expectedTask := createTestTask()
	expectedTask.Title = taskDTO.Title
	expectedTask.Description = *taskDTO.Description
	expectedTask.CreatorID = creatorID
	expectedTask.ExecutorID = executorID

	mockTaskStatusRepo.On("GetByName", "created").Return(status, nil)
	mockUserClient.On("GetUserByID", &creatorID).Return(creator, nil)
	mockUserClient.On("GetUserByID", &executorID).Return(executor, nil)
	mockFileClient.On("GetFileByID", 1).Return(createTestFile(), nil)
	mockFileClient.On("GetFileByID", 2).Return(createTestFile(), nil)
	mockTaskRepo.On("Create", mock.MatchedBy(func(task *models.Task) bool {
		return task.Title == taskDTO.Title && task.CreatorID == creatorID
	})).Return(nil).Run(func(args mock.Arguments) {
		task := args.Get(0).(*models.Task)
		task.ID = expectedTask.ID
	})
	mockTaskFileRepo.On("BulkCreate", mock.AnythingOfType("[]models.TaskFile")).Return(nil)
	mockNotificationService.On("SendTaskCreatedNotification",
		expectedTask.ID,
		taskDTO.Title,
		creator.User.Username,
		executorID,
		executor.User.Email,
	).Return(nil)

	// Act
	result, err := controller.Create(taskDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, taskDTO.Title, result.Title)
	assert.Equal(t, *taskDTO.Description, result.Description)
	assert.Equal(t, creatorID, result.CreatorID)
	assert.Equal(t, executorID, result.ExecutorID)

	mockTaskStatusRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
	mockTaskFileRepo.AssertExpectations(t)
	mockNotificationService.AssertExpectations(t)
}

func TestTaskController_Create_WithoutExecutor(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	creatorID := uuid.New()
	status := createTestTaskStatus()
	creator := createTestUserResponse()

	taskDTO := &dto.CreateTaskDTO{
		Title:       "Test Task",
		Description: stringPtr("Test Description"),
		CreatorID:   creatorID,
		ExecutorID:  uuid.Nil,
		FileIDs:     []int{},
	}

	expectedTask := createTestTask()
	expectedTask.Title = taskDTO.Title
	expectedTask.Description = *taskDTO.Description
	expectedTask.CreatorID = creatorID
	expectedTask.ExecutorID = uuid.Nil

	mockTaskStatusRepo.On("GetByName", "created").Return(status, nil)
	mockUserClient.On("GetUserByID", &creatorID).Return(creator, nil)
	mockTaskRepo.On("Create", mock.MatchedBy(func(task *models.Task) bool {
		return task.Title == taskDTO.Title && task.CreatorID == creatorID
	})).Return(nil).Run(func(args mock.Arguments) {
		task := args.Get(0).(*models.Task)
		task.ID = expectedTask.ID
	})

	// Act
	result, err := controller.Create(taskDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uuid.Nil, result.ExecutorID)

	mockTaskStatusRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
	mockNotificationService.AssertNotCalled(t, "SendTaskCreatedNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestTaskController_Create_StatusNotFound(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	creatorID := uuid.New()
	taskDTO := &dto.CreateTaskDTO{
		Title:      "Test Task",
		CreatorID:  creatorID,
		ExecutorID: uuid.Nil,
	}

	mockTaskStatusRepo.On("GetByName", "created").Return(nil, gorm.ErrRecordNotFound)

	// Act
	result, err := controller.Create(taskDTO)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	var statusNotFoundErr *custom_errors.TaskStatusNotFoundError
	assert.True(t, errors.As(err, &statusNotFoundErr))

	mockTaskStatusRepo.AssertExpectations(t)
	mockUserClient.AssertNotCalled(t, "GetUserByID", mock.Anything)
}

func TestTaskController_Create_GetCreatorError(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	creatorID := uuid.New()
	status := createTestTaskStatus()
	userError := errors.New("user not found")

	taskDTO := &dto.CreateTaskDTO{
		Title:      "Test Task",
		CreatorID:  creatorID,
		ExecutorID: uuid.Nil,
	}

	mockTaskStatusRepo.On("GetByName", "created").Return(status, nil)
	mockUserClient.On("GetUserByID", &creatorID).Return(nil, userError)

	// Act
	result, err := controller.Create(taskDTO)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	var getUserErr *custom_errors.GetUserHTTPError
	assert.True(t, errors.As(err, &getUserErr))

	mockTaskStatusRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockTaskRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestTaskController_Create_GetExecutorError(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	creatorID := uuid.New()
	executorID := uuid.New()
	status := createTestTaskStatus()
	creator := createTestUserResponse()
	userError := errors.New("user not found")

	taskDTO := &dto.CreateTaskDTO{
		Title:      "Test Task",
		CreatorID:  creatorID,
		ExecutorID: executorID,
	}

	mockTaskStatusRepo.On("GetByName", "created").Return(status, nil)
	mockUserClient.On("GetUserByID", &creatorID).Return(creator, nil)
	mockUserClient.On("GetUserByID", &executorID).Return(nil, userError)

	// Act
	result, err := controller.Create(taskDTO)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	var getUserErr *custom_errors.GetUserHTTPError
	assert.True(t, errors.As(err, &getUserErr))

	mockTaskStatusRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockTaskRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestTaskController_Create_GetChatError(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	creatorID := uuid.New()
	chatID := uuid.New()
	status := createTestTaskStatus()
	creator := createTestUserResponse()
	chatError := errors.New("chat not found")

	taskDTO := &dto.CreateTaskDTO{
		Title:      "Test Task",
		CreatorID:  creatorID,
		ExecutorID: uuid.Nil,
		ChatID:     chatID,
	}

	mockTaskStatusRepo.On("GetByName", "created").Return(status, nil)
	mockUserClient.On("GetUserByID", &creatorID).Return(creator, nil)
	mockChatClient.On("GetChatByID", chatID.String()).Return(nil, chatError)

	// Act
	result, err := controller.Create(taskDTO)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	var getChatErr *custom_errors.GetChatHTTPError
	assert.True(t, errors.As(err, &getChatErr))

	mockTaskStatusRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockChatClient.AssertExpectations(t)
	mockTaskRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestTaskController_Create_GetFileError(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	creatorID := uuid.New()
	status := createTestTaskStatus()
	creator := createTestUserResponse()
	fileError := errors.New("file not found")

	taskDTO := &dto.CreateTaskDTO{
		Title:      "Test Task",
		CreatorID:  creatorID,
		ExecutorID: uuid.Nil,
		FileIDs:    []int{1},
	}

	mockTaskStatusRepo.On("GetByName", "created").Return(status, nil)
	mockUserClient.On("GetUserByID", &creatorID).Return(creator, nil)
	mockFileClient.On("GetFileByID", 1).Return(nil, fileError)

	// Act
	result, err := controller.Create(taskDTO)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	var getFileErr *custom_errors.GetFileHTTPError
	assert.True(t, errors.As(err, &getFileErr))

	mockTaskStatusRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
	mockTaskRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestTaskController_Create_TaskRepoError(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	creatorID := uuid.New()
	status := createTestTaskStatus()
	creator := createTestUserResponse()
	repoError := errors.New("database error")

	taskDTO := &dto.CreateTaskDTO{
		Title:      "Test Task",
		CreatorID:  creatorID,
		ExecutorID: uuid.Nil,
		FileIDs:    []int{},
	}

	mockTaskStatusRepo.On("GetByName", "created").Return(status, nil)
	mockUserClient.On("GetUserByID", &creatorID).Return(creator, nil)
	mockTaskRepo.On("Create", mock.Anything).Return(repoError)

	// Act
	result, err := controller.Create(taskDTO)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoError, err)

	mockTaskStatusRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskController_Create_TaskFileRepoError(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	creatorID := uuid.New()
	status := createTestTaskStatus()
	creator := createTestUserResponse()
	fileRepoError := errors.New("database error")

	taskDTO := &dto.CreateTaskDTO{
		Title:      "Test Task",
		CreatorID:  creatorID,
		ExecutorID: uuid.Nil,
		FileIDs:    []int{1},
	}

	mockTaskStatusRepo.On("GetByName", "created").Return(status, nil)
	mockUserClient.On("GetUserByID", &creatorID).Return(creator, nil)
	mockFileClient.On("GetFileByID", 1).Return(createTestFile(), nil)
	mockTaskRepo.On("Create", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		task := args.Get(0).(*models.Task)
		task.ID = 1
	})
	mockTaskFileRepo.On("BulkCreate", mock.Anything).Return(fileRepoError)

	// Act
	result, err := controller.Create(taskDTO)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, fileRepoError, err)

	mockTaskStatusRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
	mockTaskFileRepo.AssertExpectations(t)
}

func TestTaskController_Create_NotificationError(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	creatorID := uuid.New()
	executorID := uuid.New()
	status := createTestTaskStatus()
	creator := createTestUserResponse()
	executor := createTestUserResponseWithEmail("executor@example.com")
	notificationError := errors.New("kafka error")

	taskDTO := &dto.CreateTaskDTO{
		Title:      "Test Task",
		CreatorID:  creatorID,
		ExecutorID: executorID,
		FileIDs:    []int{},
	}

	mockTaskStatusRepo.On("GetByName", "created").Return(status, nil)
	mockUserClient.On("GetUserByID", &creatorID).Return(creator, nil)
	mockUserClient.On("GetUserByID", &executorID).Return(executor, nil)
	mockTaskRepo.On("Create", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		task := args.Get(0).(*models.Task)
		task.ID = 1
	})
	mockNotificationService.On("SendTaskCreatedNotification",
		mock.Anything,
		taskDTO.Title,
		creator.User.Username,
		executorID,
		executor.User.Email,
	).Return(notificationError)

	// Act
	result, err := controller.Create(taskDTO)

	// Assert
	// Ошибка уведомления не должна прерывать создание задачи
	require.NoError(t, err)
	assert.NotNil(t, result)

	mockTaskStatusRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
	mockNotificationService.AssertExpectations(t)
}

func TestTaskController_Create_WithChatID(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	creatorID := uuid.New()
	chatID := uuid.New()
	status := createTestTaskStatus()
	creator := createTestUserResponse()
	chat := createTestChat()

	taskDTO := &dto.CreateTaskDTO{
		Title:      "Test Task",
		CreatorID:  creatorID,
		ExecutorID: uuid.Nil,
		ChatID:     chatID,
		FileIDs:    []int{},
	}

	mockTaskStatusRepo.On("GetByName", "created").Return(status, nil)
	mockUserClient.On("GetUserByID", &creatorID).Return(creator, nil)
	mockChatClient.On("GetChatByID", chatID.String()).Return(chat, nil)
	mockTaskRepo.On("Create", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		task := args.Get(0).(*models.Task)
		task.ID = 1
	})

	// Act
	result, err := controller.Create(taskDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, chatID, result.ChatID)

	mockTaskStatusRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockChatClient.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskController_Create_NilDescription(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	creatorID := uuid.New()
	status := createTestTaskStatus()
	creator := createTestUserResponse()

	taskDTO := &dto.CreateTaskDTO{
		Title:       "Test Task",
		Description: nil,
		CreatorID:   creatorID,
		ExecutorID:  uuid.Nil,
		FileIDs:     []int{},
	}

	mockTaskStatusRepo.On("GetByName", "created").Return(status, nil)
	mockUserClient.On("GetUserByID", &creatorID).Return(creator, nil)
	mockTaskRepo.On("Create", mock.MatchedBy(func(task *models.Task) bool {
		return task.Description == ""
	})).Return(nil).Run(func(args mock.Arguments) {
		task := args.Get(0).(*models.Task)
		task.ID = 1
	})

	// Act
	result, err := controller.Create(taskDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "", result.Description)

	mockTaskStatusRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

// Тесты для TaskController.UpdateStatus

func TestTaskController_UpdateStatus_Success(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	taskID := 1
	statusID := 2
	status := createTestTaskStatusWithID(statusID, "in_progress")

	mockTaskStatusRepo.On("GetByID", statusID).Return(status, nil)
	mockTaskRepo.On("UpdateStatus", taskID, statusID).Return(nil)

	// Act
	err := controller.UpdateStatus(taskID, statusID)

	// Assert
	require.NoError(t, err)

	mockTaskStatusRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskController_UpdateStatus_StatusNotFound(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	taskID := 1
	statusID := 999

	mockTaskStatusRepo.On("GetByID", statusID).Return(nil, gorm.ErrRecordNotFound)

	// Act
	err := controller.UpdateStatus(taskID, statusID)

	// Assert
	require.Error(t, err)
	var statusNotFoundErr *custom_errors.TaskStatusNotFoundError
	assert.True(t, errors.As(err, &statusNotFoundErr))

	mockTaskStatusRepo.AssertExpectations(t)
	mockTaskRepo.AssertNotCalled(t, "UpdateStatus", mock.Anything, mock.Anything)
}

func TestTaskController_UpdateStatus_TaskRepoError(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	taskID := 1
	statusID := 2
	status := createTestTaskStatusWithID(statusID, "in_progress")
	repoError := errors.New("database error")

	mockTaskStatusRepo.On("GetByID", statusID).Return(status, nil)
	mockTaskRepo.On("UpdateStatus", taskID, statusID).Return(repoError)

	// Act
	err := controller.UpdateStatus(taskID, statusID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, repoError, err)

	mockTaskStatusRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

// Тесты для TaskController.GetByID

func TestTaskController_GetByID_Success(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	taskID := 1
	task := createTestTask()
	task.Files = []models.TaskFile{
		{TaskID: taskID, FileID: 1},
		{TaskID: taskID, FileID: 2},
	}
	file1 := createTestFile()
	file1.ID = 1
	file2 := createTestFile()
	file2.ID = 2

	mockTaskRepo.On("GetByID", taskID).Return(task, nil)
	mockFileClient.On("GetFileByID", 1).Return(file1, nil)
	mockFileClient.On("GetFileByID", 2).Return(file2, nil)

	// Act
	result, err := controller.GetByID(taskID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, task, result.Task)
	assert.NotNil(t, result.Files)
	assert.Len(t, *result.Files, 2)

	mockTaskRepo.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
}

func TestTaskController_GetByID_TaskNotFound(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	taskID := 999
	mockTaskRepo.On("GetByID", taskID).Return(nil, gorm.ErrRecordNotFound)

	// Act
	result, err := controller.GetByID(taskID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)

	mockTaskRepo.AssertExpectations(t)
	mockFileClient.AssertNotCalled(t, "GetFileByID", mock.Anything)
}

func TestTaskController_GetByID_GetFileError(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	taskID := 1
	task := createTestTask()
	task.Files = []models.TaskFile{
		{TaskID: taskID, FileID: 1},
	}
	fileError := errors.New("file not found")

	mockTaskRepo.On("GetByID", taskID).Return(task, nil)
	mockFileClient.On("GetFileByID", 1).Return(nil, fileError)

	// Act
	result, err := controller.GetByID(taskID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	var getFileErr *custom_errors.GetFileHTTPError
	assert.True(t, errors.As(err, &getFileErr))

	mockTaskRepo.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
}

func TestTaskController_GetByID_NoFiles(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	taskID := 1
	task := createTestTask()
	task.Files = []models.TaskFile{}

	mockTaskRepo.On("GetByID", taskID).Return(task, nil)

	// Act
	result, err := controller.GetByID(taskID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, task, result.Task)
	assert.NotNil(t, result.Files)
	assert.Len(t, *result.Files, 0)

	mockTaskRepo.AssertExpectations(t)
	mockFileClient.AssertNotCalled(t, "GetFileByID", mock.Anything)
}

// Тесты для TaskController.GetUserTasks

func TestTaskController_GetUserTasks_Success(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	userID := uuid.New().String()
	limit := 10
	offset := 0
	expectedTasks := createTestTaskToList()

	mockTaskRepo.On("GetUserTasks", userID, limit, offset).Return(expectedTasks, nil)

	// Act
	result, err := controller.GetUserTasks(userID, limit, offset)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedTasks, result)

	mockTaskRepo.AssertExpectations(t)
}

func TestTaskController_GetUserTasks_RepositoryError(t *testing.T) {
	// Arrange
	mockTaskRepo := new(MockTaskRepository)
	mockTaskStatusRepo := new(MockTaskStatusRepository)
	mockTaskFileRepo := new(MockTaskFileRepository)
	mockNotificationService := new(MockNotificationService)
	mockUserClient := new(MockUserClient)
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)

	controller := controllers.NewTaskControllerWithClients(
		mockTaskRepo,
		mockTaskStatusRepo,
		mockTaskFileRepo,
		mockNotificationService,
		mockUserClient,
		mockChatClient,
		mockFileClient,
	)

	userID := uuid.New().String()
	limit := 10
	offset := 0
	repoError := errors.New("database error")

	mockTaskRepo.On("GetUserTasks", userID, limit, offset).Return(nil, repoError)

	// Act
	result, err := controller.GetUserTasks(userID, limit, offset)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoError, err)

	mockTaskRepo.AssertExpectations(t)
}
