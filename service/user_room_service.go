package service

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/repository"
	"gorm.io/gorm"
)

type UserRoomService struct {
	UserRoomRepository *repository.UserRoomRepository
}

// region GET PRIVATE ROOM SERVICE
func (s *UserRoomService) GetPrivateRoom(userId1, userId2 string) (string, error) {
	return s.UserRoomRepository.GetPrivateRoom(userId1, userId2)
}

//endregion

// region CREATE USER ROOM SERVICE
func (s *UserRoomService) CreateUserRoom(tx *gorm.DB, userRoom *models.UserRoom) error {
	return s.UserRoomRepository.CreateUserRoom(tx, userRoom)
}

//endregion
