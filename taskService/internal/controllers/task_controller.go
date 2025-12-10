package controllers

import (
	fc "common/contracts/file-contracts"
	"log"
	"strconv"
	customErrors "taskService/internal/custom_errors"
	"taskService/internal/handlers/dto"
	"taskService/internal/http_clients"
	"taskService/internal/models"
	"taskService/internal/repositories"
	"taskService/internal/services"

	"github.com/google/uuid"
)

type TaskController struct {
	TaskRepo            repositories.TaskRepository
	TaskStatusRepo      repositories.TaskStatusRepository
	TaskFileRepo        repositories.TaskFileRepository
	NotificationService services.NotificationServiceInterface
	UserClient          http_clients.UserClientInterface
	ChatClient          http_clients.ChatClientInterface
	FileClient          http_clients.FileClientInterface
}

func NewTaskController(
	taskRepo repositories.TaskRepository,
	taskStatusRepo repositories.TaskStatusRepository,
	taskFileRepo repositories.TaskFileRepository,
	notificationService services.NotificationServiceInterface,
) *TaskController {
	return &TaskController{
		TaskRepo:            taskRepo,
		TaskStatusRepo:      taskStatusRepo,
		TaskFileRepo:        taskFileRepo,
		NotificationService: notificationService,
		UserClient:          http_clients.NewUserClientAdapter(),
		ChatClient:          http_clients.NewChatClientAdapter(),
		FileClient:          http_clients.NewFileClientAdapter(),
	}
}

// NewTaskControllerWithClients создает контроллер с указанными HTTP клиентами (для тестирования)
func NewTaskControllerWithClients(
	taskRepo repositories.TaskRepository,
	taskStatusRepo repositories.TaskStatusRepository,
	taskFileRepo repositories.TaskFileRepository,
	notificationService services.NotificationServiceInterface,
	userClient http_clients.UserClientInterface,
	chatClient http_clients.ChatClientInterface,
	fileClient http_clients.FileClientInterface,
) *TaskController {
	return &TaskController{
		TaskRepo:            taskRepo,
		TaskStatusRepo:      taskStatusRepo,
		TaskFileRepo:        taskFileRepo,
		NotificationService: notificationService,
		UserClient:          userClient,
		ChatClient:          chatClient,
		FileClient:          fileClient,
	}
}

func (c *TaskController) Create(taskDTO *dto.CreateTaskDTO) (*models.Task, error) {
	status, err := c.TaskStatusRepo.GetByName("created")
	if err != nil {
		return nil, customErrors.NewTaskStatusNotFoundError("created")
	}

	// Получаем информацию о создателе
	creator, errUser := c.UserClient.GetUserByID(&taskDTO.CreatorID)
	if errUser != nil {
		return nil, customErrors.NewGetUserHTTPError(taskDTO.CreatorID.String(), errUser.Error())
	}

	var executorEmail string
	if taskDTO.ExecutorID != uuid.Nil {
		executor, errTask := c.UserClient.GetUserByID(&taskDTO.ExecutorID)
		if errTask != nil {
			return nil, customErrors.NewGetUserHTTPError(taskDTO.ExecutorID.String(), errTask.Error())
		}

		if executor.User.Email != "" {
			executorEmail = executor.User.Email
		}
	}

	if taskDTO.ChatID != uuid.Nil {
		if _, errChat := c.ChatClient.GetChatByID(taskDTO.ChatID.String()); errChat != nil {
			return nil, customErrors.NewGetChatHTTPError(taskDTO.ChatID.String(), errChat.Error())
		}
	}

	var taskFiles []models.TaskFile
	for _, fileID := range taskDTO.FileIDs {
		if _, errFile := c.FileClient.GetFileByID(fileID); errFile != nil {
			return nil, customErrors.NewGetFileHTTPError(fileID, errFile.Error())
		}
	}

	var desc string
	if taskDTO.Description == nil {
		desc = ""
	} else {
		desc = *taskDTO.Description
	}

	task := &models.Task{
		Title:       taskDTO.Title,
		Description: desc,
		CreatorID:   taskDTO.CreatorID,
		ExecutorID:  taskDTO.ExecutorID,
		ChatID:      taskDTO.ChatID,
		Status:      status,
		StatusID:    status.ID,
	}

	if err := c.TaskRepo.Create(task); err != nil {
		return nil, err
	}

	for _, fileID := range taskDTO.FileIDs {
		taskFiles = append(taskFiles, models.TaskFile{
			TaskID: task.ID,
			FileID: fileID,
		})
	}
	if len(taskFiles) > 0 {
		if err := c.TaskFileRepo.BulkCreate(taskFiles); err != nil {
			return nil, err
		} else {
			task.Files = taskFiles
		}
	}

	// Отправляем уведомление о новой задаче, если есть исполнитель
	if taskDTO.ExecutorID != uuid.Nil && executorEmail != "" {
		creatorName := "Unknown user"
		if creator.User.Username != "" {
			creatorName = creator.User.Username
		}

		if err := c.NotificationService.SendTaskCreatedNotification(
			task.ID,
			task.Title,
			creatorName,
			taskDTO.ExecutorID,
			executorEmail,
		); err != nil {
			// Логируем ошибку, но не прерываем процесс создания задачи
			log.Printf("Failed to send task notification: %v", err)
		}
	}

	return task, nil
}

func (c *TaskController) UpdateStatus(taskID, statusID int) error {
	if _, err := c.TaskStatusRepo.GetByID(statusID); err != nil {
		return customErrors.NewTaskStatusNotFoundError(strconv.Itoa(statusID))
	}
	return c.TaskRepo.UpdateStatus(taskID, statusID)
}

func (c *TaskController) GetByID(taskID int) (*dto.TaskResponse, error) {
	task, err := c.TaskRepo.GetByID(taskID)
	if err != nil {
		return nil, err
	}
	var files []fc.File
	for _, taskFile := range task.Files {
		file, errFile := c.FileClient.GetFileByID(taskFile.FileID)
		if errFile != nil {
			return nil, customErrors.NewGetFileHTTPError(taskFile.FileID, errFile.Error())
		}
		files = append(files, *file)
	}
	return &dto.TaskResponse{Task: task, Files: &files}, nil
}

func (c *TaskController) GetUserTasks(userID string, limit, offset int) (*[]dto.TaskToList, error) {
	return c.TaskRepo.GetUserTasks(userID, limit, offset)
}
