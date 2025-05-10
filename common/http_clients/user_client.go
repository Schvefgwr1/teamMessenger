package http_clients

import (
	cuc "common/contracts/user-contracts"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

// GetUserByID делает HTTP-запрос к файловому сервису и возвращает структуру File
func GetUserByID(userID *uuid.UUID) (*cuc.Response, error) {
	url := fmt.Sprintf("http://localhost:8082/api/v1/users/%s", (*userID).String())

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error in request's processing: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("can't get user: " + resp.Status)
	}

	var dtoResp cuc.Response
	if err := json.NewDecoder(resp.Body).Decode(&dtoResp); err != nil {
		return nil, fmt.Errorf("error of JSON encoding: %w", err)
	}

	return &dtoResp, nil
}
