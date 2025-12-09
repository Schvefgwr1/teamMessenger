package controllers

import (
	"errors"
	ctrl "fileService/internal/controllers"
	"fileService/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Тесты для FileTypeController

func TestFileTypeController_CreateFileType_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockFileTypeRepository)
	controller := ctrl.NewFileTypeController(mockRepo)

	name := "application/json"
	fileType := &models.FileType{
		ID:   1,
		Name: name,
	}

	// Настраиваем моки
	mockRepo.On("CreateFileType", mock.MatchedBy(func(ft *models.FileType) bool {
		return ft.Name == name
	})).Return(nil).Run(func(args mock.Arguments) {
		ft := args.Get(0).(*models.FileType)
		ft.ID = fileType.ID
	})

	// Act
	result, err := controller.CreateFileType(name)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, name, result.Name)

	mockRepo.AssertExpectations(t)
}

func TestFileTypeController_CreateFileType_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockFileTypeRepository)
	controller := ctrl.NewFileTypeController(mockRepo)

	name := "application/json"
	repoError := errors.New("database error")

	// Настраиваем моки
	mockRepo.On("CreateFileType", mock.Anything).Return(repoError)

	// Act
	result, err := controller.CreateFileType(name)

	// Assert
	require.Error(t, err)
	// Контроллер всегда возвращает fileType, даже при ошибке
	assert.NotNil(t, result)
	assert.Equal(t, name, result.Name)
	assert.Equal(t, repoError, err)

	mockRepo.AssertExpectations(t)
}

func TestFileTypeController_GetFileTypeByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockFileTypeRepository)
	controller := ctrl.NewFileTypeController(mockRepo)

	fileTypeID := 1
	expectedFileType := &models.FileType{
		ID:   fileTypeID,
		Name: "text/plain",
	}

	// Настраиваем моки
	mockRepo.On("GetFileTypeByID", fileTypeID).Return(expectedFileType, nil)

	// Act
	result, err := controller.GetFileTypeByID(fileTypeID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedFileType, result)

	mockRepo.AssertExpectations(t)
}

func TestFileTypeController_GetFileTypeByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockFileTypeRepository)
	controller := ctrl.NewFileTypeController(mockRepo)

	fileTypeID := 999
	repoError := errors.New("record not found")

	// Настраиваем моки
	mockRepo.On("GetFileTypeByID", fileTypeID).Return(nil, repoError)

	// Act
	result, err := controller.GetFileTypeByID(fileTypeID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoError, err)

	mockRepo.AssertExpectations(t)
}

func TestFileTypeController_GetFileTypeByName_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockFileTypeRepository)
	controller := ctrl.NewFileTypeController(mockRepo)

	name := "plain"
	fullName := "application/plain"
	expectedFileType := &models.FileType{
		ID:   1,
		Name: fullName,
	}

	// Настраиваем моки - контроллер добавляет "application/" к имени
	mockRepo.On("GetFileTypeByName", fullName).Return(expectedFileType, nil)

	// Act
	result, err := controller.GetFileTypeByName(name)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedFileType, result)

	mockRepo.AssertExpectations(t)
}

func TestFileTypeController_GetFileTypeByName_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockFileTypeRepository)
	controller := ctrl.NewFileTypeController(mockRepo)

	name := "nonexistent"
	fullName := "application/nonexistent"
	repoError := errors.New("record not found")

	// Настраиваем моки
	mockRepo.On("GetFileTypeByName", fullName).Return(nil, repoError)

	// Act
	result, err := controller.GetFileTypeByName(name)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoError, err)

	mockRepo.AssertExpectations(t)
}

func TestFileTypeController_DeleteFileTypeByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockFileTypeRepository)
	controller := ctrl.NewFileTypeController(mockRepo)

	fileTypeID := 1

	// Настраиваем моки
	mockRepo.On("DeleteFileTypeByID", fileTypeID).Return(nil)

	// Act
	err := controller.DeleteFileTypeByID(fileTypeID)

	// Assert
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestFileTypeController_DeleteFileTypeByID_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockFileTypeRepository)
	controller := ctrl.NewFileTypeController(mockRepo)

	fileTypeID := 999
	repoError := errors.New("database error")

	// Настраиваем моки
	mockRepo.On("DeleteFileTypeByID", fileTypeID).Return(repoError)

	// Act
	err := controller.DeleteFileTypeByID(fileTypeID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, repoError, err)

	mockRepo.AssertExpectations(t)
}
