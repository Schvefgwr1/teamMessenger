package controllers

import (
	fc "common/contracts/file-contracts"
	"common/http_clients"
	"github.com/google/uuid"
	"strconv"
	customErrors "taskService/internal/custom_errors"
	"taskService/internal/handlers/dto"
	"taskService/internal/models"
	"taskService/internal/repositories"
)

type TaskController struct {
	TaskRepo       repositories.TaskRepository
	TaskStatusRepo repositories.TaskStatusRepository
	TaskFileRepo   repositories.TaskFileRepository
}

func NewTaskController(
	taskRepo repositories.TaskRepository,
	taskStatusRepo repositories.TaskStatusRepository,
	taskFileRepo repositories.TaskFileRepository,
) *TaskController {
	return &TaskController{
		TaskRepo:       taskRepo,
		TaskStatusRepo: taskStatusRepo,
		TaskFileRepo:   taskFileRepo,
	}
}

func (c *TaskController) Create(taskDTO *dto.CreateTaskDTO) (*models.Task, error) {
	status, err := c.TaskStatusRepo.GetByName("created")
	if err != nil {
		return nil, customErrors.NewTaskStatusNotFoundError("created")
	}

	if _, errUser := http_clients.GetUserByID(&taskDTO.CreatorID); errUser != nil {
		return nil, customErrors.NewGetUserHTTPError(taskDTO.CreatorID.String(), err.Error())
	}

	if taskDTO.ExecutorID != uuid.Nil {
		if _, errTask := http_clients.GetUserByID(&taskDTO.ExecutorID); errTask != nil {
			return nil, customErrors.NewGetUserHTTPError(taskDTO.ExecutorID.String(), err.Error())
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
