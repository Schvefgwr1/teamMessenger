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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

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

	// Инвалидация кеша задач для создателя
	_ = ctrl.cacheService.DeleteUserTasksCache(ctx, creatorID.String())

	// Инвалидация кеша задач для исполнителя (если указан)
	if executorID != nil && *executorID != uuid.Nil {
		_ = ctrl.cacheService.DeleteUserTasksCache(ctx, executorID.String())
	}

	return taskResp, nil
}

func (ctrl *TaskController) UpdateTaskStatus(taskID, statusID int) error {
	err := ctrl.taskClient.UpdateTaskStatus(taskID, statusID)
	if err != nil {
		return err
	}

	// Инвалидация кеша задачи
	ctx := context.Background()
	_ = ctrl.cacheService.DeleteTaskCache(ctx, taskID)

	return nil
}

func (ctrl *TaskController) GetTaskByID(taskID int) (*at.TaskServiceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Пытаемся получить из кеша
	var cachedTask at.TaskServiceResponse
	err := ctrl.cacheService.GetTaskCache(ctx, taskID, &cachedTask)
	if err == nil {
		log.Printf("Task %d found in cache", taskID)
		return &cachedTask, nil
	}

	// Получаем из сервиса
	task, err := ctrl.taskClient.GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кеш
	if err := ctrl.cacheService.SetTaskCache(ctx, taskID, task); err != nil {
		log.Printf("Failed to cache task %d: %v", taskID, err)
	}

	return task, nil
}

func (ctrl *TaskController) GetUserTasks(userID string, limit, offset int) (*[]at.TaskToList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Кешируем только первую страницу (offset = 0)
	if offset == 0 && limit <= 20 {
		var cachedTasks []at.TaskToList
		err := ctrl.cacheService.GetUserTasksCache(ctx, userID, &cachedTasks)
		if err == nil {
			log.Printf("User tasks for %s found in cache", userID)
			if limit < len(cachedTasks) {
				result := cachedTasks[:limit]
				return &result, nil
			}
			return &cachedTasks, nil
		}

		// Получаем из сервиса (всегда запрашиваем 20 для кеша)
		tasks, err := ctrl.taskClient.GetUserTasks(userID, 20, 0)
		if err != nil {
			return nil, err
		}

		// Сохраняем в кеш
		if tasks != nil {
			if err := ctrl.cacheService.SetUserTasksCache(ctx, userID, *tasks); err != nil {
				log.Printf("Failed to cache user tasks for %s: %v", userID, err)
			}
		}

		// Возвращаем запрошенное количество
		if tasks != nil && limit < len(*tasks) {
			result := (*tasks)[:limit]
			return &result, nil
		}
		return tasks, nil
	}

	// Для остальных запросов идём напрямую в сервис
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
