package service

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/repository"
	"gorm.io/gorm"
)

type FriendService struct {
	FriendRepository *repository.FriendRepository
}

// region INSERT NEW FRIEND SERVICE
func (s *FriendService) Insert(tx *gorm.DB, friend *models.Friend) error {
	return s.FriendRepository.Insert(tx, friend)
}

//endregion

// region DELETE FRIEND BY MAIL SERVICE
func (s *FriendService) Delete(friend *models.Friend) error {
	return s.FriendRepository.Delete(friend)
}

//endregion

// region GET FRIENDS BY MAIL SERVICE
func (s *FriendService) GetFriends(userMail string) ([]*models.Friend, error) {
	return s.FriendRepository.GetFriends(userMail)
}

//endregion

// region GET BLOCKED FRIENDS BY MAIL SERVICE
func (s *FriendService) GetBlocked(userId string) ([]*models.Friend, error) {
	return s.FriendRepository.GetBlocked(userId)
}

//endregion

// region BLOCK FRIEND BY MAIL SERVICE
func (s *FriendService) Block(friend *models.Friend) error {
	return s.FriendRepository.Block(friend)
}

//endregion

func (s *FriendService) IsBlocked(userMail, otherUserMail string) (bool, error) {
	return s.FriendRepository.IsBlocked(userMail, otherUserMail)
}
