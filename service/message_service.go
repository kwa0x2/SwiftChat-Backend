package service

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/repository"
)

type MessageService struct {
	MessageRepository *repository.MessageRepository
}

func (s *MessageService) InsertMessage(message *models.Message) (*models.Message, error) {
	return s.MessageRepository.InsertMessage(message)
}