package service

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/repository"
	"github.com/kwa0x2/realtime-chat-backend/types"
	"gorm.io/gorm"
)

type IFriendService interface {
	Create(tx *gorm.DB, friend *models.Friend) error
	Update(tx *gorm.DB, whereFriend *models.Friend, updates *models.Friend) error
	UpdateFriendStatusByMail(tx *gorm.DB, userEmail, userEmail2 string, friendStatus types.FriendStatus) error
	Delete(UserEmail, UserEmail2 string) error
	GetFriends(userEmail string, isUnFriendStatusAllow bool) ([]*models.Friend, error)
	GetSpecificFriend(userEmail, userEmail2 string) (*models.Friend, error)
	GetBlockedUsers(userEmail string) ([]*models.Friend, error)
	Block(userEmail, userEmail2 string) (string, error)
	IsBlocked(userMail, otherUserMail string) (bool, error)
}

type friendService struct {
	FriendRepository repository.IFriendRepository
}

func NewFriendService(friendRepository repository.IFriendRepository) IFriendService {
	return &friendService{
		FriendRepository: friendRepository,
	}
}

// region "Create" adds a new friend to the database. If a friendship already exists, it updates the status
func (s *friendService) Create(tx *gorm.DB, friend *models.Friend) error {
	// Check if a specific friend relationship exists between the two users using the provided email addresses.
	existingFriend, err := s.GetSpecificFriend(friend.UserMail, friend.UserMail2)
	if err != nil {
		return err
	}

	// If an existing friendship is found and the second user's email matches the provided friend's email,
	if existingFriend != nil && existingFriend.UserMail == friend.UserMail2 {
		// Update the status of the friendship to "Friend" (or a defined status) for the two users.
		if updateErr := s.UpdateFriendStatusByMail(nil, friend.UserMail, existingFriend.UserMail, types.Friend); updateErr != nil {
			return updateErr
		}
		return nil
	}

	// If no existing friendship was found, create a new friend record in the repository.
	return s.FriendRepository.Create(tx, friend)
}

// endregion

// region "Update" modifies the fields of a friend in the database based on specified conditions
func (s *friendService) Update(tx *gorm.DB, whereFriend *models.Friend, updates *models.Friend) error {
	return s.FriendRepository.Update(tx, whereFriend, updates)
}

// endregion

// region "UpdateFriendStatusByMail" updates the deletedAt field and friendStatus for given user emails
func (s *friendService) UpdateFriendStatusByMail(tx *gorm.DB, userEmail, userEmail2 string, friendStatus types.FriendStatus) error {
	return s.FriendRepository.UpdateFriendStatusByMail(tx, userEmail, userEmail2, friendStatus)
}

// endregion

// region "Delete" removes a friendship between two users
func (s *friendService) Delete(UserEmail, UserEmail2 string) error {
	return s.FriendRepository.Delete(UserEmail, UserEmail2)
}

// endregion

// region "GetFriends" retrieves a list of friends for a given user email
func (s *friendService) GetFriends(userEmail string, isUnFriendStatusAllow bool) ([]*models.Friend, error) {
	return s.FriendRepository.GetFriends(userEmail, isUnFriendStatusAllow)
}

// endregion

// region "GetSpecificFriend" retrieves a specific friend relationship between two users
func (s *friendService) GetSpecificFriend(userEmail, userEmail2 string) (*models.Friend, error) {
	return s.FriendRepository.GetSpecificFriend(userEmail, userEmail2)
}

// endregion

// region "GetBlockedUsers" retrieves a list of blocked users for a given email
func (s *friendService) GetBlockedUsers(userEmail string) ([]*models.Friend, error) {
	return s.FriendRepository.GetBlockedUsers(userEmail)
}

// endregion

// region "Block" updates the status of a friendship to blocked
func (s *friendService) Block(userEmail, userEmail2 string) (string, error) {
	return s.FriendRepository.Block(userEmail, userEmail2)
}

// endregion

// region "IsBlocked" checks if a user is blocked by another user
func (s *friendService) IsBlocked(userMail, otherUserMail string) (bool, error) {
	return s.FriendRepository.IsBlocked(userMail, otherUserMail)
}

// endregion
