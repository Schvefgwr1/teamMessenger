package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"chatService/internal/custom_errors"
	"chatService/internal/handlers"
	"chatService/internal/handlers/dto"
	"chatService/internal/models"
	ac "common/contracts/api-chat"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockMessageController - мок для MessageControllerInterface
type MockMessageController struct {
	mock.Mock
}

func (m *MockMessageController) SendMessage(senderID, chatID uuid.UUID, dto *dto.CreateMessageDTO) (*models.Message, error) {
	args := m.Called(senderID, chatID, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Message), args.Error(1)
}

func (m *MockMessageController) GetChatMessages(chatID uuid.UUID, offset, limit int) (*[]dto.GetChatMessage, error) {
	args := m.Called(chatID, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]dto.GetChatMessage), args.Error(1)
}

func (m *MockMessageController) SearchMessages(userID, chatID uuid.UUID, query string, limit, offset int) (*ac.GetSearchResponse, error) {
	args := m.Called(userID, chatID, query, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ac.GetSearchResponse), args.Error(1)
}

// helpers
func createTestMessageModel() *models.Message {
	id := uuid.New()
	chatID := uuid.New()
	senderID := uuid.New()
	return &models.Message{
		ID:       id,
		ChatID:   chatID,
		SenderID: &senderID,
		Content:  "hi",
	}
}

func createTestGetChatMessage() dto.GetChatMessage {
	msg := createTestMessageModel()
	return dto.GetChatMessage{
		ID:        msg.ID,
		ChatID:    msg.ChatID,
		SenderID:  msg.SenderID,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
		Files:     nil,
	}
}

// --- Tests ---

func TestMessageHandler_SendMessage_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	chatID := uuid.New()
	senderID := uuid.New()
	msg := createTestMessageModel()
	reqDTO := dto.CreateMessageDTO{Content: "hi"}
	payload, _ := json.Marshal(reqDTO)

	mockController.On("SendMessage", senderID, chatID, mock.AnythingOfType("*dto.CreateMessageDTO")).Return(msg, nil)

	router := gin.New()
	router.POST("/chats/messages/:chat_id", handler.SendMessage)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats/messages/"+chatID.String(), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", senderID.String())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockController.AssertExpectations(t)
}

func TestMessageHandler_SendMessage_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.POST("/chats/messages/:chat_id", handler.SendMessage)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats/messages/"+uuid.New().String(), bytes.NewBuffer([]byte(`{"content":"hi"}`)))
	req.Header.Set("X-User-ID", "invalid")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "SendMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_SendMessage_InvalidChatID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.POST("/chats/messages/:chat_id", handler.SendMessage)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats/messages/invalid", bytes.NewBuffer([]byte(`{"content":"hi"}`)))
	req.Header.Set("X-User-ID", uuid.New().String())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "SendMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_SendMessage_InvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.POST("/chats/messages/:chat_id", handler.SendMessage)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats/messages/"+uuid.New().String(), bytes.NewBuffer([]byte(`{`)))
	req.Header.Set("X-User-ID", uuid.New().String())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "SendMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_SendMessage_ControllerErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"invalid credentials", custom_errors.ErrInvalidCredentials, http.StatusUnauthorized},
		{"file not found", custom_errors.NewFileNotFoundError(1), http.StatusBadRequest},
		{"get file http", custom_errors.NewGetFileHTTPError(1, "err"), http.StatusBadGateway},
		{"db error", custom_errors.NewDatabaseError("db"), http.StatusInternalServerError},
		{"user err", custom_errors.NewUserClientError("err"), http.StatusBadRequest},
		{"unknown", errors.New("unknown"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockController := new(MockMessageController)
			handler := handlers.NewMessageHandler(mockController)

			chatID := uuid.New()
			senderID := uuid.New()
			reqDTO := dto.CreateMessageDTO{Content: "hi"}
			payload, _ := json.Marshal(reqDTO)

			mockController.On("SendMessage", senderID, chatID, mock.Anything).Return(nil, tt.err)

			router := gin.New()
			router.POST("/chats/messages/:chat_id", handler.SendMessage)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/chats/messages/"+chatID.String(), bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-ID", senderID.String())
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			mockController.AssertExpectations(t)
		})
	}
}

// GetChatMessages

func TestMessageHandler_GetChatMessages_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	chatID := uuid.New()
	msg := createTestGetChatMessage()
	resp := []dto.GetChatMessage{msg}

	mockController.On("GetChatMessages", chatID, 0, 20).Return(&resp, nil)

	router := gin.New()
	router.GET("/chats/messages/:chat_id", handler.GetChatMessages)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/messages/"+chatID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []dto.GetChatMessage
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Len(t, got, 1)
	mockController.AssertExpectations(t)
}

func TestMessageHandler_GetChatMessages_EmptyList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	chatID := uuid.New()
	emptyMessages := []dto.GetChatMessage{}
	mockController.On("GetChatMessages", chatID, 0, 20).Return(&emptyMessages, nil)

	router := gin.New()
	router.GET("/chats/messages/:chat_id", handler.GetChatMessages)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/messages/"+chatID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []dto.GetChatMessage
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp, 0)
	mockController.AssertExpectations(t)
}

