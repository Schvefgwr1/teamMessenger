package handlers

import (
	"apiService/internal/dto"
	"apiService/internal/handlers"
	"bytes"
	ac "common/contracts/api-chat"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Тесты для ChatHandler.CreateChat

func TestChatHandler_CreateChat_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	ownerID := uuid.New()
	userID1 := uuid.New()
	userID2 := uuid.New()
	chatID := uuid.New()

	expectedChat := &dto.CreateChatResponse{
		ID:   chatID,
		Name: "Test Chat",
	}

	mockController.On("CreateChat", mock.Anything, ownerID, mock.Anything).Return(expectedChat, nil)

	router := gin.New()
	router.POST("/chats", handler.CreateChat)

	// Act - создаем multipart форму
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "Test Chat")
	writer.WriteField("description", "Test Description")
	writer.WriteField("ownerID", ownerID.String())
	writer.WriteField("userIDs", userID1.String())
	writer.WriteField("userIDs", userID2.String())
	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	mockController.AssertExpectations(t)
}

func TestChatHandler_CreateChat_InvalidUUID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.POST("/chats", handler.CreateChat)

	// Act
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "Test Chat")
	writer.WriteField("ownerID", "invalid-uuid")
	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chats", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "CreateChat", mock.Anything, mock.Anything, mock.Anything)
}

// Тесты для ChatHandler.BanUser

func TestChatHandler_BanUser_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()
	ownerID := uuid.New()

	mockController.On("BanUser", chatID, userID, ownerID).Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", ownerID)
		c.Next()
	})
	router.PATCH("/chats/:chat_id/ban/:user_id", handler.BanUser)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/chats/"+chatID.String()+"/ban/"+userID.String(), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "user banned successfully", response["message"])

	mockController.AssertExpectations(t)
}

func TestChatHandler_BanUser_InvalidChatID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.PATCH("/chats/:chat_id/ban/:user_id", handler.BanUser)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/chats/invalid/ban/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "BanUser", mock.Anything, mock.Anything, mock.Anything)
}

func TestChatHandler_BanUser_MissingUserID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()

	router := gin.New()
	router.PATCH("/chats/:chat_id/ban/:user_id", handler.BanUser)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/chats/"+chatID.String()+"/ban/invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "BanUser", mock.Anything, mock.Anything, mock.Anything)
}

// Тесты для ChatHandler.ChangeUserRole

func TestChatHandler_ChangeUserRole_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	ownerID := uuid.New()
	userID := uuid.New()

	mockController.On("ChangeUserRole", chatID, ownerID, mock.Anything).Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", ownerID)
		c.Next()
	})
	router.PATCH("/chats/:chat_id/roles/change", handler.ChangeUserRole)

	// Act
	reqBody := `{"user_id":"` + userID.String() + `","role_id":2}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/chats/"+chatID.String()+"/roles/change", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "user role changed successfully", response["message"])

	mockController.AssertExpectations(t)
}

// Тесты для ChatHandler.GetMyRoleInChat

func TestChatHandler_GetMyRoleInChat_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()
	expectedRole := &ac.MyRoleResponse{
		RoleID:   1,
		RoleName: "admin",
	}

	mockController.On("GetMyRoleInChat", chatID, userID).Return(expectedRole, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.GET("/chats/:chat_id/me/role", handler.GetMyRoleInChat)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/"+chatID.String()+"/me/role", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	mockController.AssertExpectations(t)
}

func TestChatHandler_GetMyRoleInChat_NotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	chatID := uuid.New()
	userID := uuid.New()
	notFoundError := errors.New("user not found in chat")

	mockController.On("GetMyRoleInChat", chatID, userID).Return(nil, notFoundError)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.GET("/chats/:chat_id/me/role", handler.GetMyRoleInChat)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chats/"+chatID.String()+"/me/role", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	mockController.AssertExpectations(t)
}
