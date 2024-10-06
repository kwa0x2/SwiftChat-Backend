package repository

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/types"
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
			message_read_status,
			message_type,
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

func (r *MessageRepository) StarMessageById(messageId string) error {
	return r.DB.Model(&models.Message{}).Where("message_id = ?", messageId).UpdateColumns(&models.Message{MessageType: types.StarredText}).Error
}

func (r *MessageRepository) ReadMessageByRoomId(connectedUserID, roomId string, messageId *string) error {
	query := r.DB.Model(&models.Message{}).Unscoped().Where("sender_id != ? AND room_id = ?", connectedUserID, roomId)

	if messageId != nil {
		query = query.Where("message_id = ?", *messageId)
	}

	return query.UpdateColumns(&models.Message{MessageReadStatus: types.Readed}).Error
}
