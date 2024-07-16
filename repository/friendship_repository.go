package repository

import (
	"fmt"

	"github.com/kwa0x2/realtime-chat-backend/models"
	"gorm.io/gorm"
)

type FriendshipRepository struct {
	DB *gorm.DB
}

func (r *FriendshipRepository) SendFriendRequest(friendship *models.Friendship) (string, error) {
	if err := r.DB.Table("FRIENDSHIP").Create(&friendship).Error; err != nil {
		return "", err
	}
	return string(friendship.FriendshipStatus), nil
}

type ComingRequests struct {
	SenderID  string `json:"sender_id"`
	UserName  string `json:"user_name"`
	UserPhoto string `json:"user_photo"`
}

func (r *FriendshipRepository) GetComingRequests(receiverId string) ([]*models.Friendship, error) {
	var friendship []*models.Friendship

	if err := r.DB.Table("FRIENDSHIP").
		Preload("User").
		Where("receiver_id = ? AND friendship_status = ?", receiverId, "pending").
		Find(&friendship).Error; err != nil {
		return nil, err
	}

	return friendship, nil
}

func (r *FriendshipRepository) GetFriends(userId string) ([]*models.Friendship, error) {
	var friendship []*models.Friendship

	if err := r.DB.Table("FRIENDSHIP").
		Preload("User").
		Where("(receiver_id = ? OR sender_id = ?) AND friendship_status = ?", userId,userId, "accepted").
		Find(&friendship).Error; err != nil {
		return nil, err
	}

	return friendship, nil
}

func (r *FriendshipRepository) GetBlockeds(userId string) ([]*models.Friendship, error) {
	var friendship []*models.Friendship

	if err := r.DB.Table("FRIENDSHIP").
		Preload("User").
		Where("(receiver_id = ? OR sender_id = ?) AND friendship_status = ?", userId, userId, "blocked").
		Find(&friendship).Error; err != nil {
		return nil, err
	}

	return friendship, nil
}

func (r *FriendshipRepository) Block(friendship *models.Friendship) error {
	if err := r.DB.Table("FRIENDSHIP").Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", friendship.SenderId, friendship.ReceiverId, friendship.ReceiverId, friendship.SenderId).Update("friendship_status", "blocked").Error; err != nil {
		return err
	}

	return nil
}

func (r *FriendshipRepository) Delete(friendship *models.Friendship) error {
	
	if err := r.DB.Table("FRIENDSHIP").Where("(sender_id = ? AND receiver_id = ?) OR (receiver_id = ? AND sender_id = ?) AND friendship_status = ?", friendship.SenderId, friendship.ReceiverId, friendship.ReceiverId, friendship.SenderId, "accepted").Delete("FRIENDSHIP").Error; err != nil {
		return err
	}
	return nil
}

func (r *FriendshipRepository) Accept(friendship *models.Friendship) error {
	if err := r.DB.Table("FRIENDSHIP").Where("sender_id = ? AND receiver_id = ?", friendship.SenderId, friendship.ReceiverId).Update("friendship_status", "accepted").Error; err != nil {
		return err
	}
	return nil
}

func (r *FriendshipRepository) Reject(friendship *models.Friendship) error {
	if err := r.DB.Table("FRIENDSHIP").Where("sender_id = ? AND receiver_id = ?", friendship.SenderId, friendship.ReceiverId).Delete("FRIENDSHIP").Error; err != nil {
		return err
	}
	return nil
}
