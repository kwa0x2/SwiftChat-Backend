package repository

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"gorm.io/gorm"
)

type MessageRepository struct {
	DB *gorm.DB
}

func (r *MessageRepository) CreateMessage(tx *gorm.DB, message *models.Message) (*models.Message, error) {
	db := r.DB
	if tx != nil {
		db = tx
	}

	if err := db.Create(&message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

func (r *MessageRepository) GetPrivateConversation(senderId, receiverId string) ([]*models.Message, error) {
	var messages []*models.Message
	if err := r.DB.Where(
		"(message_sender_id = ? AND message_receiver_id = ?) OR (message_receiver_id = ? AND message_sender_id = ?)",
		senderId, receiverId, senderId, receiverId,
	).Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MessageRepository) GetMessageHistoryByRoomID(roomId string) ([]*models.Message, error) {
	var messages []*models.Message
	if err := r.DB.Unscoped().
		Select(`
			message_id, 
			sender_id, 
			room_id, 
			message_status, 
			"createdAt", 
			"updatedAt",
			"deletedAt",
			CASE WHEN "deletedAt" IS NOT NULL THEN '' ELSE message END as message
		`).
		Where("room_id = ?", roomId).
		Order(`"createdAt" ASC`).
		Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MessageRepository) DeleteById(messageId string) error {
	return r.DB.Delete(&models.Message{}, "message_id = ?", messageId).Error
}

func (r *MessageRepository) UpdateMessageByIdBody(messageId, message string) error {
	return r.DB.Model(&models.Message{}).Where("message_id = ?", messageId).Update("message", message).Error
}
