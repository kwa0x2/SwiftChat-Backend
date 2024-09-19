package service

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/repository"
)

type UserService struct {
	UserRepository *repository.UserRepository
}

// region IS USERNAME UNIQUE SERVICE
func (s *UserService) IsUsernameUnique(username string) bool {
	return s.UserRepository.IsUsernameUnique(username)
}

//endregion

// region IS EMAIL UNIQUE SERVICE
func (s *UserService) IsEmailUnique(email string) bool {
	return s.UserRepository.IsEmailUnique(email)
}

//endregion

// region INSERT NEW USER SERVICE
func (s *UserService) Insert(user *models.User) (*models.User, error) {
	return s.UserRepository.Insert(user)
}

//endregion

// region GET USER BY EMAIL SERVICE
func (s *UserService) GetByEmail(email string) (*models.User, error) {
	return s.UserRepository.GetByEmail(email)
}

//endregion

// region GET USER BY ID SERVICE
func (s *UserService) GetUserById(id string) (*models.User, error) {
	return s.UserRepository.GetUserById(id)
}

//endregion

// region IS ID UNIQUE SERVICE
func (s *UserService) IsIdUnique(id string) bool {
	return s.UserRepository.IsIdUnique(id)
}

//endregion

func (s *UserService) UpdateUsernameByMail(userName, userEmail string) error {
	return s.UserRepository.UpdateUsernameByMail(userName, userEmail)
}

func (s *UserService) UpdateUserPhotoByMail(userPhoto, userEmail string) error {
	return s.UserRepository.UpdateUserPhotoByMail(userPhoto, userEmail)
}
