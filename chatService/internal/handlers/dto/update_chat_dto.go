package dto

import "github.com/google/uuid"

type UpdateChatDTO struct {
	Name          *string     `json:"name,omitempty"`
	Description   *string     `json:"description,omitempty"`
	AvatarFileID  *int        `json:"avatarFileID,omitempty"`
	AddUserIDs    []uuid.UUID `json:"addUserIDs,omitempty"`
	RemoveUserIDs []uuid.UUID `json:"removeUserIDs,omitempty"`
}
