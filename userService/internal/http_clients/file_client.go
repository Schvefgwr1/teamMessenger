package http_clients

import (
	ufc "common/contracts/user-file-contracts"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// GetFileByID делает HTTP-запрос к файловому сервису и возвращает структуру File
func GetFileByID(fileID int) (*ufc.File, error) {
	url := fmt.Sprintf("http://localhost:8080/api/v1/files/%d", fileID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error in request's processing: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("can't get correct file: " + resp.Status)
	}

	var file ufc.File
	if err := json.NewDecoder(resp.Body).Decode(&file); err != nil {
		return nil, fmt.Errorf("error of JSON encoding: %w", err)
	}

	return &file, nil
}
