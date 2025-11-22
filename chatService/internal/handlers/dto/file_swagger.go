package dto

import (
	"time"
)

// FileSwagger — подмена структуры для Swagger, аналог fc.File
type FileSwagger struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	FileTypeID int       `json:"file_type_id"`
	URL        string    `json:"url"`
	CreatedAt  time.Time `json:"created_at"`
	FileType   any       `json:"file_type"`
}

type GetSearchResponse struct {
	Messages *[]GetChatMessage `json:"messages"`
	Total    *int64            `json:"total"`
}
