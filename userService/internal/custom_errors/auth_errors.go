package custom_errors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenGeneration    = errors.New("token generation failed")
)

type UserEmailConflictError struct {
	email string `binding:"email"`
}

func (e *UserEmailConflictError) Error() string {
	return fmt.Sprintf("user with email %s already exists in database", e.email)
}
func NewUserEmailConflictError(email string) *UserEmailConflictError {
	return &UserEmailConflictError{email}
}

type UserUsernameConflictError struct {
	username string
}

func (e *UserUsernameConflictError) Error() string {
	return fmt.Sprintf("user with username %s already exists in database", e.username)
}
func NewUserUsernameConflictError(username string) *UserUsernameConflictError {
	return &UserUsernameConflictError{username}
}

type RoleNotFoundError struct {
	roleId int
}

func (e *RoleNotFoundError) Error() string {
	return fmt.Sprintf("role with id %d does not exist", e.roleId)
}
func NewRoleNotFoundError(roleId int) *RoleNotFoundError {
	return &RoleNotFoundError{roleId}
}

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
