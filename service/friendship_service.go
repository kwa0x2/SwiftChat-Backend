package service

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/repository"
)

type FriendshipService struct {
	FriendshipRepository *repository.FriendshipRepository
}

func (s *FriendshipService) SendFriendRequest(friendship *models.Friendship) (string,error) {
	return s.FriendshipRepository.SendFriendRequest(friendship)
}

func (s *FriendshipService) GetComingRequests(receiverId string) ([]*models.Friendship,error){
	return s.FriendshipRepository.GetComingRequests(receiverId)
}

func (s *FriendshipService) GetFriends(userId string) ([]*models.Friendship,error) {
	return s.FriendshipRepository.GetFriends(userId)
}

func (s *FriendshipService) GetBlockeds(userId string) ([]*models.Friendship, error) {
	return s.FriendshipRepository.GetBlockeds(userId)
}

func (s *FriendshipService) Block(friendship *models.Friendship) (error){
	return s.FriendshipRepository.Block(friendship)
}

func (s *FriendshipService) Delete(friendship *models.Friendship) (error){
	return s.FriendshipRepository.Delete(friendship)
}

func (s *FriendshipService) Accept(friendship *models.Friendship) (error){
	return s.FriendshipRepository.Accept(friendship)
}
func (s *FriendshipService) Reject(friendship *models.Friendship) (error){
	return s.FriendshipRepository.Reject(friendship)
}