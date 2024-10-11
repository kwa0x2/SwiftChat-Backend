package repository

import (
	"errors"
	"fmt"
	"github.com/kwa0x2/swiftchat-backend/models"
	"github.com/kwa0x2/swiftchat-backend/types"
	"gorm.io/gorm"
)

type IFriendRepository interface {
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

type friendRepository struct {
	DB *gorm.DB
}

func NewFriendRepository(db *gorm.DB) IFriendRepository {
	return &friendRepository{
		DB: db,
	}
}

// region "Create" adds a new friend to the database
func (r *friendRepository) Create(tx *gorm.DB, friend *models.Friend) error {
	db := r.DB // Use the repository's DB connection
	if tx != nil {
		db = tx // If a transaction is provided, use it
	}

	if err := db.Create(&friend).Error; err != nil {
		return err
	}
	return nil
}

// endregion

// region "Update" modifies the fields of a friend in the database based on specified conditions
func (r *friendRepository) Update(tx *gorm.DB, whereFriend *models.Friend, updates *models.Friend) error {
	db := r.DB
	if tx != nil {
		db = tx
	}

	result := db.Model(&models.Friend{}).Unscoped().Where(whereFriend).Updates(updates)

	if result.Error != nil {
		return result.Error // Return any error that occurs during the update
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Return an error if no records were affected
	}

	return nil
}

// endregion

// region "UpdateFriendStatusByMail" updates the deletedAt field and friendStatus for given user emails
func (r *friendRepository) UpdateFriendStatusByMail(tx *gorm.DB, userEmail, userEmail2 string, friendStatus types.FriendStatus) error {
	db := r.DB
	if tx != nil {
		db = tx
	}

	result := db.Model(&models.Friend{}).Unscoped().
		Where("(user_email = ? AND user_email2 = ?) OR (user_email2 = ? AND user_email = ?)", userEmail, userEmail2, userEmail, userEmail2).
		Updates(map[string]interface{}{
			"friend_status": friendStatus,
			"deletedAt":     nil, // Reset the deletedAt field
		})

	if result.Error != nil {
		return result.Error // Return any error that occurs during the update
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Return an error if no records were affected
	}

	return nil
}

// endregion

// region "Delete" removes a friend relationship from the database
func (r *friendRepository) Delete(UserEmail, UserEmail2 string) error {
	if err := r.DB.
		Where("(user_email = ? AND user_email2 = ?) OR (user_email = ? AND user_email2 = ?)",
			UserEmail,
			UserEmail2,
			UserEmail2,
			UserEmail).
		Updates(&models.Friend{FriendStatus: types.UnFriend}).
		Delete(&models.Friend{}).Error; err != nil {
		return err
	}
	return nil
}

// endregion

// region "GetFriends" retrieves a list of friends based on user email
func (r *friendRepository) GetFriends(userEmail string, isUnFriendStatusAllow bool) ([]*models.Friend, error) {
	var friends []*models.Friend

	query := r.DB

	// Adjust the query based on the isUnFriendStatusAllow flag
	if isUnFriendStatusAllow {
		query = query.Unscoped().Where("user_email = ? AND (friend_status = ? OR friend_status = ?)", userEmail, types.Friend, types.UnFriend).
			Or("user_email2 = ? AND (friend_status = ? OR friend_status = ?)", userEmail, types.Friend, types.UnFriend)
	} else {
		query = query.Where("user_email = ? AND friend_status = ?", userEmail, types.Friend).
			Or("user_email2 = ? AND friend_status = ?", userEmail, types.Friend)
	}

	if err := query.Preload("User").
		Select("CASE WHEN user_email = ? THEN user_email2 ELSE user_email END as user_email", userEmail).
		Find(&friends).Error; err != nil {
		return nil, err
	}

	return friends, nil
}

// endregion

// region "GetSpecificFriend" retrieves a specific friend relationship based on user emails
func (r *friendRepository) GetSpecificFriend(userEmail, userEmail2 string) (*models.Friend, error) {
	var friend models.Friend

	if err := r.DB.Unscoped().
		Select("friend_status, CASE WHEN user_email = ? THEN user_email2 ELSE user_email END as user_email", userEmail).
		Where("(user_email = ? AND user_email2 = ?)", userEmail, userEmail2).
		Or("(user_email2 = ? AND user_email = ?)", userEmail, userEmail2).
		First(&friend).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &friend, nil
}

// endregion

// region "GetBlockedUsers" retrieves a list of blocked users for a given email
func (r *friendRepository) GetBlockedUsers(userEmail string) ([]*models.Friend, error) {
	var friends []*models.Friend

	if err := r.DB.
		Preload("User").
		Select("CASE WHEN user_email = ? THEN user_email2 ELSE user_email END as user_email", userEmail).
		Where("user_email = ? AND friend_status = ?", userEmail, types.Blocked).
		Or("user_email2 = ? AND friend_status = ?", userEmail, types.Blocked).
		Find(&friends).Error; err != nil {
		return nil, err
	}

	return friends, nil
}

// endregion

// region "Block" updates the status of a friendship to blocked
func (r *friendRepository) Block(userEmail, userEmail2 string) (string, error) {
	// Check if the friendship exists
	if err := r.DB.
		Where("(user_email = ? AND user_email2 = ?) OR (user_email = ? AND user_email2 = ?)", userEmail, userEmail2, userEmail2, userEmail).
		Updates(&models.Friend{FriendStatus: types.Blocked}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("friendship not found")
		}
		return "", err
	}

	return "", nil
}

// endregion

// region "IsBlocked" checks if a user is blocked by another user
func (r *friendRepository) IsBlocked(userEmail, userEmail2 string) (bool, error) {
	var count int64

	if err := r.DB.Model(&models.Friend{}).
		Where("(user_email = ? OR user_email2 = ?) AND (user_email = ? OR user_email2 = ?)",
			userEmail, userEmail, userEmail2, userEmail2).
		Where("friend_status = ?", types.Blocked).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// endregion
