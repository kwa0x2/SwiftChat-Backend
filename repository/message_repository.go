package repository

import (
	"github.com/google/uuid"
	"github.com/kwa0x2/swiftchat-backend/models"
	"github.com/kwa0x2/swiftchat-backend/types"
	"gorm.io/gorm"
)

type IMessageRepository interface {
	Create(tx *gorm.DB, message *models.Message) (*models.Message, error)
	UpdateExceptUpdatedAt(whereMessage *models.Message, updateMessage *models.Message, isUnscoped bool) error
	Delete(whereMessage *models.Message) error
	ReadMessageByRoomId(connectedUserID, roomId string, messageId *string) error
	GetMessageHistoryByRoomID(roomId uuid.UUID) ([]*models.Message, error)
	GetDB() *gorm.DB
}

type messageRepository struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) IMessageRepository {
	return &messageRepository{
		DB: db,
	}
}

// region "Create" adds a new message to the database
func (r *messageRepository) Create(tx *gorm.DB, message *models.Message) (*models.Message, error) {
	db := r.DB
	if tx != nil {
		db = tx // Use the provided transaction if available
	}

	if err := db.Create(&message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

// endregion

// region "Update" modifies the fields of a message in the database based on specified conditions
func (r *messageRepository) UpdateExceptUpdatedAt(whereMessage *models.Message, updateMessage *models.Message, isUnscoped bool) error {
	query := r.DB.Model(&models.Message{}).Where(whereMessage)

	if isUnscoped {
		query = query.Unscoped() // Include soft-deleted messages in the update
	}

	return query.UpdateColumns(updateMessage).Error
}

// endregion

// region "Delete" removes a message from the database
func (r *messageRepository) Delete(whereMessage *models.Message) error {
	return r.DB.Delete(whereMessage).Error
}

// endregion

// region "ReadMessageByRoomId" marks a message as read for a specific user and room
func (r *messageRepository) ReadMessageByRoomId(connectedUserID, roomId string, messageId *string) error {
	query := r.DB.Model(&models.Message{}).Unscoped().Where("sender_id != ? AND room_id = ?", connectedUserID, roomId)

	if messageId != nil {
		query = query.Where("message_id = ?", *messageId) // Filter by message ID if provided
	}

	return query.UpdateColumns(&models.Message{MessageReadStatus: types.Readed}).Error
}

// endregion

// region "GetMessageHistoryByRoomID" retrieves the message history for a specific room
func (r *messageRepository) GetMessageHistoryByRoomID(roomId uuid.UUID) ([]*models.Message, error) {
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
		Where(&models.Message{RoomID: roomId}).
		Order(`"createdAt" ASC`).
		Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

// endregion

// region "GetDB" returns the underlying gorm.DB instance
func (r *messageRepository) GetDB() *gorm.DB {
	return r.DB // Return the database instance
}

// endregion
