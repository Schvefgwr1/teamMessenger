package repositories

import (
	"chatService/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatUserRepository interface {
	AddUserToChat(chatUser *models.ChatUser) error
	ChangeUserRole(chatID, userID uuid.UUID, roleID int) error
	GetUserRole(chatID, userID uuid.UUID) (*models.ChatRole, error)
	GetChatUserWithRoleAndPermissions(userID, chatID uuid.UUID) (*models.ChatUser, error)
	GetChatUser(userID, chatID uuid.UUID) (*models.ChatUser, error)
	GetChatUsers(chatID uuid.UUID) ([]models.ChatUser, error)
	RemoveUserFromChat(chatID, userID uuid.UUID) error
	DeleteChatUsersByChatID(chatID uuid.UUID) error
}

type chatUserRepository struct {
	db *gorm.DB
}

func NewChatUserRepository(db *gorm.DB) ChatUserRepository {
	return &chatUserRepository{db}
}

func (r *chatUserRepository) AddUserToChat(chatUser *models.ChatUser) error {
	return r.db.Create(chatUser).Error
}

func (r *chatUserRepository) ChangeUserRole(chatID, userID uuid.UUID, roleID int) error {
	return r.db.Model(&models.ChatUser{}).
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Update("role_id", roleID).Error
}

func (r *chatUserRepository) GetUserRole(chatID, userID uuid.UUID) (*models.ChatRole, error) {
	var chatUser models.ChatUser
	err := r.db.Preload("Role").First(&chatUser, "chat_id = ? AND user_id = ?", chatID, userID).Error
	if err != nil {
		return nil, err
	}
	return &chatUser.Role, nil
}

func (r *chatUserRepository) GetChatUserWithRoleAndPermissions(userID, chatID uuid.UUID) (*models.ChatUser, error) {
	var chatUser models.ChatUser
	err := r.db.Preload("Role.Permissions").
		Where("user_id = ? AND chat_id = ?", userID, chatID).
		First(&chatUser).Error
	if err != nil {
		return nil, err
	}
	return &chatUser, nil
}

func (r *chatUserRepository) GetChatUser(userID, chatID uuid.UUID) (*models.ChatUser, error) {
	var chatUser models.ChatUser
	err := r.db.
		Where("user_id = ? AND chat_id = ?", userID, chatID).
		First(&chatUser).Error
	if err != nil {
		return nil, err
	}

	return &chatUser, nil
}

func (r *chatUserRepository) RemoveUserFromChat(chatID, userID uuid.UUID) error {
	return r.db.Where("chat_id = ? AND user_id = ?", chatID, userID).
		Delete(&models.ChatUser{}).Error
}

func (r *chatUserRepository) GetChatUsers(chatID uuid.UUID) ([]models.ChatUser, error) {
	var chatUsers []models.ChatUser
	err := r.db.Preload("Role").
		Where("chat_id = ?", chatID).
		Find(&chatUsers).Error
	if err != nil {
		return nil, err
	}
	return chatUsers, nil
}

func (r *chatUserRepository) DeleteChatUsersByChatID(chatID uuid.UUID) error {
	return r.db.Where("chat_id = ?", chatID).Delete(&models.ChatUser{}).Error
}
