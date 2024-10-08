package service

import (
	"github.com/google/uuid"
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/repository"
	"github.com/kwa0x2/realtime-chat-backend/types"
	"gorm.io/gorm"
)

type IMessageService interface {
	Create(tx *gorm.DB, message *models.Message) (*models.Message, error)
	InsertAndUpdateRoom(message *models.Message) (*models.Message, error)
	GetMessageHistoryByRoomID(roomId uuid.UUID) ([]*models.Message, error)
	DeleteById(messageId uuid.UUID) error
	UpdateMessageById(messageId uuid.UUID, message string) error
	StarMessageById(messageId uuid.UUID) error
	ReadMessageByRoomId(connectedUserID, roomId string, messageId *string) error
}

type messageService struct {
	MessageRepository repository.IMessageRepository
	RoomService       IRoomService
}

func NewMessageService(messageRepo repository.IMessageRepository, roomService IRoomService) IMessageService {
	return &messageService{
		MessageRepository: messageRepo,
		RoomService:       roomService,
	}
}

// region "Create" adds a new message to the database
func (s *messageService) Create(tx *gorm.DB, message *models.Message) (*models.Message, error) {
	return s.MessageRepository.Create(tx, message)
}

// endregion

// region "InsertAndUpdateRoom" creates a new message and updates the corresponding room
func (s *messageService) InsertAndUpdateRoom(message *models.Message) (*models.Message, error) {
	// Start a new database transaction.
	tx := s.MessageRepository.GetDB().Begin()
	if tx.Error != nil {
		// If starting the transaction failed, return the error.
		return nil, tx.Error
	}

	// Create a new message and check for errors.
	addedMessage, err := s.Create(tx, message)
	if err != nil {
		// Rollback the transaction in case of an error.
		tx.Rollback()
		return nil, err
	}

	// Prepare the room data for updating.
	whereRoom := &models.Room{
		RoomID: message.RoomID, // Specify the room to update using the room ID from the message.
	}

	updateRoom := &models.Room{
		LastMessageID: message.MessageID, // Set the last message ID to the new message's ID.
		LastMessage:   message.Message,   // Update the last message text.
	}

	// Update the room with the new last message details.
	if updateErr := s.RoomService.Update(tx, whereRoom, updateRoom); updateErr != nil {
		// Rollback the transaction if the update fails.
		tx.Rollback()
		return nil, updateErr
	}

	// Commit the transaction if everything went smoothly.
	if commitErr := tx.Commit().Error; commitErr != nil {
		return nil, commitErr // Return the error if committing the transaction fails.
	}

	// Return the added message.
	return addedMessage, nil
}

// endregion

// region "GetMessageHistoryByRoomID" retrieves the message history for a specific room
func (s *messageService) GetMessageHistoryByRoomID(roomId uuid.UUID) ([]*models.Message, error) {
	return s.MessageRepository.GetMessageHistoryByRoomID(roomId)
}

// endregion

// region "DeleteById" removes a message by its ID from the database
func (s *messageService) DeleteById(messageId uuid.UUID) error {
	// Prepare the message data for deletion.
	whereMessage := &models.Message{
		MessageID: messageId, // Specify the message to delete using its ID.
	}

	return s.MessageRepository.Delete(whereMessage)
}

// endregion

// region "UpdateMessageById" updates the content of a message identified by its ID
func (s *messageService) UpdateMessageById(messageId uuid.UUID, message string) error {
	// Prepare the message data for updating.
	whereMessage := &models.Message{
		MessageID: messageId, // Specify the message to update using its ID.
	}

	updateMessage := &models.Message{
		Message: message, // Set the new message content.
	}

	return s.MessageRepository.UpdateExceptUpdatedAt(whereMessage, updateMessage, false)
}

// endregion

// region "StarMessageById" marks a message as starred by its ID
func (s *messageService) StarMessageById(messageId uuid.UUID) error {
	// Prepare the message data for starring.
	whereMessage := &models.Message{
		MessageID: messageId, // Specify the message to star using its ID.
	}

	updateMessage := &models.Message{
		MessageType: types.StarredText, // Set the message type to "starred".
	}

	return s.MessageRepository.UpdateExceptUpdatedAt(whereMessage, updateMessage, false)
}

// endregion

// region "ReadMessageByRoomId" marks a message as read for a specific user and room
func (s *messageService) ReadMessageByRoomId(connectedUserID, roomId string, messageId *string) error {
	return s.MessageRepository.ReadMessageByRoomId(connectedUserID, roomId, messageId)
}

// endregion
