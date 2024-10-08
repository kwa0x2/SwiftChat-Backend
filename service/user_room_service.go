package service

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/repository"
	"gorm.io/gorm"
)

type IUserRoomService interface {
	Create(tx *gorm.DB, userRoom *models.UserRoom) error
	GetPrivateRoom(userId1, userId2 string) (string, error)
}

type userRoomService struct {
	UserRoomRepository repository.IUserRoomRepository
}

func NewUserRoomService(UserRoomRepository repository.IUserRoomRepository) IUserRoomService {
	return &userRoomService{
		UserRoomRepository: UserRoomRepository,
	}
}

// region "Create" adds a new user room to the database
func (s *userRoomService) Create(tx *gorm.DB, userRoom *models.UserRoom) error {
	return s.UserRoomRepository.Create(tx, userRoom)
}

//endregion

// region "GetPrivateRoom" fetches the room ID of a private room for the specified user IDs
func (s *userRoomService) GetPrivateRoom(userId1, userId2 string) (string, error) {
	return s.UserRoomRepository.GetPrivateRoom(userId1, userId2)
}

//endregion
