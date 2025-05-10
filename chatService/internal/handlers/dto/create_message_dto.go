package dto

type CreateMessageDTO struct {
	Content string `json:"content" binding:"required"`
	FileIDs []int  `json:"fileIDs"`
}
