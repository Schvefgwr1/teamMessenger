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
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockChatController - мок для ChatController
type MockChatController struct {
	mock.Mock
}

func (m *MockChatController) ChangeUserRole(chatID, userID uuid.UUID, roleID int) error {
	args := m.Called(chatID, userID, roleID)
	return args.Error(0)
}

func (m *MockChatController) GetUserChats(userID uuid.UUID) (*[]dto.ChatResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]dto.ChatResponse), args.Error(1)
}

func (m *MockChatController) CreateChat(createDTO *dto.CreateChatDTO) (*uuid.UUID, error) {
	args := m.Called(createDTO)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*uuid.UUID), args.Error(1)
}

func (m *MockChatController) UpdateChat(chatID uuid.UUID, updateDTO *dto.UpdateChatDTO) (*dto.UpdateChatResponse, error) {
	args := m.Called(chatID, updateDTO)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UpdateChatResponse), args.Error(1)
}

func (m *MockChatController) DeleteChat(chatID uuid.UUID) error {
	args := m.Called(chatID)
	return args.Error(0)
}

func (m *MockChatController) BanUser(chatID, userID uuid.UUID) error {
	args := m.Called(chatID, userID)
	return args.Error(0)
}

func (m *MockChatController) GetUserRoleInChat(chatID, userID, requesterID uuid.UUID) (string, error) {
	args := m.Called(chatID, userID, requesterID)
	return args.String(0), args.Error(1)
}

func (m *MockChatController) GetMyRoleWithPermissions(chatID, userID uuid.UUID) (*models.ChatRole, error) {
	args := m.Called(chatID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatRole), args.Error(1)
}

func (m *MockChatController) GetChatByID(chatID uuid.UUID) (*dto.ChatResponse, error) {
	args := m.Called(chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ChatResponse), args.Error(1)
}

func (m *MockChatController) GetChatMembers(chatID uuid.UUID) ([]models.ChatUser, error) {
	args := m.Called(chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ChatUser), args.Error(1)
}

// Вспомогательные функции
func createTestChatResponse() dto.ChatResponse {
	return dto.ChatResponse{
		ID:   uuid.New(),
		Name: "Test Chat",
	}
}

func createTestUpdateChatResponse() *dto.UpdateChatResponse {
	return &dto.UpdateChatResponse{
		Chat: createTestChatResponse(),
		UpdateUsers: []dto.UpdateUser{
			{UserID: uuid.New(), State: "created"},
		},
	}
}

func createTestChatRoleWithPermissions() *models.ChatRole {
	return &models.ChatRole{
		ID:   1,
		Name: "main",
		Permissions: []models.ChatPermission{
			{ID: 1, Name: "write"},
		},
	}
}

func stringPtr(s string) *string {
	return &s
}

// --- Tests ---

func TestChatHandler_ChangeUserRole_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()
	body := map[string]interface{}{
		"user_id": userID,
		"role_id": 2,
	}
	payload, _ := json.Marshal(body)

	mockController.On("ChangeUserRole", chatID, userID, 2).Return(nil)

	router := gin.New()
	router.PATCH("/chats/:chat_id/roles/change", handler.ChangeUserRole)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/chats/"+chatID.String()+"/roles/change", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_ChangeUserRole_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.PATCH("/chats/:chat_id/roles/change", handler.ChangeUserRole)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/chats/"+uuid.New().String()+"/roles/change", bytes.NewBuffer([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "ChangeUserRole", mock.Anything, mock.Anything, mock.Anything)
}

func TestChatHandler_ChangeUserRole_InvalidChatID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	body := map[string]interface{}{
		"user_id": uuid.New(),
		"role_id": 2,
	}
	payload, _ := json.Marshal(body)

	router := gin.New()
	router.PATCH("/chats/:chat_id/roles/change", handler.ChangeUserRole)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/chats/invalid/roles/change", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "ChangeUserRole", mock.Anything, mock.Anything, mock.Anything)
}

