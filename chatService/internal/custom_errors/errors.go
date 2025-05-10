package custom_errors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInternalServerError = errors.New("internal server error")
	ErrEmptyQuery          = errors.New("query parameter cannot be empty")
	ErrChatNotFound        = errors.New("chat with provided ID not found")
	ErrUnauthorizedChat    = errors.New("user is not a member of this chat")
)

type GetFileHTTPError struct {
	httpError string
	fileId    int
}

func (e *GetFileHTTPError) Error() string {
	return fmt.Sprintf("can't get file with id: %d with error: %s", e.fileId, e.httpError)
}
func NewGetFileHTTPError(fileId int, httpError string) *GetFileHTTPError {
	return &GetFileHTTPError{fileId: fileId, httpError: httpError}
}

type FileNotFoundError struct {
	fileId int
}

func (e *FileNotFoundError) Error() string {
	return fmt.Sprintf("can't get file with incorrect id: %d", e.fileId)
}
func NewFileNotFoundError(fileId int) *FileNotFoundError {
	return &FileNotFoundError{fileId}
}

type UserClientError struct {
	description string
}

func (e *UserClientError) Error() string {
	return fmt.Sprintf("incorrect work of user http client: %s", e.description)

}
func NewUserClientError(error string) *UserClientError {
	return &UserClientError{error}
}

type DatabaseError struct {
	description string
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("error of db: %s", e.description)

}
func NewDatabaseError(error string) *DatabaseError {
	return &DatabaseError{error}
}
