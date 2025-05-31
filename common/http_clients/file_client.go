package http_clients

import (
	"common/config"
	fc "common/contracts/file-contracts"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// GetFileByID делает HTTP-запрос к файловому сервису и возвращает структуру File
func GetFileByID(fileID int) (*fc.File, error) {
	baseURL := config.GetEnvOrDefault("FILE_SERVICE_URL", "http://localhost:8080")
	url := fmt.Sprintf("%s/api/v1/files/%d", baseURL, fileID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error in request's processing: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("can't get correct file: " + resp.Status)
	}

	var file fc.File
	if err := json.NewDecoder(resp.Body).Decode(&file); err != nil {
		return nil, fmt.Errorf("error of JSON encoding: %w", err)
	}

	return &file, nil
}