func TestChatHandler_ChangeUserRole_ControllerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()
	body := map[string]interface{}{
		"user_id": userID,
		"role_id": 2,
	}
	payload, _ := json.Marshal(body)

	mockController.On("ChangeUserRole", chatID, userID, 2).Return(errors.New("db"))

	router := gin.New()
	router.PATCH("/chats/:chat_id/roles/change", handler.ChangeUserRole)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/chats/"+chatID.String()+"/roles/change", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_GetUserChats_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	userID := uuid.New()
	chats := []dto.ChatResponse{createTestChatResponse()}
	mockController.On("GetUserChats", userID).Return(&chats, nil)

	router := gin.New()
	router.GET("/chats/user/:user_id", handler.GetUserChats)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/user/"+userID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []dto.ChatResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp, 1)

	mockController.AssertExpectations(t)
}

func TestChatHandler_GetUserChats_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.GET("/chats/user/:user_id", handler.GetUserChats)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/user/invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "GetUserChats", mock.Anything)
}

func TestChatHandler_GetUserChats_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	userID := uuid.New()
	mockController.On("GetUserChats", userID).Return(nil, errors.New("db"))

	router := gin.New()
	router.GET("/chats/user/:user_id", handler.GetUserChats)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/user/"+userID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_GetUserChats_EmptyList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	userID := uuid.New()
	emptyChats := []dto.ChatResponse{}
	mockController.On("GetUserChats", userID).Return(&emptyChats, nil)

	router := gin.New()
	router.GET("/chats/user/:user_id", handler.GetUserChats)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/user/"+userID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []dto.ChatResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp, 0)

	mockController.AssertExpectations(t)
}

func TestChatHandler_GetChatMembers_EmptyList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	emptyMembers := []models.ChatUser{}
	mockController.On("GetChatMembers", chatID).Return(emptyMembers, nil)

	router := gin.New()
	router.GET("/chats/:chat_id/members", handler.GetChatMembers)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/"+chatID.String()+"/members", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []dto.ChatMemberResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp, 0)

	mockController.AssertExpectations(t)
}

func TestChatHandler_CreateChat_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	reqDTO := dto.CreateChatDTO{
		Name:    "chat",
		OwnerID: uuid.New(),
		UserIDs: []uuid.UUID{},
	}
	payload, _ := json.Marshal(reqDTO)

	mockController.On("CreateChat", mock.AnythingOfType("*dto.CreateChatDTO")).Return(&chatID, nil)

	router := gin.New()
	router.POST("/chats", handler.CreateChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, chatID.String(), resp["chat_id"])

	mockController.AssertExpectations(t)
}

func TestChatHandler_CreateChat_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.POST("/chats", handler.CreateChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats", bytes.NewBuffer([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "CreateChat", mock.Anything)
}

