package controllers

import (
	cc "common/contracts/chat-contracts"
	fc "common/contracts/file-contracts"
	cuc "common/contracts/user-contracts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"taskService/internal/handlers/dto"
	"taskService/internal/models"
)

// MockTaskRepository - мок для TaskRepository
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) UpdateStatus(taskID int, statusID int) error {
	args := m.Called(taskID, statusID)
	return args.Error(0)
}

func (m *MockTaskRepository) GetByID(taskID int) (*models.Task, error) {
	args := m.Called(taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) GetUserTasks(userID string, limit, offset int) (*[]dto.TaskToList, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]dto.TaskToList), args.Error(1)
}

// MockTaskStatusRepository - мок для TaskStatusRepository
type MockTaskStatusRepository struct {
	mock.Mock
}

func (m *MockTaskStatusRepository) GetByID(id int) (*models.TaskStatus, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskStatus), args.Error(1)
}

func (m *MockTaskStatusRepository) GetByName(name string) (*models.TaskStatus, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskStatus), args.Error(1)
}

func (m *MockTaskStatusRepository) Create(name string) (*models.TaskStatus, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskStatus), args.Error(1)
}

func (m *MockTaskStatusRepository) DeleteByID(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTaskStatusRepository) GetAll() ([]models.TaskStatus, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.TaskStatus), args.Error(1)
}

// MockTaskFileRepository - мок для TaskFileRepository
type MockTaskFileRepository struct {
	mock.Mock
}

func (m *MockTaskFileRepository) BulkCreate(taskFiles []models.TaskFile) error {
	args := m.Called(taskFiles)
	return args.Error(0)
}

// MockNotificationService - мок для NotificationServiceInterface
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendTaskCreatedNotification(
	taskID int,
	taskTitle string,
	creatorName string,
	executorID uuid.UUID,
	executorEmail string,
) error {
	args := m.Called(taskID, taskTitle, creatorName, executorID, executorEmail)
	return args.Error(0)
}

func (m *MockNotificationService) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Вспомогательные функции для создания тестовых данных
func createTestTask() *models.Task {
	taskID := 1
	statusID := 1
	creatorID := uuid.New()
	executorID := uuid.New()
	return &models.Task{
		ID:          taskID,
		Title:       "Test Task",
		Description: "Test Description",
		StatusID:    statusID,
		CreatorID:   creatorID,
		ExecutorID:  executorID,
		Status: &models.TaskStatus{
			ID:   statusID,
			Name: "created",
		},
		Files: []models.TaskFile{},
	}
}

func createTestTaskStatus() *models.TaskStatus {
	return &models.TaskStatus{
		ID:   1,
		Name: "created",
	}
}

func createTestTaskStatusWithID(id int, name string) *models.TaskStatus {
	return &models.TaskStatus{
		ID:   id,
		Name: name,
	}
}

func createTestTaskToList() *[]dto.TaskToList {
	return &[]dto.TaskToList{
		{
			ID:     1,
			Title:  "Test Task 1",
			Status: "created",
		},
		{
			ID:     2,
			Title:  "Test Task 2",
			Status: "in_progress",
		},
	}
}

func createTestFile() *fc.File {
	return &fc.File{
		ID:   1,
		Name: "test.txt",
		URL:  "http://example.com/test.txt",
	}
}

func stringPtr(s string) *string {
	return &s
}

// MockUserClient - мок для UserClientInterface
type MockUserClient struct {
	mock.Mock
}

func (m *MockUserClient) GetUserByID(userID *uuid.UUID) (*cuc.Response, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cuc.Response), args.Error(1)
}

// MockChatClient - мок для ChatClientInterface
type MockChatClient struct {
	mock.Mock
}

func (m *MockChatClient) GetChatByID(chatID string) (*cc.Chat, error) {
	args := m.Called(chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cc.Chat), args.Error(1)
}

// MockFileClient - мок для FileClientInterface
type MockFileClient struct {
	mock.Mock
}

func (m *MockFileClient) GetFileByID(fileID int) (*fc.File, error) {
	args := m.Called(fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fc.File), args.Error(1)
}

// Вспомогательные функции для создания тестовых данных HTTP клиентов
func createTestUserResponse() *cuc.Response {
	userID := uuid.New()
	username := "testuser"
	email := "test@example.com"
	return &cuc.Response{
		User: &cuc.User{
			ID:       userID,
			Username: username,
			Email:    email,
		},
		File:  nil,
		Error: nil,
	}
}

func createTestUserResponseWithEmail(email string) *cuc.Response {
	user := createTestUserResponse()
	user.User.Email = email
	return user
}

func createTestChat() *cc.Chat {
	chatID := uuid.New()
	return &cc.Chat{
		ID:   chatID,
		Name: "Test Chat",
	}
}
