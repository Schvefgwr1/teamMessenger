package dto

import "github.com/google/uuid"

type UpdateUser struct {
	UserID uuid.UUID `json:"userID"`
	State  string    `json:"state"`
}

type UpdateChatResponse struct {
	Chat        ChatResponse `json:"chat"`
	UpdateUsers []UpdateUser `json:"updateUsers"`
}
