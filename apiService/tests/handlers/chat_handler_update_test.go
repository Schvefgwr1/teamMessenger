package handlers

import (
	"apiService/internal/handlers"
	"bytes"
	ac "common/contracts/api-chat"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Тесты для ChatHandler.UpdateChat

func TestChatHandler_UpdateChat_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()
	expectedResponse := &ac.UpdateChatResponse{
		Chat: ac.ChatResponse{
			ID:   chatID,
			Name: "Updated Chat",
		},
	}

	mockController.On("UpdateChat", chatID, mock.Anything, mock.Anything, userID).Return(expectedResponse, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.PATCH("/chats/:chat_id", handler.UpdateChat)

	// Act - создаем multipart форму
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "Updated Chat")
	writer.WriteField("description", "Updated Description")
	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/chats/"+chatID.String(), body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	mockController.AssertExpectations(t)
}

func TestChatHandler_UpdateChat_InvalidChatID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.PATCH("/chats/:chat_id", handler.UpdateChat)

	// Act
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "Updated Chat")
	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/chats/invalid", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "UpdateChat", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestChatHandler_UpdateChat_MissingUserID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()

	router := gin.New()
	router.PATCH("/chats/:chat_id", handler.UpdateChat)

	// Act
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "Updated Chat")
	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/chats/"+chatID.String(), body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockController.AssertNotCalled(t, "UpdateChat", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}
