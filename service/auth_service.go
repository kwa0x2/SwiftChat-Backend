package service

import "github.com/kwa0x2/realtime-chat-backend/repository"

type AuthService struct {
	AuthRepository *repository.AuthRepository
}

func (s *AuthService) IsIdUnique(id string) bool {
	return s.AuthRepository.IsIdUnique(id)
}
