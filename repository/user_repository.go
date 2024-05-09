package repository

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) IsUsernameUnique(username string) bool {
	var count int64
	r.DB.Table("USER").Where("user_name = ?", username).Count(&count)
	return count == 0
}



func (r *UserRepository) InsertUser(user *models.User) error {
    if err := r.DB.Table("USER").Create(&user).Error; err != nil {
        return err
    }
    return nil
}


