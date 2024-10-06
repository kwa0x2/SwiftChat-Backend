package repository

import (
	"errors"
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/types"
	"gorm.io/gorm"
)

type FriendRepository struct {
	DB *gorm.DB
}

// region INSERT NEW FRIEND REPOSITORY
func (r *FriendRepository) Insert(tx *gorm.DB, friend *models.Friend) error {
	db := r.DB
	if tx != nil {
		db = tx
	}

	if err := db.Create(&friend).Error; err != nil {
		return err
	}
	return nil
}

//endregion

func (r *FriendRepository) Update(tx *gorm.DB, filter map[string]interface{}, updates map[string]interface{}) error {
	db := r.DB
	if tx != nil {
		db = tx
	}

	result := db.Model(&models.Friend{}).Unscoped().Where(filter).Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *FriendRepository) UpdateDeletedAtByMail(tx *gorm.DB, userMail, userMail2 string, friendStatus types.FriendStatus) error {
	db := r.DB
	if tx != nil {
		db = tx
	}

	result := db.Model(&models.Friend{}).Unscoped().
		Where("(user_mail = ? AND user_mail2 = ?)", userMail, userMail2).
		Or("(user_mail2 = ? AND user_mail = ?)", userMail, userMail2).
		Updates(map[string]interface{}{
			"friend_status": friendStatus,
			"deletedAt":     nil,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// region DELETE FRIEND BY MAIL REPOSITORY
func (r *FriendRepository) Delete(friend *models.Friend) error {
	if err := r.DB.
		Where("(user_mail = ? AND user_mail2 = ?) OR (user_mail = ? AND user_mail2 = ?)",
			friend.UserMail,
			friend.UserMail2,
			friend.UserMail2,
			friend.UserMail).
		Updates(&models.Friend{FriendStatus: types.UnFriend}).
		Delete(&models.Friend{}).Error; err != nil {
		return err
	}
	return nil
}

//endregion

// region GET FRIENDS BY MAIL REPOSITORY
func (r *FriendRepository) GetFriends(userMail string, isUnFriendStatusAllow bool) ([]*models.Friend, error) {
	var friends []*models.Friend

	query := r.DB

	if isUnFriendStatusAllow {
		query = query.Unscoped().
			Where("user_mail = ? AND (friend_status = ? OR friend_status = ?)", userMail, "friend", "unfriend").
			Or("user_mail2 = ? AND (friend_status = ? OR friend_status = ?)", userMail, "friend", "unfriend")
	} else {
		query = query.Where("user_mail = ? AND friend_status = ?", userMail, "friend").
			Or("user_mail2 = ? AND friend_status = ?", userMail, "friend")
	}

	if err := query.
		Preload("User").
		Select("CASE WHEN user_mail = ? THEN user_mail2 ELSE user_mail END as user_mail", userMail).
		Find(&friends).Error; err != nil {
		return nil, err
	}

	return friends, nil
}

//endregion

func (r *FriendRepository) GetFriend(userMail, userMail2 string) (*models.Friend, error) {
	var friend *models.Friend

	if err := r.DB.Unscoped().
		Select("friend_status, CASE WHEN user_mail = ? THEN user_mail2 ELSE user_mail END as user_mail", userMail).
		Where("(user_mail = ? AND user_mail2 = ?)", userMail, userMail2).
		Or("(user_mail2 = ? AND user_mail = ?)", userMail, userMail2).
		First(&friend).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return friend, nil
}

// region GET BLOCKED FRIENDS BY MAIL SERVICE
func (r *FriendRepository) GetBlocked(userMail string) ([]*models.Friend, error) {
	var friends []*models.Friend

	if err := r.DB.
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
func (r *FriendRepository) Block(friend *models.Friend) (string, error) {
	var existingFriend models.Friend

	if err := r.DB.
		Where("(user_mail = ? AND user_mail2 = ?) OR (user_mail = ? AND user_mail2 = ?)", friend.UserMail, friend.UserMail2, friend.UserMail2, friend.UserMail).
		First(&existingFriend).Error; err != nil {
		return "", err

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

	if err := r.DB.Save(&existingFriend).Error; err != nil {
		return "", err
	}

	return string(existingFriend.FriendStatus), nil
}

//endregion

func (r *FriendRepository) IsBlocked(userMail, otherUserMail string) (bool, error) {
	var count int64

	if err := r.DB.Model(&models.Friend{}).
		Where("(user_mail = ? OR user_mail2 = ?) AND (user_mail = ? OR user_mail2 = ?)",
			userMail, userMail, otherUserMail, otherUserMail).
		Where("friend_status != ?", "friend").
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
