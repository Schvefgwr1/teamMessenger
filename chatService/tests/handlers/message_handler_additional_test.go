package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"chatService/internal/custom_errors"
	"chatService/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMessageHandler_GetChatMessages_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/messages/:chat_id", handler.GetChatMessages)

	chatID := uuid.New()
	mockController.On("GetChatMessages", chatID, 0, 20).Return(nil, custom_errors.ErrInvalidCredentials)

	req, _ := http.NewRequest("GET", "/chats/messages/"+chatID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertExpectations(t)
}

func TestMessageHandler_SendMessage_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.POST("/chats/messages/:chat_id", handler.SendMessage)

	senderID := uuid.New()
	chatID := uuid.New()
	payload := []byte(`{"content":"test","file_ids":[]}`)

	mockController.On("SendMessage", senderID, chatID, mock.Anything).Return(nil, custom_errors.ErrInvalidCredentials)

	req, _ := http.NewRequest("POST", "/chats/messages/"+chatID.String(), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", senderID.String())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockController.AssertExpectations(t)
}

func TestMessageHandler_SendMessage_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.POST("/chats/messages/:chat_id", handler.SendMessage)

	senderID := uuid.New()
	chatID := uuid.New()
	payload := []byte(`{"content":"test","file_ids":[]}`)

	mockController.On("SendMessage", senderID, chatID, mock.Anything).Return(nil, custom_errors.NewDatabaseError("database error"))

	req, _ := http.NewRequest("POST", "/chats/messages/"+chatID.String(), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", senderID.String())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockController.AssertExpectations(t)
}

func TestMessageHandler_GetChatMessages_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/messages/:chat_id", handler.GetChatMessages)

	chatID := uuid.New()
	mockController.On("GetChatMessages", chatID, 0, 20).Return(nil, custom_errors.NewDatabaseError("database error"))

	req, _ := http.NewRequest("GET", "/chats/messages/"+chatID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockController.AssertExpectations(t)
}
