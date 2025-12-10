package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"chatService/internal/handlers"
	ac "common/contracts/api-chat"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageHandler_SearchMessages_Success_WithResults(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()
	messages := []ac.GetChatMessage{
		{ID: uuid.New(), Content: "test message"},
	}
	total := int64(1)
	response := &ac.GetSearchResponse{
		Messages: &messages,
		Total:    &total,
	}

	mockController.On("SearchMessages", userID, chatID, "test", 20, 0).Return(response, nil)

	router := gin.New()
	router.GET("/chats/search/:chat_id", handler.SearchMessages)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/search/"+chatID.String()+"?query=test", nil)
	req.Header.Set("X-User-ID", userID.String())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp ac.GetSearchResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotNil(t, resp.Messages)
	assert.Len(t, *resp.Messages, 1)
	mockController.AssertExpectations(t)
}

func TestMessageHandler_SearchMessages_UnknownError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()

	mockController.On("SearchMessages", userID, chatID, "test", 20, 0).Return(nil, errors.New("unknown error"))

	router := gin.New()
	router.GET("/chats/search/:chat_id", handler.SearchMessages)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/search/"+chatID.String()+"?query=test", nil)
	req.Header.Set("X-User-ID", userID.String())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockController.AssertExpectations(t)
}
