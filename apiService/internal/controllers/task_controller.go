package controllers

import (
	"apiService/internal/dto"
	"apiService/internal/http_clients"
	at "common/contracts/api-task"
	"errors"
	"github.com/google/uuid"
	"log"
)

type TaskController struct {
	taskClient http_clients.TaskClient
	fileClient http_clients.FileClient
}

func NewTaskController(taskClient http_clients.TaskClient, fileClient http_clients.FileClient) *TaskController {
	return &TaskController{
		taskClient: taskClient,
		fileClient: fileClient,
	}
}

func (ctrl *TaskController) CreateTask(req *dto.CreateTaskRequestGateway, creatorID uuid.UUID) (*at.TaskResponse, error) {
	// Парсим UUID из строк
	executorID, chatID, err := req.ParseUUIDs()
	if err != nil {
		return nil, err
	}

	// Загружаем файлы, если есть
	var fileIDs []int
	if len(req.Files) > 0 {
		fileIDs = make([]int, 0, len(req.Files))
		for _, file := range req.Files {
			uploadedFile, err := ctrl.fileClient.UploadFile(file)
			if err != nil {
				log.Printf("failed to upload file %s: %v\n", file.Filename, err)
				continue
			}
			if uploadedFile.ID != nil {
				fileIDs = append(fileIDs, *uploadedFile.ID)
			}
		}
	}

	// Создаем запрос к Task Service
	createReq := &at.CreateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
		CreatorID:   creatorID,
		FileIDs:     fileIDs,
	}

	// Устанавливаем ExecutorID (используем uuid.Nil если не указан)
	if executorID != nil {
		createReq.ExecutorID = *executorID
	} else {
		createReq.ExecutorID = uuid.Nil
	}

	// Устанавливаем ChatID (используем uuid.Nil если не указан)
	if chatID != nil {
		createReq.ChatID = *chatID
	} else {
		createReq.ChatID = uuid.Nil
	}

	taskResp, err := ctrl.taskClient.CreateTask(createReq)
	if err != nil {
		return nil, errors.New("error of task client")
	}

	return taskResp, nil
}

func (ctrl *TaskController) UpdateTaskStatus(taskID, statusID int) error {
	return ctrl.taskClient.UpdateTaskStatus(taskID, statusID)
}

func (ctrl *TaskController) GetTaskByID(taskID int) (*at.TaskServiceResponse, error) {
	return ctrl.taskClient.GetTaskByID(taskID)
}

func (ctrl *TaskController) GetUserTasks(userID string, limit, offset int) (*[]at.TaskToList, error) {
	return ctrl.taskClient.GetUserTasks(userID, limit, offset)
}
