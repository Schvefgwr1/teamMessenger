package api_task

import (
	fc "common/contracts/file-contracts"
	"github.com/google/uuid"
	"time"
)

// CreateTaskRequest - запрос на создание задачи (должен соответствовать CreateTaskDTO в taskService)
type CreateTaskRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	CreatorID   uuid.UUID `json:"creator_id" binding:"required"`
	ExecutorID  uuid.UUID `json:"executor_id"`
	ChatID      uuid.UUID `json:"chat_id"`
	FileIDs     []int     `json:"file_ids"`
}

// TaskResponse - ответ с задачей
type TaskResponse struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatorID   uuid.UUID  `json:"creatorID"`
	ExecutorID  *uuid.UUID `json:"executor_id,omitempty"`
	ChatID      *uuid.UUID `json:"chat_id,omitempty"`
	Status      TaskStatus `json:"status"`
	Files       []TaskFile `json:"files,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
}

type TaskFile struct {
	TaskID int `json:"taskID"`
	FileID int `json:"fileID"`
}

// TaskStatus - статус задачи
type TaskStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TaskToList - задача для списка
type TaskToList struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// TaskListResponse - ответ со списком задач
type TaskListResponse struct {
	Tasks []TaskToList `json:"tasks"`
	Total int64        `json:"total"`
}

// TaskServiceResponse - реальная структура ответа от taskService
type TaskServiceResponse struct {
	Task  *TaskResponse `json:"task"`
	Files *[]fc.File    `json:"files"`
}

// CreateStatusRequest - запрос на создание статуса задачи
type CreateStatusRequest struct {
	Name string `json:"name" binding:"required"`
}
