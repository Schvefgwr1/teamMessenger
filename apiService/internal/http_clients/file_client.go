package http_clients

import (
	"bytes"
	af "common/contracts/api-file"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type FileClient interface {
	UploadFile(file *multipart.FileHeader) (*af.FileUploadResponse, error)
}

type fileClient struct {
	host string
}

func NewFileClient(host string) FileClient {
	return &fileClient{host: host}
}

func (c *fileClient) UploadFile(file *multipart.FileHeader) (*af.FileUploadResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	if _, err = io.Copy(part, src); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", c.host+"/api/v1/files/upload", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("file upload request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("file upload failed: %s", string(body))
	}

	var uploadedFile af.FileUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadedFile); err != nil {
		return nil, fmt.Errorf("failed to parse file service response: %w", err)
	}

	return &uploadedFile, nil
}
