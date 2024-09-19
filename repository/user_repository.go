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

// region GET USER BY ID SERVICE
func (r *UserRepository) GetUserById(id string) (*models.User, error) {
	var user *models.User
	if err := r.DB.Table("USER").Where("user_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

//endregion

// region IS ID UNIQUE SERVICE
func (r *UserRepository) IsIdUnique(id string) bool {
	var count int64
	r.DB.Table("USER").Where("user_id = ?", id).Count(&count)
	return count == 0
}

//endregion

func (r *UserRepository) UpdateUsernameByMail(userName, userEmail string) error {
	if err := r.DB.Model(&models.User{}).Where("user_email = ?", userEmail).Update("user_name", userName).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateUserPhotoByMail(userPhoto, userEmail string) error {
	if err := r.DB.Model(&models.User{}).Where("user_email = ?", userEmail).Update("user_photo", userPhoto).Error; err != nil {
		return err
	}
	return nil
}
