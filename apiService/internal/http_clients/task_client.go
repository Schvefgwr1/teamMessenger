package http_clients

import (
	"bytes"
	at "common/contracts/api-task"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TaskClient interface {
	CreateTask(req *at.CreateTaskRequest) (*at.TaskResponse, error)
	UpdateTaskStatus(taskID, statusID int) error
	GetTaskByID(taskID int) (*at.TaskServiceResponse, error)
	GetUserTasks(userID string, limit, offset int) (*[]at.TaskToList, error)
	GetAllStatuses() ([]at.TaskStatus, error)
	CreateStatus(req *at.CreateStatusRequest) (*at.TaskStatus, error)
	GetStatusByID(statusID int) (*at.TaskStatus, error)
	DeleteStatus(statusID int) error
}

type taskClient struct {
	host string
}

func NewTaskClient(host string) TaskClient {
	return &taskClient{host: host}
}

func (c *taskClient) CreateTask(req *at.CreateTaskRequest) (*at.TaskResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	resp, err := http.Post(c.host+"/api/v1/tasks", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("request to task service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("task service returned error: %s", string(bodyBytes))
	}

	var serviceTask at.TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&serviceTask); err != nil {
		return nil, fmt.Errorf("failed to decode task service response: %w", err)
	}

	return &serviceTask, nil
}

func (c *taskClient) UpdateTaskStatus(taskID, statusID int) error {
	url := fmt.Sprintf("%s/api/v1/tasks/%d/status/%d", c.host, taskID, statusID)

	httpReq, err := http.NewRequest(http.MethodPatch, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request to task service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("task service returned error: %s", string(bodyBytes))
	}

	return nil
}

func (c *taskClient) GetTaskByID(taskID int) (*at.TaskServiceResponse, error) {
	url := fmt.Sprintf("%s/api/v1/tasks/%d", c.host, taskID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("task service returned error: %s", string(bodyBytes))
	}

	var serviceResp at.TaskServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&serviceResp); err != nil {
		return nil, fmt.Errorf("failed to decode task response: %w", err)
	}

	return &serviceResp, nil
}

func (c *taskClient) GetUserTasks(userID string, limit, offset int) (*[]at.TaskToList, error) {
	url := fmt.Sprintf("%s/api/v1/users/%s/tasks?limit=%d&offset=%d", c.host, userID, limit, offset)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tasks: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("task service returned error: %s", string(bodyBytes))
	}

	var tasks []at.TaskToList
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, fmt.Errorf("failed to decode tasks response: %w", err)
	}

	return &tasks, nil
}

// GetAllStatuses - получить все статусы задач
func (c *taskClient) GetAllStatuses() ([]at.TaskStatus, error) {
	url := fmt.Sprintf("%s/api/v1/tasks/statuses", c.host)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get statuses: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("task service returned error: %s", string(bodyBytes))
	}

	var statuses []at.TaskStatus
	if err := json.NewDecoder(resp.Body).Decode(&statuses); err != nil {
		return nil, fmt.Errorf("failed to decode statuses response: %w", err)
	}

	return statuses, nil
}

// CreateStatus - создать новый статус задачи
func (c *taskClient) CreateStatus(req *at.CreateStatusRequest) (*at.TaskStatus, error) {
	url := fmt.Sprintf("%s/api/v1/tasks/statuses", c.host)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("request to task service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("task service returned error: %s", string(bodyBytes))
	}

	var status at.TaskStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode status response: %w", err)
	}

	return &status, nil
}

// GetStatusByID - получить статус по ID
func (c *taskClient) GetStatusByID(statusID int) (*at.TaskStatus, error) {
	url := fmt.Sprintf("%s/api/v1/tasks/statuses/%d", c.host, statusID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("task service returned error: %s", string(bodyBytes))
	}

	var status at.TaskStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode status response: %w", err)
	}

	return &status, nil
}

// DeleteStatus - удалить статус задачи
func (c *taskClient) DeleteStatus(statusID int) error {
	url := fmt.Sprintf("%s/api/v1/tasks/statuses/%d", c.host, statusID)

	httpReq, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request to task service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("task service returned error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
