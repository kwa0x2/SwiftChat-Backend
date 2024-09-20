package service

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/repository"
	"gorm.io/gorm"
)

type MessageService struct {
	MessageRepository *repository.MessageRepository
	RoomService       *RoomService
}

func (s *MessageService) CreateMessage(tx *gorm.DB, message *models.Message) (*models.Message, error) {
	return s.MessageRepository.CreateMessage(tx, message)
}

func (s *MessageService) GetPrivateConversation(senderId, receiverId string) ([]*models.Message, error) {
	return s.MessageRepository.GetPrivateConversation(senderId, receiverId)
}

func (s *MessageService) InsertAndUpdateRoom(message *models.Message) (*models.Message, error) {
	tx := s.MessageRepository.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	addedMessage, err := s.CreateMessage(tx, message)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	roomObj := &models.Room{
		RoomID:      message.RoomID,
		LastMessage: message.Message,
	}

	if err := s.RoomService.Update(tx, roomObj); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return addedMessage, nil
}

func (s *MessageService) GetMessageHistoryByRoomID(roomId string) ([]*models.Message, error) {
	return s.MessageRepository.GetMessageHistoryByRoomID(roomId)
}

func (s *MessageService) DeleteById(messageId string) error {
	return s.MessageRepository.DeleteById(messageId)
}

func (s *MessageService) UpdateMessageByIdBody(messageId, message string) error {
	return s.MessageRepository.UpdateMessageByIdBody(messageId, message)
}
