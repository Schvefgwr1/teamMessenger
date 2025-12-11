//go:build integration
// +build integration

package integration

import (
	"apiService/internal/dto"
	ac "common/contracts/api-chat"
	at "common/contracts/api-task"
	au "common/contracts/api-user"

	"github.com/google/uuid"
)

// NewTestRegisterRequest создает тестовый запрос на регистрацию
func NewTestRegisterRequest() *dto.RegisterUserRequestGateway {
	return &dto.RegisterUserRequestGateway{
		Username: "testuser_" + uuid.New().String()[:8],
		Email:    "test_" + uuid.New().String()[:8] + "@example.com",
		Password: "password123",
		Age:      25,
		RoleID:   1,
	}
}

// NewTestUpdateUserRequest создает тестовый запрос на обновление пользователя
func NewTestUpdateUserRequest() *dto.UpdateUserRequestGateway {
	username := "updated_user_" + uuid.New().String()[:8]
	age := 30
	return &dto.UpdateUserRequestGateway{
		Username: &username,
		Age:      &age,
	}
}

// NewTestCreateChatRequest создает тестовый запрос на создание чата
func NewTestCreateChatRequest() *dto.CreateChatRequestGateway {
	name := "test_chat_" + uuid.New().String()[:8]
	description := "Test chat description"
	ownerID := uuid.New()
	return &dto.CreateChatRequestGateway{
		Name:        name,
		Description: &description,
		OwnerID:     ownerID.String(),
		UserIDs:     []string{uuid.New().String(), uuid.New().String()},
	}
}

// NewTestUpdateChatRequest создает тестовый запрос на обновление чата
func NewTestUpdateChatRequest() *dto.UpdateChatRequestGateway {
	name := "updated_chat_" + uuid.New().String()[:8]
	return &dto.UpdateChatRequestGateway{
		Name:        &name,
		Description: stringPtr("Updated description"),
	}
}

// NewTestCreateTaskRequest создает тестовый запрос на создание задачи
func NewTestCreateTaskRequest() *dto.CreateTaskRequestGateway {
	title := "test_task_" + uuid.New().String()[:8]
	description := "Test task description"
	executorID := uuid.New().String()
	return &dto.CreateTaskRequestGateway{
		Title:       title,
		Description: &description,
		ExecutorID:  executorID,
	}
}

// NewTestChangeRoleRequest создает тестовый запрос на изменение роли
func NewTestChangeRoleRequest(userID uuid.UUID) *ac.ChangeRoleRequest {
	return &ac.ChangeRoleRequest{
		UserID: userID,
		RoleID: 2,
	}
}

// NewTestCreateRoleRequest создает тестовый запрос на создание роли
func NewTestCreateRoleRequest() *au.CreateRoleRequest {
	name := "test_role_" + uuid.New().String()[:8]
	description := "Test role description"
	return &au.CreateRoleRequest{
		Name:        name,
		Description: description,
	}
}

// NewTestCreateChatRoleRequest создает тестовый запрос на создание роли чата
func NewTestCreateChatRoleRequest() *ac.CreateRoleRequest {
	name := "test_chat_role_" + uuid.New().String()[:8]
	return &ac.CreateRoleRequest{
		Name: name,
	}
}

// NewTestCreateChatPermissionRequest создает тестовый запрос на создание разрешения чата
func NewTestCreateChatPermissionRequest() *ac.CreatePermissionRequest {
	name := "test_permission_" + uuid.New().String()[:8]
	return &ac.CreatePermissionRequest{
		Name: name,
	}
}

// NewTestCreateStatusRequest создает тестовый запрос на создание статуса задачи
func NewTestCreateStatusRequest() *at.CreateStatusRequest {
	name := "test_status_" + uuid.New().String()[:8]
	return &at.CreateStatusRequest{
		Name: name,
	}
}

// stringPtr возвращает указатель на строку
func stringPtr(s string) *string {
	return &s
}
