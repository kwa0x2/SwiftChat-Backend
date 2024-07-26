package service

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/repository"
)

type FriendService struct {
	FriendRepository *repository.FriendRepository
}

func (s *FriendService) Insert(friend *models.Friend) error {
	return s.FriendRepository.Insert(friend)
}

func (s *FriendService) Delete(friend *models.Friend) error {
	return s.FriendRepository.Delete(friend)
}

func (s *FriendService) GetFriends(userMail string) ([]*models.Friend, error) {
	return s.FriendRepository.GetFriends(userMail)
}

func (s *FriendService) GetBlockeds(userId string) ([]*models.Friend, error) {
	return s.FriendRepository.GetBlockeds(userId)
}

func (s *FriendService) Block(friend *models.Friend) error {
	return s.FriendRepository.Block(friend)
}
