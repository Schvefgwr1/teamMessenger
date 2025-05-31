package controllers

import (
	fc "common/contracts/file-contracts"
	"common/http_clients"
	"github.com/google/uuid"
	"log"
	"strconv"
	customErrors "taskService/internal/custom_errors"
	"taskService/internal/handlers/dto"
	"taskService/internal/models"
	"taskService/internal/repositories"
	"taskService/internal/services"
)

type TaskController struct {
	TaskRepo            repositories.TaskRepository
	TaskStatusRepo      repositories.TaskStatusRepository
	TaskFileRepo        repositories.TaskFileRepository
	NotificationService *services.NotificationService
}

func NewTaskController(
	taskRepo repositories.TaskRepository,
	taskStatusRepo repositories.TaskStatusRepository,
	taskFileRepo repositories.TaskFileRepository,
	notificationService *services.NotificationService,
) *TaskController {
	return &TaskController{
		TaskRepo:            taskRepo,
		TaskStatusRepo:      taskStatusRepo,
		TaskFileRepo:        taskFileRepo,
		NotificationService: notificationService,
	}
}

func (c *TaskController) Create(taskDTO *dto.CreateTaskDTO) (*models.Task, error) {
	status, err := c.TaskStatusRepo.GetByName("created")
	if err != nil {
		return nil, customErrors.NewTaskStatusNotFoundError("created")
	}

	// Получаем информацию о создателе
	creator, errUser := http_clients.GetUserByID(&taskDTO.CreatorID)
	if errUser != nil {
		return nil, customErrors.NewGetUserHTTPError(taskDTO.CreatorID.String(), err.Error())
	}

	var executorEmail string
	if taskDTO.ExecutorID != uuid.Nil {
		executor, errTask := http_clients.GetUserByID(&taskDTO.ExecutorID)
		if errTask != nil {
			return nil, customErrors.NewGetUserHTTPError(taskDTO.ExecutorID.String(), err.Error())
		}

		if executor.User.Email != "" {
			executorEmail = executor.User.Email
		}
	}

	if taskDTO.ChatID != uuid.Nil {
		if _, errChat := http_clients.GetChatByID(taskDTO.ChatID.String()); errChat != nil {
			return nil, customErrors.NewGetChatHTTPError(taskDTO.ChatID.String(), err.Error())
		}
	}

	var taskFiles []models.TaskFile
	for _, fileID := range taskDTO.FileIDs {
		if _, errFile := http_clients.GetFileByID(fileID); errFile != nil {
			return nil, customErrors.NewGetFileHTTPError(fileID, err.Error())
		}
	}

	task := &models.Task{
		Title:       taskDTO.Title,
		Description: taskDTO.Description,
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
		file, err := http_clients.GetFileByID(taskFile.FileID)
		if err != nil {
			return nil, customErrors.NewGetFileHTTPError(taskFile.FileID, err.Error())
		}
		files = append(files, *file)
	}
	return &dto.TaskResponse{Task: task, Files: &files}, nil
}

func (c *TaskController) GetUserTasks(userID string, limit, offset int) (*[]dto.TaskToList, error) {
	return c.TaskRepo.GetUserTasks(userID, limit, offset)
}
