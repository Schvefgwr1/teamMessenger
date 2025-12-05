package dto

import "github.com/google/uuid"

type CreateTaskDTO struct {
	Title       string    `json:"title" binding:"required"`
	Description *string   `json:"description"`
	CreatorID   uuid.UUID `json:"creator_id" binding:"required"`
	ExecutorID  uuid.UUID `json:"executor_id" binding:"required"`
	ChatID      uuid.UUID `json:"chat_id"`
	FileIDs     []int     `json:"file_ids"`
}
