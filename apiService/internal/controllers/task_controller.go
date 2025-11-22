package controllers

import (
	"apiService/internal/dto"
	"apiService/internal/http_clients"
	"apiService/internal/services"
	at "common/contracts/api-task"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

type TaskController struct {
	taskClient   http_clients.TaskClient
	fileClient   http_clients.FileClient
	cacheService *services.CacheService
}

func NewTaskController(taskClient http_clients.TaskClient, fileClient http_clients.FileClient, cacheService *services.CacheService) *TaskController {
	return &TaskController{
		taskClient:   taskClient,
		fileClient:   fileClient,
		cacheService: cacheService,
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

// GetAllStatuses - получить все статусы задач с кешированием
func (ctrl *TaskController) GetAllStatuses() ([]at.TaskStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Пытаемся получить из кеша
	cacheKey := "task:statuses"
	var cachedStatuses []at.TaskStatus
	err := ctrl.cacheService.Get(ctx, cacheKey, &cachedStatuses)
	if err == nil {
		log.Printf("Task statuses found in cache")
		return cachedStatuses, nil
	}

	// Получаем из сервиса
	statuses, err := ctrl.taskClient.GetAllStatuses()
	if err != nil {
		return nil, err
	}

	// Сохраняем в кеш на 1 час
	if err := ctrl.cacheService.Set(ctx, cacheKey, statuses, time.Hour); err != nil {
		log.Printf("Failed to cache task statuses: %v", err)
	}

	return statuses, nil
}

// CreateStatus - создать новый статус задачи с инвалидацией кеша
func (ctrl *TaskController) CreateStatus(name string) (*at.TaskStatus, error) {
	req := &at.CreateStatusRequest{Name: name}
	status, err := ctrl.taskClient.CreateStatus(req)
	if err != nil {
		return nil, err
	}

	// Инвалидация кеша списка статусов
	ctx := context.Background()
	cacheKey := "task:statuses"
	_ = ctrl.cacheService.Delete(ctx, cacheKey)

	return status, nil
}

// GetStatusByID - получить статус по ID с кешированием
func (ctrl *TaskController) GetStatusByID(statusID int) (*at.TaskStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Пытаемся получить из кеша
	cacheKey := fmt.Sprintf("task:status:%d", statusID)
	var cachedStatus at.TaskStatus
	err := ctrl.cacheService.Get(ctx, cacheKey, &cachedStatus)
	if err == nil {
		log.Printf("Task status %d found in cache", statusID)
		return &cachedStatus, nil
	}

	// Получаем из сервиса
	status, err := ctrl.taskClient.GetStatusByID(statusID)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кеш на 1 час
	if err := ctrl.cacheService.Set(ctx, cacheKey, status, time.Hour); err != nil {
		log.Printf("Failed to cache task status %d: %v", statusID, err)
	}

	return status, nil
}

// DeleteStatus - удалить статус задачи с инвалидацией кеша
func (ctrl *TaskController) DeleteStatus(statusID int) error {
	err := ctrl.taskClient.DeleteStatus(statusID)
	if err != nil {
		return err
	}

	// Инвалидация кеша
	ctx := context.Background()
	statusKey := fmt.Sprintf("task:status:%d", statusID)
	_ = ctrl.cacheService.Delete(ctx, statusKey)

	listKey := "task:statuses"
	_ = ctrl.cacheService.Delete(ctx, listKey)

	return nil
}
