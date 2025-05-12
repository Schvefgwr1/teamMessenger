package custom_errors

import (
	"errors"
	"fmt"
)

var (
	ErrStatusAlreadyExists = errors.New("task status with this name already exists")
)

// ============ File ============

type GetFileHTTPError struct {
	httpError string
	FileID    int
}

func (e *GetFileHTTPError) Error() string {
	return fmt.Sprintf("can't get file with id: %d, error: %s", e.FileID, e.httpError)
}

func NewGetFileHTTPError(fileID int, httpError string) *GetFileHTTPError {
	return &GetFileHTTPError{FileID: fileID, httpError: httpError}
}

// ============ User ============

type GetUserHTTPError struct {
	httpError string
	UserID    string
}

func (e *GetUserHTTPError) Error() string {
	return fmt.Sprintf("can't get user with id: %s, error: %s", e.UserID, e.httpError)
}

func NewGetUserHTTPError(userID string, httpError string) *GetUserHTTPError {
	return &GetUserHTTPError{UserID: userID, httpError: httpError}
}

// ============ Chat ============

type GetChatHTTPError struct {
	httpError string
	ChatID    string
}

func (e *GetChatHTTPError) Error() string {
	return fmt.Sprintf("can't get chat with id: %s, error: %s", e.ChatID, e.httpError)
}

func NewGetChatHTTPError(chatID string, httpError string) *GetChatHTTPError {
	return &GetChatHTTPError{ChatID: chatID, httpError: httpError}
}

// ============ Task Status ============

type TaskStatusNotFoundError struct {
	StatusName string
}

func (e *TaskStatusNotFoundError) Error() string {
	return fmt.Sprintf("task status with ID %s not found", e.StatusName)
}

func NewTaskStatusNotFoundError(statusName string) error {
	return &TaskStatusNotFoundError{StatusName: statusName}
}

// ============ Task Status ============

type TaskNotFoundError struct {
	TaskID int
}

func (e *TaskNotFoundError) Error() string {
	return fmt.Sprintf("task with id %d not found", e.TaskID)
}

func NewTaskNotFoundError(taskID int) error {
	return &TaskNotFoundError{TaskID: taskID}
}
