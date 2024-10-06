package service

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/repository"
	"github.com/kwa0x2/realtime-chat-backend/types"
	"gorm.io/gorm"
)

type RoomService struct {
	RoomRepository  *repository.RoomRepository
	UserRoomService *UserRoomService
}

// region CREATE ROOM SERVICE
func (s *RoomService) CreateRoom(tx *gorm.DB, room *models.Room) (*models.Room, error) {
	return s.RoomRepository.CreateRoom(tx, room)
}

//endregion

func (s *RoomService) Update(tx *gorm.DB, room *models.Room) error {
	return s.RoomRepository.Update(tx, room)
}

// region CREATE ROOM AND ADD USERS IN USER ROOM WITH TRANSACTION SERVICE
func (s *RoomService) CreateAndAddUsers(createdUserId string, userId2 string, roomType types.RoomType) (string, error) {
	tx := s.RoomRepository.DB.Begin()
	if tx.Error != nil {
		return "", tx.Error
	}

	roomObj := &models.Room{
		CreatedUserID: createdUserId,
		RoomType:      roomType,
	}

	room, err := s.CreateRoom(tx, roomObj)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	for _, userId := range []string{createdUserId, userId2} {
		userRoom := &models.UserRoom{
			UserID: userId,
			RoomID: room.RoomID,
		}
		if err := s.UserRoomService.CreateUserRoom(tx, userRoom); err != nil {
			tx.Rollback()
			return "", err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	return room.RoomID.String(), nil
}

//endregion

func (s *RoomService) GetChatList(userId, userMail string) ([]*repository.ChatList, error) {
	return s.RoomRepository.GetChatList(userId, userMail)
}
