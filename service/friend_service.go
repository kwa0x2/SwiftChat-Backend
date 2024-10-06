package service

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/repository"
	"github.com/kwa0x2/realtime-chat-backend/types"
	"gorm.io/gorm"
)

type FriendService struct {
	FriendRepository *repository.FriendRepository
}

// region INSERT NEW FRIEND SERVICE
func (s *FriendService) Insert(tx *gorm.DB, friend *models.Friend) error {
	existingFriend, err := s.GetFriend(friend.UserMail, friend.UserMail2)
	if err != nil {
		return err
	}

	if existingFriend != nil && existingFriend.UserMail == friend.UserMail2 {

		if err := s.UpdateDeletedAtByMail(nil, friend.UserMail, existingFriend.UserMail, types.Friend); err != nil {
			return err
		}

		return nil
	}

	return s.FriendRepository.Insert(tx, friend)
}

//endregion

func (s *FriendService) Update(tx *gorm.DB, filter map[string]interface{}, updates map[string]interface{}) error {

	return s.FriendRepository.Update(tx, filter, updates)
}

func (s *FriendService) UpdateDeletedAtByMail(tx *gorm.DB, userMail, userMail2 string, friendStatus types.FriendStatus) error {

	return s.FriendRepository.UpdateDeletedAtByMail(tx, userMail, userMail2, friendStatus)
}

// region DELETE FRIEND BY MAIL SERVICE
func (s *FriendService) Delete(friend *models.Friend) error {
	return s.FriendRepository.Delete(friend)
}

//endregion

// region GET FRIENDS BY MAIL SERVICE
func (s *FriendService) GetFriends(userMail string, isUnFriendStatusAllow bool) ([]*models.Friend, error) {
	return s.FriendRepository.GetFriends(userMail, isUnFriendStatusAllow)
}

//endregion

func (s *FriendService) GetFriend(userMail, userMail2 string) (*models.Friend, error) {
	return s.FriendRepository.GetFriend(userMail, userMail2)
}

// region GET BLOCKED FRIENDS BY MAIL SERVICE
func (s *FriendService) GetBlocked(userEmail string) ([]*models.Friend, error) {
	return s.FriendRepository.GetBlocked(userEmail)
}

//endregion

// region BLOCK FRIEND BY MAIL SERVICE
func (s *FriendService) Block(friend *models.Friend) (string, error) {
	return s.FriendRepository.Block(friend)
}

//endregion

func (s *FriendService) IsBlocked(userMail, otherUserMail string) (bool, error) {
	return s.FriendRepository.IsBlocked(userMail, otherUserMail)
}
