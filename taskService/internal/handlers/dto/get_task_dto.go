package dto

import (
	fc "common/contracts/file-contracts"
	"taskService/internal/models"
)

type TaskResponse struct {
	Task  *models.Task `json:"task"`
	Files *[]fc.File   `json:"files"`
}
