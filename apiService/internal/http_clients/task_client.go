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