func TestChatHandler_CreateChat_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	reqDTO := dto.CreateChatDTO{Name: "chat", OwnerID: uuid.New()}
	payload, _ := json.Marshal(reqDTO)

	mockController.On("CreateChat", mock.Anything).Return(nil, custom_errors.ErrInvalidCredentials)

	router := gin.New()
	router.POST("/chats", handler.CreateChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_CreateChat_GetFileHTTPError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	reqDTO := dto.CreateChatDTO{Name: "chat", OwnerID: uuid.New()}
	payload, _ := json.Marshal(reqDTO)
	mockController.On("CreateChat", mock.Anything).Return(nil, custom_errors.NewGetFileHTTPError(1, "err"))

	router := gin.New()
	router.POST("/chats", handler.CreateChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadGateway, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_CreateChat_FileNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	reqDTO := dto.CreateChatDTO{Name: "chat", OwnerID: uuid.New()}
	payload, _ := json.Marshal(reqDTO)
	mockController.On("CreateChat", mock.Anything).Return(nil, custom_errors.NewFileNotFoundError(1))

	router := gin.New()
	router.POST("/chats", handler.CreateChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_CreateChat_UserClientError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	reqDTO := dto.CreateChatDTO{Name: "chat", OwnerID: uuid.New()}
	payload, _ := json.Marshal(reqDTO)
	mockController.On("CreateChat", mock.Anything).Return(nil, custom_errors.NewUserClientError("err"))

	router := gin.New()
	router.POST("/chats", handler.CreateChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadGateway, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_CreateChat_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	reqDTO := dto.CreateChatDTO{Name: "chat", OwnerID: uuid.New()}
	payload, _ := json.Marshal(reqDTO)
	mockController.On("CreateChat", mock.Anything).Return(nil, custom_errors.NewDatabaseError("db"))

	router := gin.New()
	router.POST("/chats", handler.CreateChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_CreateChat_UnknownError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	reqDTO := dto.CreateChatDTO{Name: "chat", OwnerID: uuid.New()}
	payload, _ := json.Marshal(reqDTO)
	mockController.On("CreateChat", mock.Anything).Return(nil, errors.New("unknown"))

	router := gin.New()
	router.POST("/chats", handler.CreateChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_UpdateChat_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	reqDTO := dto.UpdateChatDTO{Name: stringPtr("new")}
	payload, _ := json.Marshal(reqDTO)
	updateResp := createTestUpdateChatResponse()

	mockController.On("UpdateChat", chatID, mock.AnythingOfType("*dto.UpdateChatDTO")).Return(updateResp, nil)

	router := gin.New()
	router.PUT("/chats/:chat_id", handler.UpdateChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/chats/"+chatID.String(), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.UpdateChatResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, updateResp.Chat.ID, resp.Chat.ID)

	mockController.AssertExpectations(t)
}

func TestChatHandler_UpdateChat_InvalidChatID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	reqDTO := dto.UpdateChatDTO{Name: stringPtr("new")}
	payload, _ := json.Marshal(reqDTO)

	router := gin.New()
	router.PUT("/chats/:chat_id", handler.UpdateChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/chats/invalid", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "UpdateChat", mock.Anything, mock.Anything)
}

func TestChatHandler_UpdateChat_InvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.PUT("/chats/:chat_id", handler.UpdateChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/chats/"+uuid.New().String(), bytes.NewBuffer([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "UpdateChat", mock.Anything, mock.Anything)
}

func TestChatHandler_UpdateChat_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	reqDTO := dto.UpdateChatDTO{Name: stringPtr("new")}
	payload, _ := json.Marshal(reqDTO)

	mockController.On("UpdateChat", chatID, mock.Anything).Return(nil, custom_errors.ErrInvalidCredentials)

	router := gin.New()
	router.PUT("/chats/:chat_id", handler.UpdateChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/chats/"+chatID.String(), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_UpdateChat_GetFileHTTPError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	reqDTO := dto.UpdateChatDTO{Name: stringPtr("new")}
	payload, _ := json.Marshal(reqDTO)

	mockController.On("UpdateChat", chatID, mock.Anything).Return(nil, custom_errors.NewGetFileHTTPError(1, "err"))

	router := gin.New()
	router.PUT("/chats/:chat_id", handler.UpdateChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/chats/"+chatID.String(), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadGateway, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_UpdateChat_FileNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	reqDTO := dto.UpdateChatDTO{Name: stringPtr("new")}
	payload, _ := json.Marshal(reqDTO)

	mockController.On("UpdateChat", chatID, mock.Anything).Return(nil, custom_errors.NewFileNotFoundError(1))

	router := gin.New()
	router.PUT("/chats/:chat_id", handler.UpdateChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/chats/"+chatID.String(), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_UpdateChat_UserClientError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	reqDTO := dto.UpdateChatDTO{Name: stringPtr("new")}
	payload, _ := json.Marshal(reqDTO)

	mockController.On("UpdateChat", chatID, mock.Anything).Return(nil, custom_errors.NewUserClientError("err"))

	router := gin.New()
	router.PUT("/chats/:chat_id", handler.UpdateChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/chats/"+chatID.String(), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadGateway, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_UpdateChat_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	reqDTO := dto.UpdateChatDTO{Name: stringPtr("new")}
	payload, _ := json.Marshal(reqDTO)

	mockController.On("UpdateChat", chatID, mock.Anything).Return(nil, custom_errors.NewDatabaseError("db"))

	router := gin.New()
	router.PUT("/chats/:chat_id", handler.UpdateChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/chats/"+chatID.String(), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_UpdateChat_UnknownError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	reqDTO := dto.UpdateChatDTO{Name: stringPtr("new")}
	payload, _ := json.Marshal(reqDTO)

	mockController.On("UpdateChat", chatID, mock.Anything).Return(nil, errors.New("unknown"))

	router := gin.New()
	router.PUT("/chats/:chat_id", handler.UpdateChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/chats/"+chatID.String(), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_DeleteChat_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	mockController.On("DeleteChat", chatID).Return(nil)

	router := gin.New()
	router.DELETE("/chats/:chat_id", handler.DeleteChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/chats/"+chatID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_DeleteChat_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	mockController.On("DeleteChat", chatID).Return(errors.New("db"))

	router := gin.New()
	router.DELETE("/chats/:chat_id", handler.DeleteChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/chats/"+chatID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_BanUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()

	mockController.On("BanUser", chatID, userID).Return(nil)

	router := gin.New()
	router.POST("/chats/:chat_id/ban/:user_id", handler.BanUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats/"+chatID.String()+"/ban/"+userID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_BanUser_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.POST("/chats/:chat_id/ban/:user_id", handler.BanUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats/invalid/ban/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("POST", "/chats/"+uuid.New().String()+"/ban/invalid", nil)
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestChatHandler_BanUser_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()
	mockController.On("BanUser", chatID, userID).Return(errors.New("db"))

	router := gin.New()
	router.POST("/chats/:chat_id/ban/:user_id", handler.BanUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats/"+chatID.String()+"/ban/"+userID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_GetUserRoleInChat_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()
	requesterID := uuid.New()

	mockController.On("GetUserRoleInChat", chatID, userID, requesterID).Return("admin", nil)

	router := gin.New()
	router.GET("/chats/:chat_id/user-roles/:user_id", handler.GetUserRoleInChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/"+chatID.String()+"/user-roles/"+userID.String(), nil)
	req.Header.Set("X-User-ID", requesterID.String())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.UserRoleResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "admin", resp.RoleName)

	mockController.AssertExpectations(t)
}

func TestChatHandler_GetUserRoleInChat_InvalidUUIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.GET("/chats/:chat_id/user-roles/:user_id", handler.GetUserRoleInChat)

	// invalid chat
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/invalid/user-roles/"+uuid.New().String(), nil)
	req.Header.Set("X-User-ID", uuid.New().String())
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// invalid user
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/chats/"+uuid.New().String()+"/user-roles/invalid", nil)
	req2.Header.Set("X-User-ID", uuid.New().String())
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)

	// invalid requester
	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest("GET", "/chats/"+uuid.New().String()+"/user-roles/"+uuid.New().String(), nil)
	req3.Header.Set("X-User-ID", "invalid")
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusBadRequest, w3.Code)

	mockController.AssertNotCalled(t, "GetUserRoleInChat", mock.Anything, mock.Anything, mock.Anything)
}

func TestChatHandler_GetUserRoleInChat_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()
	requesterID := uuid.New()

	mockController.On("GetUserRoleInChat", chatID, userID, requesterID).Return("", custom_errors.ErrUnauthorizedChat)

	router := gin.New()
	router.GET("/chats/:chat_id/user-roles/:user_id", handler.GetUserRoleInChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/"+chatID.String()+"/user-roles/"+userID.String(), nil)
	req.Header.Set("X-User-ID", requesterID.String())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_GetUserRoleInChat_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()
	requesterID := uuid.New()

	mockController.On("GetUserRoleInChat", chatID, userID, requesterID).Return("", custom_errors.ErrUserNotInChat)

	router := gin.New()
	router.GET("/chats/:chat_id/user-roles/:user_id", handler.GetUserRoleInChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/"+chatID.String()+"/user-roles/"+userID.String(), nil)
	req.Header.Set("X-User-ID", requesterID.String())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_GetMyRoleInChat_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()
	role := createTestChatRoleWithPermissions()

	mockController.On("GetMyRoleWithPermissions", chatID, userID).Return(role, nil)

	router := gin.New()
	router.GET("/chats/:chat_id/me/role", handler.GetMyRoleInChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/"+chatID.String()+"/me/role", nil)
	req.Header.Set("X-User-ID", userID.String())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.UserRoleWithPermissionsResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, role.Name, resp.RoleName)
	assert.Len(t, resp.Permissions, 1)

	mockController.AssertExpectations(t)
}

func TestChatHandler_GetMyRoleInChat_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.GET("/chats/:chat_id/me/role", handler.GetMyRoleInChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/invalid/me/role", nil)
	req.Header.Set("X-User-ID", uuid.New().String())
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/chats/"+uuid.New().String()+"/me/role", nil)
	req2.Header.Set("X-User-ID", "invalid")
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)

	mockController.AssertNotCalled(t, "GetMyRoleWithPermissions", mock.Anything, mock.Anything)
}

func TestChatHandler_GetMyRoleInChat_NotInChat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()

	mockController.On("GetMyRoleWithPermissions", chatID, userID).Return(nil, custom_errors.ErrUserNotInChat)

	router := gin.New()
	router.GET("/chats/:chat_id/me/role", handler.GetMyRoleInChat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/"+chatID.String()+"/me/role", nil)
	req.Header.Set("X-User-ID", userID.String())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_GetChatByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	resp := createTestChatResponse()
	resp.ID = chatID

	mockController.On("GetChatByID", chatID).Return(&resp, nil)

	router := gin.New()
	router.GET("/chats/:chat_id", handler.GetChatByID)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/"+chatID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_GetChatByID_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.GET("/chats/:chat_id", handler.GetChatByID)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "GetChatByID", mock.Anything)
}

func TestChatHandler_GetChatByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	mockController.On("GetChatByID", chatID).Return(nil, custom_errors.ErrChatNotFound)

	router := gin.New()
	router.GET("/chats/:chat_id", handler.GetChatByID)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/"+chatID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_GetChatByID_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	mockController.On("GetChatByID", chatID).Return(nil, errors.New("db"))

	router := gin.New()
	router.GET("/chats/:chat_id", handler.GetChatByID)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/"+chatID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_GetChatMembers_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	chatUsers := []models.ChatUser{
		{UserID: uuid.New(), RoleID: 1, Role: models.ChatRole{Name: "main"}},
	}

	mockController.On("GetChatMembers", chatID).Return(chatUsers, nil)

	router := gin.New()
	router.GET("/chats/:chat_id/members", handler.GetChatMembers)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/"+chatID.String()+"/members", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []dto.ChatMemberResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp, 1)

	mockController.AssertExpectations(t)
}

func TestChatHandler_GetChatMembers_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.GET("/chats/:chat_id/members", handler.GetChatMembers)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/invalid/members", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertNotCalled(t, "GetChatMembers", mock.Anything)
}

func TestChatHandler_GetChatMembers_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	mockController.On("GetChatMembers", chatID).Return(nil, errors.New("db"))

	router := gin.New()
	router.GET("/chats/:chat_id/members", handler.GetChatMembers)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/"+chatID.String()+"/members", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockController.AssertExpectations(t)
}
