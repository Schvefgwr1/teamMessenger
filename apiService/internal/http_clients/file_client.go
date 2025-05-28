package http_clients

import (
	"apiService/internal/custom_errors"
	"bytes"
	af "common/contracts/api-file"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"
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

	// Определим MIME-тип (можно на основе расширения)
	contentType := mime.TypeByExtension(filepath.Ext(file.Filename))
	if contentType == "" {
		// fallback
		contentType = "application/octet-stream"
	}

	// Создаём кастомную часть с Content-Disposition и Content-Type
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, file.Filename))
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file part: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	if _, err := io.Copy(part, src); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

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
		responseBody, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusConflict {
			if strings.Contains(string(responseBody), "already exists in database") {
				return nil, custom_errors.NewFileServiceConflictError(file.Filename, custom_errors.FileSourceDB)
			}
			if strings.Contains(string(responseBody), "already exists in MinIO") {
				return nil, custom_errors.NewFileServiceConflictError(file.Filename, custom_errors.FileSourceMinIO)
			}
		}
		return nil, fmt.Errorf("file upload failed: %s", string(responseBody))
	}

	var uploadedFile af.FileUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadedFile); err != nil {
		return nil, fmt.Errorf("failed to parse file service response: %w", err)
	}

	return &uploadedFile, nil
}
