package service

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/repository"
)

type IUserService interface {
	IsUsernameUnique(userName string) bool
	IsIdUnique(userId string) bool
	IsEmailExists(email string) bool
	Create(user *models.User) (*models.User, error)
	GetByEmail(userEmail string) (*models.User, error)
	GetUserById(userId string) (*models.User, error)
	UpdateUserNameByMail(userName, userEmail string) error
	UpdateUserPhotoByMail(userPhoto, userEmail string) error
}

type userService struct {
	UserRepository repository.IUserRepository
}

func NewUserService(userRepo repository.IUserRepository) IUserService {
	return &userService{
		UserRepository: userRepo,
	}
}

// region "IsUsernameUnique" checks if the given username is unique
func (s *userService) IsUsernameUnique(userName string) bool {
	return s.UserRepository.IsFieldUnique(&models.User{UserName: userName})
}

// endregion

// region "IsIdUnique" checks if the given user ID is unique
func (s *userService) IsIdUnique(userId string) bool {
	return s.UserRepository.IsFieldUnique(&models.User{UserID: userId})
}

// endregion

// region "IsEmailExists" checks if the given email already exists in the database
func (s *userService) IsEmailExists(email string) bool {
	return s.UserRepository.IsFieldExists(&models.User{UserEmail: email})
}

// endregion

// region "Create" adds a new user to the database
func (s *userService) Create(user *models.User) (*models.User, error) {
	return s.UserRepository.Create(user)
}

// endregion

// region "GetByEmail" retrieves a user from the database by their email
func (s *userService) GetByEmail(userEmail string) (*models.User, error) {
	return s.UserRepository.GetUser(&models.User{UserEmail: userEmail})
}

// endregion

// region "GetUserById" retrieves a user from the database by their ID
func (s *userService) GetUserById(userId string) (*models.User, error) {
	return s.UserRepository.GetUser(&models.User{UserID: userId})
}

// endregion

// region "UpdateUserNameByMail" updates the user's name based on their email
func (s *userService) UpdateUserNameByMail(userName, userEmail string) error {
	whereUser := &models.User{
		UserEmail: userEmail, // User to find based on email.
	}
	updates := &models.User{
		UserName: userName, // New username to set.
	}

	return s.UserRepository.Update(whereUser, updates)
}

// endregion

// region "UpdateUserPhotoByMail" updates the user's photo based on their email
func (s *userService) UpdateUserPhotoByMail(userPhoto, userEmail string) error {
	whereUser := &models.User{
		UserEmail: userEmail, // User to find based on email.
	}
	updates := &models.User{
		UserPhoto: userPhoto, // New photo URL to set.
	}

	return s.UserRepository.Update(whereUser, updates)
}

// endregion
