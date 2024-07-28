package repository

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"gorm.io/gorm"
)

type FriendRepository struct {
	DB *gorm.DB
}

// region INSERT NEW FRIEND REPOSITORY
func (r *FriendRepository) Insert(friend *models.Friend) error {
	if err := r.DB.Create(&friend).Error; err != nil {
		return err
	}
	return nil
}

//endregion

// region DELETE FRIEND BY MAIL REPOSITORY
func (r *FriendRepository) Delete(friend *models.Friend) error {
	if err := r.DB.Debug().
		Where("(user_mail = ? AND user_mail2 = ?) OR (user_mail = ? AND user_mail2 = ?)",
			friend.UserMail,
			friend.UserMail2,
			friend.UserMail2,
			friend.UserMail).
		Delete(&models.Friend{}).Error; err != nil {
		return err
	}
	return nil
}

//endregion

// region GET FRIENDS BY MAIL REPOSITORY
func (r *FriendRepository) GetFriends(userMail string) ([]*models.Friend, error) {
	var friends []*models.Friend

	if err := r.DB.Debug().
		Preload("User").
		Select("CASE WHEN user_mail = ? THEN user_mail2 ELSE user_mail END as user_mail", userMail).
		Where("user_mail = ? AND friend_status = ?", userMail, "friend").
		Or("user_mail2 = ? AND friend_status = ?", userMail, "friend").
		Find(&friends).Error; err != nil {
		return nil, err
	}

	return friends, nil
}

//endregion

// region GET BLOCKED FRIENDS BY MAIL SERVICE
func (r *FriendRepository) GetBlocked(userMail string) ([]*models.Friend, error) {
	var friends []*models.Friend

	if err := r.DB.Debug().
		Preload("User").
		Select("CASE WHEN user_mail = ? THEN user_mail2 ELSE user_mail END as user_mail", userMail).
		Where("user_mail = ? AND (friend_status = ? OR friend_status = ?)", userMail, "block_both", "block_first_second").
		Or("user_mail2 = ? AND (friend_status = ? OR friend_status = ?)", userMail, "block_both", "block_second_first").
		Find(&friends).Error; err != nil {
		return nil, err
	}

	return friends, nil
}

//endregion

// region BLOCK FRIEND BY MAIL SERVICE
func (r *FriendRepository) Block(friend *models.Friend) error {
	var existingFriend models.Friend

	if err := r.DB.Debug().
		Where("(user_mail = ? AND user_mail2 = ?) OR (user_mail = ? AND user_mail2 = ?)", friend.UserMail, friend.UserMail2, friend.UserMail2, friend.UserMail).
		First(&existingFriend).Error; err != nil {
		return err

	}

	switch existingFriend.FriendStatus {
	case "friend":
		if existingFriend.UserMail2 == friend.UserMail {
			existingFriend.FriendStatus = "block_first_second"
		} else {
			existingFriend.FriendStatus = "block_second_first"
		}
	case "block_first_second":
		if existingFriend.UserMail == friend.UserMail {
			existingFriend.FriendStatus = "block_both"
		}
	case "block_second_first":
		if existingFriend.UserMail == friend.UserMail2 {
			existingFriend.FriendStatus = "block_both"
		}
	}

	return r.DB.Debug().Save(&existingFriend).Error
}

//endregion
