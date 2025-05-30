package repositories

import (
	"chatService/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository interface {
	CreateMessage(message *models.Message) error
	CreateMessageFile(msgFile *models.MessageFile) error
	GetMessageWithFile(msgID uuid.UUID) (*models.Message, error)
	GetChatMessages(chatID uuid.UUID, offset, limit int) ([]models.Message, error)
	SearchMessages(userID, chatID uuid.UUID, text string, limit, offset int) ([]models.Message, int64, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db}
}

func (r *messageRepository) CreateMessage(message *models.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepository) CreateMessageFile(msgFile *models.MessageFile) error {
	return r.db.Create(msgFile).Error
}

func (r *messageRepository) GetMessageWithFile(msgID uuid.UUID) (*models.Message, error) {
	var message models.Message
	err := r.db.Preload("Files").
		Where("id = ?", msgID).
		First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *messageRepository) GetChatMessages(chatID uuid.UUID, offset, limit int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Where("chat_id = ?", chatID).
		Offset(offset).Limit(limit).
		Preload("Files").
		Order("created_at desc").
		Find(&messages).Error
	return messages, err
}

func (r *messageRepository) SearchMessages(userID, chatID uuid.UUID, text string, limit, offset int) ([]models.Message, int64, error) {
	var messages []models.Message
	var total int64

	query := r.db.Model(&models.Message{}).
		Joins("JOIN chat_service.chats ch ON ch.id = messages.chat_id").
		Joins("JOIN chat_service.chat_user cu ON cu.chat_id = ch.id").
		Where("cu.user_id = ? AND ch.id = ?", userID, chatID).
		Where("messages.content ILIKE ?", "%"+text+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("messages.created_at DESC").
		Limit(limit).Offset(offset).
		Find(&messages).Error

	return messages, total, err
}
