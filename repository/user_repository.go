package repository

import (
	"github.com/kwa0x2/swiftchat-backend/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	IsFieldUnique(whereUser *models.User) bool
	IsFieldExists(whereUser *models.User) bool
	Create(user *models.User) (*models.User, error)
	GetUser(whereUser *models.User) (*models.User, error)
	Update(whereUser *models.User, updates *models.User) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{
		DB: db,
	}
}

// region "IsFieldUnique" checks if the specified fields in the User model are unique.
func (r *userRepository) IsFieldUnique(whereUser *models.User) bool {
	var count int64
	r.DB.Model(&models.User{}).Where(whereUser).Count(&count)
	return count == 0
}

// endregion

// region "IsFieldExists" checks if the specified fields in the User model exist in the database.
func (r *userRepository) IsFieldExists(whereUser *models.User) bool {
	var count int64
	r.DB.Model(&models.User{}).Where(whereUser).Count(&count)
	return count > 0
}

// endregion

// region "Create" adds a new user to the database and returns the created user.
func (r *userRepository) Create(user *models.User) (*models.User, error) {
	if err := r.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// endregion

// region "GetUser" retrieves a user from the database based on the specified conditions.
func (r *userRepository) GetUser(whereUser *models.User) (*models.User, error) {
	var user *models.User
	if err := r.DB.Where(whereUser).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// endregion

// region "Update" modifies the fields of a user in the database based on specified conditions.
func (r *userRepository) Update(whereUser *models.User, updates *models.User) error {
	return r.DB.Model(&models.User{}).Where(whereUser).Updates(updates).Error
}

// endregion
