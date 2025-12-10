package handlers

import (
	"apiService/internal/handlers"
	ac "common/contracts/api-chat"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Тесты для ChatHandler.GetUserChats

func TestChatHandler_GetUserChats_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	userID := uuid.New()
	expectedChats := []*ac.ChatResponse{
		{ID: uuid.New(), Name: "Test Chat"},
	}

	mockController.On("GetUserChats", userID).Return(expectedChats, nil)

	router := gin.New()
	router.GET("/chats/:user_id", handler.GetUserChats)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/"+userID.String(), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []*ac.ChatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, expectedChats[0].Name, response[0].Name)

	mockController.AssertExpectations(t)
}

func TestChatHandler_GetUserChats_InvalidUserID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.GET("/chats/:user_id", handler.GetUserChats)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/invalid-uuid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "invalid user ID")

	mockController.AssertNotCalled(t, "GetUserChats", mock.Anything)
}

func TestChatHandler_GetUserChats_ControllerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	userID := uuid.New()
	serviceError := errors.New("service error")

	mockController.On("GetUserChats", userID).Return(nil, serviceError)

	router := gin.New()
	router.GET("/chats/:user_id", handler.GetUserChats)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/"+userID.String(), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, serviceError.Error(), response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для ChatHandler.GetChatMessages

func TestChatHandler_GetChatMessages_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()
	expectedMessages := []*ac.GetChatMessage{
		{ID: uuid.New(), Content: "Test message"},
	}

	mockController.On("GetChatMessages", chatID, userID, 0, 20).Return(expectedMessages, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.GET("/chats/messages/:chat_id", handler.GetChatMessages)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/messages/"+chatID.String()+"?offset=0&limit=20", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []*ac.GetChatMessage
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 1)

	mockController.AssertExpectations(t)
}

func TestChatHandler_GetChatMessages_InvalidChatID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.GET("/chats/messages/:chat_id", handler.GetChatMessages)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/messages/invalid-uuid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "GetChatMessages", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestChatHandler_GetChatMessages_InvalidOffset(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.GET("/chats/messages/:chat_id", handler.GetChatMessages)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/messages/"+chatID.String()+"?offset=-1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "GetChatMessages", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestChatHandler_GetChatMessages_InvalidLimit(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.GET("/chats/messages/:chat_id", handler.GetChatMessages)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/messages/"+chatID.String()+"?limit=101", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "GetChatMessages", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

// Тесты для ChatHandler.SendMessage

func TestChatHandler_SendMessage_MissingUserID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()

	router := gin.New()
	router.POST("/chats/messages/:chat_id", handler.SendMessage)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats/messages/"+chatID.String(), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	mockController.AssertNotCalled(t, "SendMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestChatHandler_SendMessage_InvalidChatID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.POST("/chats/messages/:chat_id", handler.SendMessage)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats/messages/invalid-uuid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "SendMessage", mock.Anything, mock.Anything, mock.Anything)
}

// Тесты для ChatHandler.SearchMessages

func TestChatHandler_SearchMessages_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()
	query := "test"
	messages := []ac.GetChatMessage{
		{ID: uuid.New(), Content: "test message"},
	}
	total := int64(1)
	expectedResult := &ac.GetSearchResponse{
		Messages: &messages,
		Total:    &total,
	}

	mockController.On("SearchMessages", userID, chatID, query, 0, 20).Return(expectedResult, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.GET("/chats/search/:chat_id", handler.SearchMessages)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/search/"+chatID.String()+"?query="+query, nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response ac.GetSearchResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotNil(t, response.Messages)

	mockController.AssertExpectations(t)
}

func TestChatHandler_SearchMessages_MissingQuery(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.GET("/chats/search/:chat_id", handler.SearchMessages)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/search/"+chatID.String(), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "SearchMessages", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

// Тесты для ChatHandler.DeleteChat

func TestChatHandler_DeleteChat_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()

	mockController.On("DeleteChat", chatID, userID).Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.DELETE("/chats/:chat_id", handler.DeleteChat)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/chats/"+chatID.String(), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "chat deleted successfully", response["message"])

	mockController.AssertExpectations(t)
}

// Тесты для ChatHandler.GetChatMembers

func TestChatHandler_GetChatMembers_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	expectedMembers := []*ac.ChatMember{
		{UserID: uuid.New().String(), RoleID: 1, RoleName: "user"},
	}

	mockController.On("GetChatMembers", chatID).Return(expectedMembers, nil)

	router := gin.New()
	router.GET("/chats/:chat_id/members", handler.GetChatMembers)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/"+chatID.String()+"/members", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []*ac.ChatMember
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 1)

	mockController.AssertExpectations(t)
}