func TestMessageHandler_GetChatMessages_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/messages/:chat_id", handler.GetChatMessages)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/messages/invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "GetChatMessages", mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_GetChatMessages_InvalidPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/messages/:chat_id", handler.GetChatMessages)

	chatID := uuid.New()

	cases := []string{"/chats/messages/" + chatID.String() + "?offset=-1", "/chats/messages/" + chatID.String() + "?limit=0", "/chats/messages/" + chatID.String() + "?limit=abc"}
	for _, url := range cases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	}

	mockController.AssertNotCalled(t, "GetChatMessages", mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_GetChatMessages_ControllerErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"invalid credentials", custom_errors.ErrInvalidCredentials, http.StatusBadRequest},
		{"db error", custom_errors.NewDatabaseError("db"), http.StatusInternalServerError},
		{"unknown", errors.New("unknown"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockController := new(MockMessageController)
			handler := handlers.NewMessageHandler(mockController)

			chatID := uuid.New()
			mockController.On("GetChatMessages", chatID, 0, 20).Return(nil, tt.err)

			router := gin.New()
			router.GET("/chats/messages/:chat_id", handler.GetChatMessages)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/chats/messages/"+chatID.String(), nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			mockController.AssertExpectations(t)
		})
	}
}

// SearchMessages

func TestMessageHandler_SearchMessages_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	userID := uuid.New()
	chatID := uuid.New()
	resp := &ac.GetSearchResponse{
		Messages: func() *[]ac.GetChatMessage {
			msg := createTestGetChatMessage()
			conv := ac.GetChatMessage{
				ID:        msg.ID,
				ChatID:    msg.ChatID,
				SenderID:  msg.SenderID,
				Content:   msg.Content,
				UpdatedAt: msg.UpdatedAt,
				CreatedAt: msg.CreatedAt,
			}
			return &[]ac.GetChatMessage{conv}
		}(),
		Total: func() *int64 { v := int64(1); return &v }(),
	}

	mockController.On("SearchMessages", userID, chatID, "hi", 20, 0).Return(resp, nil)

	router := gin.New()
	router.GET("/chats/search/:chat_id", handler.SearchMessages)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/search/"+chatID.String()+"?query=hi", nil)
	req.Header.Set("X-User-ID", userID.String())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got ac.GetSearchResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, int64(1), *got.Total)
	mockController.AssertExpectations(t)
}

func TestMessageHandler_SearchMessages_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/search/:chat_id", handler.SearchMessages)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/search/invalid?query=hi", nil)
	req.Header.Set("X-User-ID", uuid.New().String())
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/chats/search/"+uuid.New().String()+"?query=hi", nil)
	req2.Header.Set("X-User-ID", "invalid")
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)

	mockController.AssertNotCalled(t, "SearchMessages", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_SearchMessages_EmptyQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()

	mockController.On("SearchMessages", userID, chatID, "", 20, 0).Return(nil, custom_errors.ErrEmptyQuery)

	router := gin.New()
	router.GET("/chats/search/:chat_id", handler.SearchMessages)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/search/"+chatID.String(), nil)
	req.Header.Set("X-User-ID", userID.String())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertExpectations(t)
}

func TestMessageHandler_SearchMessages_LimitTooLarge(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/search/:chat_id", handler.SearchMessages)

	chatID := uuid.New()
	userID := uuid.New()
	total := int64(0)
	expectedResponse := &ac.GetSearchResponse{
		Messages: &[]ac.GetChatMessage{},
		Total:    &total,
	}

	mockController.On("SearchMessages", userID, chatID, "test", 50, 0).Return(expectedResponse, nil)

	req, _ := http.NewRequest("GET", "/chats/search/"+chatID.String()+"?query=test&limit=100", nil)
	req.Header.Set("X-User-ID", userID.String())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockController.AssertExpectations(t)
}

func TestMessageHandler_SearchMessages_InvalidPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockMessageController)
	handler := handlers.NewMessageHandler(mockController)

	router := gin.New()
	router.GET("/chats/search/:chat_id", handler.SearchMessages)

	chatID := uuid.New()
	userID := uuid.New()

	urls := []string{
		"/chats/search/" + chatID.String() + "?query=hi&limit=0",
		"/chats/search/" + chatID.String() + "?query=hi&offset=-1",
		"/chats/search/" + chatID.String() + "?query=hi&limit=abc",
	}
	for _, u := range urls {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", u, nil)
		req.Header.Set("X-User-ID", userID.String())
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	}

	mockController.AssertNotCalled(t, "SearchMessages", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_SearchMessages_ControllerErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"empty", custom_errors.ErrEmptyQuery, http.StatusBadRequest},
		{"notfound", custom_errors.ErrChatNotFound, http.StatusNotFound},
		{"unauth", custom_errors.ErrUnauthorizedChat, http.StatusForbidden},
		{"unknown", errors.New("unknown"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockController := new(MockMessageController)
			handler := handlers.NewMessageHandler(mockController)

			chatID := uuid.New()
			userID := uuid.New()

			mockController.On("SearchMessages", userID, chatID, "hi", 20, 0).Return(nil, tt.err)

			router := gin.New()
			router.GET("/chats/search/:chat_id", handler.SearchMessages)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/chats/search/"+chatID.String()+"?query=hi", nil)
			req.Header.Set("X-User-ID", userID.String())
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			mockController.AssertExpectations(t)
		})
	}
}
