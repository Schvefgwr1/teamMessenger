package dto

import (
	ac "common/contracts/api-chat"
	"mime/multipart"

	"github.com/google/uuid"
)

type CreateChatRequestGateway struct {
	Name        string                `form:"name" binding:"required"`
	Description *string               `form:"description"`
	OwnerID     string                `form:"ownerID" binding:"required"`
	UserIDs     []string              `form:"userIDs" binding:"required"`
	Avatar      *multipart.FileHeader `form:"avatar"`
}

// ParseUUIDs преобразует строковые идентификаторы в UUID
func (r *CreateChatRequestGateway) ParseUUIDs() (ownerID uuid.UUID, userIDs []uuid.UUID, err error) {
	// Парсим ownerID
	ownerID, err = uuid.Parse(r.OwnerID)
	if err != nil {
		return uuid.Nil, nil, err
	}

	// Преобразуем каждый string в UUID
	userIDs = make([]uuid.UUID, len(r.UserIDs))
	for i, idStr := range r.UserIDs {
		userIDs[i], err = uuid.Parse(idStr)
		if err != nil {
			return uuid.Nil, nil, err
		}
	}

	return ownerID, userIDs, nil
}

type SendMessageRequestGateway struct {
	Content string                  `form:"content" binding:"required"`
	Files   []*multipart.FileHeader `form:"files"`
}

type CreateChatResponse struct {
	ID           uuid.UUID   `json:"id"`
	Name         string      `json:"name"`
	Description  *string     `json:"description"`
	OwnerID      uuid.UUID   `json:"ownerID"`
	UserIDs      []uuid.UUID `json:"userIDs"`
	AvatarFileID *int        `json:"avatarFileID,omitempty"`
}

type UpdateChatRequestGateway struct {
	Name          *string               `form:"name"`
	Description   *string               `form:"description"`
	Avatar        *multipart.FileHeader `form:"avatar"`
	AddUserIDs    []string              `form:"addUserIDs"`
	RemoveUserIDs []string              `form:"removeUserIDs"`
}

// ToUpdateChatRequest преобразует Gateway DTO в контракт
func (r *UpdateChatRequestGateway) ToUpdateChatRequest() (*ac.UpdateChatRequest, error) {
	req := &ac.UpdateChatRequest{
		Name:        r.Name,
		Description: r.Description,
	}

	// Парсим AddUserIDs
	if len(r.AddUserIDs) > 0 {
		req.AddUserIDs = make([]uuid.UUID, len(r.AddUserIDs))
		for i, idStr := range r.AddUserIDs {
			id, err := uuid.Parse(idStr)
			if err != nil {
				return nil, err
			}
			req.AddUserIDs[i] = id
		}
	}

	// Парсим RemoveUserIDs
	if len(r.RemoveUserIDs) > 0 {
		req.RemoveUserIDs = make([]uuid.UUID, len(r.RemoveUserIDs))
		for i, idStr := range r.RemoveUserIDs {
			id, err := uuid.Parse(idStr)
			if err != nil {
				return nil, err
			}
			req.RemoveUserIDs[i] = id
		}
	}

	return req, nil
}

type ChangeRoleRequestGateway struct {
	UserID string `json:"user_id" binding:"required"`
	RoleID int    `json:"role_id" binding:"required"`
}

// ToChangeRoleRequest преобразует Gateway DTO в контракт
func (r *ChangeRoleRequestGateway) ToChangeRoleRequest() (*ac.ChangeRoleRequest, error) {
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return nil, err
	}

	return &ac.ChangeRoleRequest{
		UserID: userID,
		RoleID: r.RoleID,
	}, nil
}

// MyRoleResponseGateway - ответ с ролью текущего пользователя и его permissions в чате
// ChatPermissionResponseGateway определён в chat_role_permission_dto.go
type MyRoleResponseGateway struct {
	RoleID      int                             `json:"roleId"`
	RoleName    string                          `json:"roleName"`
	Permissions []ChatPermissionResponseGateway `json:"permissions"`
}

// ChatMemberResponseGateway - участник чата для swagger
type ChatMemberResponseGateway struct {
	UserID   string `json:"userId"`
	RoleID   int    `json:"roleId"`
	RoleName string `json:"roleName"`
}
