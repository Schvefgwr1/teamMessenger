package repositories

import (
	"chatService/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatRepository interface {
	GetUserChats(userID uuid.UUID) ([]models.Chat, error)
	CreateChat(chat *models.Chat) error
	UpdateChat(chat *models.Chat) error
	DeleteChat(chatID uuid.UUID) error
	GetChatByID(chatID uuid.UUID) (*models.Chat, error)
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db}
}

func (r *chatRepository) GetUserChats(userID uuid.UUID) ([]models.Chat, error) {
	var chats []models.Chat
	err := r.db.Joins("JOIN chat_service.chat_user cu ON cu.chat_id = chats.id").
		Where("cu.user_id = ?", userID).
		Find(&chats).Error
	return chats, err
}

func (r *chatRepository) CreateChat(chat *models.Chat) error {
	return r.db.Create(chat).Error
}

func (r *chatRepository) UpdateChat(chat *models.Chat) error {
	return r.db.Save(chat).Error
}

func (r *chatRepository) DeleteChat(chatID uuid.UUID) error {
	return r.db.Delete(&models.Chat{}, chatID).Error
}

func (r *chatRepository) GetChatByID(chatID uuid.UUID) (*models.Chat, error) {
	var chat models.Chat
	// Явно исключаем загрузку связанных данных (Users и Messages)
	err := r.db.Omit("Users", "Messages").First(&chat, "id = ?", chatID).Error
	return &chat, err
}
