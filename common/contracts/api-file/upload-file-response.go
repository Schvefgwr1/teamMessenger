package api_file

import "time"

type FileType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type FileUploadResponse struct {
	ID         *int       `json:"id"`
	Name       *string    `json:"name"`
	FileTypeID *int       `json:"file_type_id"`
	URL        *string    `json:"url"`
	CreatedAt  *time.Time `json:"created_at"`
	FileType   *FileType  `json:"file_type"`
	Error      *string    `json:"error"`
}
