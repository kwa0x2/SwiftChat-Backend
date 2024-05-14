package service

import "github.com/kwa0x2/realtime-chat-backend/repository"

type ChatService struct {
	ChatRepository *repository.ChatRepository
}