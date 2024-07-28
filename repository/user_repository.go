package repository

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

// region IS USERNAME UNIQUE REPOSITORY
func (r *UserRepository) IsUsernameUnique(username string) bool {
	var count int64
	r.DB.Where("user_name = ?", username).Count(&count)
	return count == 0
}

//endregion

// region IS EMAIL UNIQUE REPOSITORY
func (r *UserRepository) IsEmailUnique(email string) bool {
	var count int64
	r.DB.Where("user_email = ?", email).Count(&count)
	return count == 0
}

//endregion

// region INSERT NEW USER REPOSITORY
func (r *UserRepository) Insert(user *models.User) (*models.User, error) {
	if err := r.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

//endregion

// region GET USER BY EMAIL REPOSITORY
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user *models.User
	if err := r.DB.Table("USER").Where("user_email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

//endregion

// region GET USERNAME BY ID SERVICE
func (r *UserRepository) GetUsernameById(id string) string {
	var username string
	r.DB.Table("USER").Select("user_name").Where("user_id = ?", id).Scan(&username)
	return username
}

//endregion

// region IS ID UNIQUE SERVICE
func (r *UserRepository) IsIdUnique(id string) bool {
	var count int64
	r.DB.Table("USER").Where("user_id = ?", id).Count(&count)
	return count == 0
}

//endregion
