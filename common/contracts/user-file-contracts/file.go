package user_file_contracts

import "time"

// File — описание файла
type File struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	FileTypeID int       `json:"file_type_id"`
	URL        string    `json:"url"`
	CreatedAt  time.Time `json:"created_at"`
	FileType   FileType  `json:"file_type"`
}
