package repository

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"gorm.io/gorm"
)

type FriendRepository struct {
	DB *gorm.DB
}

type ComingRequests struct {
	SenderID  string `json:"sender_id"`
	UserName  string `json:"user_name"`
	UserPhoto string `json:"user_photo"`
}

func (r *FriendRepository) Insert(friend *models.Friend) (error) {
	if err := r.DB.Create(&friend).Error; err != nil {
		return err
	}
	return nil
}

func (r *FriendRepository) Delete(friend *models.Friend) error {
	if err := r.DB.
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


func (r *FriendRepository) GetFriends(userMail string) ([]*models.Friend, error) {
	var friend []*models.Friend

	if err := r.DB.
		Preload("User").
		Select("CASE WHEN user_mail = ? THEN user_mail2 ELSE user_mail END as user_mail", userMail).
		Where("user_mail = ? AND status = ?", userMail, "friend").
		Or("user_mail2 = ? AND status = ?", userMail, "friend").
		Find(&models.Friend{}).Error; err != nil {
		return nil, err
	}

	return friend, nil
}

func (r *FriendRepository) GetBlockeds(userMail string) ([]*models.Friend, error) {
	var friend []*models.Friend

	if err := r.DB.
		Preload("User").
		Select("CASE WHEN user_mail = ? THEN user_mail2 ELSE user_mail END as user_mail", userMail).
		Where("user_mail = ? AND status = ?", userMail, "blocked").
		Or("user_mail2 = ? AND status = ?", userMail, "blocked").
		Find(&models.Friend{}).Error; err != nil {
		return nil, err
	}

	return friend, nil
}

func (r *FriendRepository) Block(friend *models.Friend) error {
	if err := r.DB.Where("(user_mail = ? AND user_mail2 = ?) OR (user_mail = ? AND user_mail2 = ?)", friend.UserMail, friend.UserMail2, friend.UserMail2, friend.UserMail).Update("status", "blocked").Error; err != nil {
		return err
	}

	return nil
}
