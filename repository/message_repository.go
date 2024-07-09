package repository

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"gorm.io/gorm"
)

type MessageRepository struct {
	DB *gorm.DB
}

func (r *MessageRepository) InsertMessage(message *models.Message) (*models.Message, error) {
	if err := r.DB.Table("MESSAGE").Create(&message).Error; err != nil {
		return nil, err	
	}
	return message, nil
}