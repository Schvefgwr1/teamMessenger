package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"chatService/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMessageHandler_SendMessage_MissingUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.POST("/chats/messages/:chat_id", handler.SendMessage)

	chatID := uuid.New()
	payload := []byte(`{"content":"test","file_ids":[]}`)

	req, _ := http.NewRequest("POST", "/chats/messages/"+chatID.String(), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "SendMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_GetChatMessages_InvalidOffset_Negative(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/messages/:chat_id", handler.GetChatMessages)

	chatID := uuid.New()

	req, _ := http.NewRequest("GET", "/chats/messages/"+chatID.String()+"?offset=-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "GetChatMessages", mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_GetChatMessages_InvalidLimit_Zero(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/messages/:chat_id", handler.GetChatMessages)

	chatID := uuid.New()

	req, _ := http.NewRequest("GET", "/chats/messages/"+chatID.String()+"?limit=0", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "GetChatMessages", mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_GetChatMessages_InvalidLimit_Negative(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/messages/:chat_id", handler.GetChatMessages)

	chatID := uuid.New()

	req, _ := http.NewRequest("GET", "/chats/messages/"+chatID.String()+"?limit=-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "GetChatMessages", mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_GetChatMessages_InvalidOffset_NonNumeric(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/messages/:chat_id", handler.GetChatMessages)

	chatID := uuid.New()

	req, _ := http.NewRequest("GET", "/chats/messages/"+chatID.String()+"?offset=abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "GetChatMessages", mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_GetChatMessages_InvalidLimit_NonNumeric(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/messages/:chat_id", handler.GetChatMessages)

	chatID := uuid.New()

	req, _ := http.NewRequest("GET", "/chats/messages/"+chatID.String()+"?limit=abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "GetChatMessages", mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_SearchMessages_MissingUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/search/:chat_id", handler.SearchMessages)

	chatID := uuid.New()

	req, _ := http.NewRequest("GET", "/chats/search/"+chatID.String()+"?query=test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "SearchMessages", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_SearchMessages_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/search/:chat_id", handler.SearchMessages)

	chatID := uuid.New()

	req, _ := http.NewRequest("GET", "/chats/search/"+chatID.String()+"?query=test", nil)
	req.Header.Set("X-User-ID", "invalid")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "SearchMessages", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_SearchMessages_InvalidOffset_Negative(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/search/:chat_id", handler.SearchMessages)

	chatID := uuid.New()
	userID := uuid.New()

	req, _ := http.NewRequest("GET", "/chats/search/"+chatID.String()+"?query=test&offset=-1", nil)
	req.Header.Set("X-User-ID", userID.String())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "SearchMessages", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_SearchMessages_InvalidLimit_Zero(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/search/:chat_id", handler.SearchMessages)

	chatID := uuid.New()
	userID := uuid.New()

	req, _ := http.NewRequest("GET", "/chats/search/"+chatID.String()+"?query=test&limit=0", nil)
	req.Header.Set("X-User-ID", userID.String())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "SearchMessages", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_SearchMessages_InvalidOffset_NonNumeric(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/search/:chat_id", handler.SearchMessages)

	chatID := uuid.New()
	userID := uuid.New()

	req, _ := http.NewRequest("GET", "/chats/search/"+chatID.String()+"?query=test&offset=abc", nil)
	req.Header.Set("X-User-ID", userID.String())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "SearchMessages", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_SearchMessages_InvalidLimit_NonNumeric(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/search/:chat_id", handler.SearchMessages)

	chatID := uuid.New()
	userID := uuid.New()

	req, _ := http.NewRequest("GET", "/chats/search/"+chatID.String()+"?query=test&limit=abc", nil)
	req.Header.Set("X-User-ID", userID.String())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "SearchMessages", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}
