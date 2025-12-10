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
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChatHandler_DeleteChat_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.DELETE("/chats/:chat_id", handler.DeleteChat)

	// DeleteChat не проверяет UUID, просто игнорирует ошибку и вызывает контроллер с нулевым UUID
	mockController.On("DeleteChat", uuid.Nil).Return(nil)

	req, _ := http.NewRequest("DELETE", "/chats/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_UpdateChat_InvalidCredentials_Additional(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.PUT("/chats/:chat_id", handler.UpdateChat)

	chatID := uuid.New()
	name := "Updated Name"
	updateDTO := dto.UpdateChatDTO{
		Name: &name,
	}

	mockController.On("UpdateChat", chatID, mock.MatchedBy(func(dto *dto.UpdateChatDTO) bool {
		return dto.Name != nil && *dto.Name == name
	})).Return(nil, custom_errors.ErrInvalidCredentials)

	body, _ := json.Marshal(updateDTO)
	req, _ := http.NewRequest("PUT", "/chats/"+chatID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_GetChatByID_ChatNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.GET("/chats/:chat_id", handler.GetChatByID)

	chatID := uuid.New()
	mockController.On("GetChatByID", chatID).Return(nil, custom_errors.ErrChatNotFound)

	req, _ := http.NewRequest("GET", "/chats/"+chatID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_GetChatMembers_AllErrorTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.GET("/chats/:chat_id/members", handler.GetChatMembers)

	chatID := uuid.New()

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"database error", custom_errors.NewDatabaseError("error"), http.StatusInternalServerError},
		{"unknown error", errors.New("unknown"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockController.On("GetChatMembers", chatID).Return(nil, tt.err)

			req, _ := http.NewRequest("GET", "/chats/"+chatID.String()+"/members", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestChatHandler_GetUserRoleInChat_UnauthorizedChat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.GET("/chats/:chat_id/users/:user_id/role", handler.GetUserRoleInChat)

	chatID := uuid.New()
	userID := uuid.New()
	requesterID := uuid.New()

	mockController.On("GetUserRoleInChat", chatID, userID, requesterID).Return("", custom_errors.ErrUnauthorizedChat)

	req, _ := http.NewRequest("GET", "/chats/"+chatID.String()+"/users/"+userID.String()+"/role", nil)
	req.Header.Set("X-User-ID", requesterID.String())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_GetMyRoleInChat_UserNotInChat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.GET("/chats/:chat_id/my-role", handler.GetMyRoleInChat)

	chatID := uuid.New()
	userID := uuid.New()

	mockController.On("GetMyRoleWithPermissions", chatID, userID).Return(nil, custom_errors.ErrUserNotInChat)

	req, _ := http.NewRequest("GET", "/chats/"+chatID.String()+"/my-role", nil)
	req.Header.Set("X-User-ID", userID.String())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockController.AssertExpectations(t)
}

func TestChatHandler_BanUser_AllErrorTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.POST("/chats/:chat_id/ban/:user_id", handler.BanUser)

	chatID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"invalid credentials", custom_errors.ErrInvalidCredentials, http.StatusInternalServerError},
		{"internal server error", custom_errors.ErrInternalServerError, http.StatusInternalServerError},
		{"unknown error", errors.New("unknown"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockController.On("BanUser", chatID, userID).Return(tt.err)

			req, _ := http.NewRequest("POST", "/chats/"+chatID.String()+"/ban/"+userID.String(), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestChatHandler_ChangeUserRole_AllErrorTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatController)
	handler := handlers.NewChatHandler(mockController)

	router := gin.New()
	router.PATCH("/chats/:chat_id/roles/change", handler.ChangeUserRole)

	chatID := uuid.New()
	request := struct {
		UserID uuid.UUID `json:"user_id"`
		RoleID int       `json:"role_id"`
	}{
		UserID: uuid.New(),
		RoleID: 1,
	}

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"invalid credentials", custom_errors.ErrInvalidCredentials, http.StatusInternalServerError},
		{"unknown error", errors.New("unknown"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockController.On("ChangeUserRole", chatID, request.UserID, request.RoleID).Return(tt.err)

			body, _ := json.Marshal(request)
			req, _ := http.NewRequest("PATCH", "/chats/"+chatID.String()+"/roles/change", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
