package custom_errors

import "fmt"

var (
	ErrNilUserInClient = "User nil error"
)

type FileSource string

const (
	FileSourceMinIO FileSource = "MIN_IO"
	FileSourceDB    FileSource = "DB"
)

type FileServiceConflictError struct {
	name   string
	source FileSource
}

func (e *FileServiceConflictError) Error() string {
	return fmt.Sprintf("File with name %s conflict in source %s", e.name, e.source)
}

func NewFileServiceConflictError(name string, source FileSource) *FileServiceConflictError {
	return &FileServiceConflictError{name: name, source: source}
}
