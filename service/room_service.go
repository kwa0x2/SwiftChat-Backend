package service

import (
	"github.com/kwa0x2/swiftchat-backend/models"
	"github.com/kwa0x2/swiftchat-backend/repository"
	"gorm.io/gorm"
)

type IRoomService interface {
	Create(tx *gorm.DB, room *models.Room) (*models.Room, error)
	Update(tx *gorm.DB, whereRoom *models.Room, updateRoom *models.Room) error
	CreateAndAddUsers(createdUserId string, userId2 string) (string, error)
	GetChatList(userId, userEmail string) ([]*repository.ChatList, error)
}

type roomService struct {
	RoomRepository  repository.IRoomRepository
	UserRoomService IUserRoomService
}

func NewRoomService(roomRepository repository.IRoomRepository, UserRoomService IUserRoomService) IRoomService {
	return &roomService{
		RoomRepository:  roomRepository,
		UserRoomService: UserRoomService,
	}
}

// region "Create" adds a new room to the database
func (s *roomService) Create(tx *gorm.DB, room *models.Room) (*models.Room, error) {
	return s.RoomRepository.Create(tx, room)
}

//endregion

// region "Update" modifies the fields of a friend in the database based on specified conditions
func (s *roomService) Update(tx *gorm.DB, whereRoom *models.Room, updateRoom *models.Room) error {
	return s.RoomRepository.Update(tx, whereRoom, updateRoom)
}

// endregion

// region "CreateAndAddUsers" creates a new room and adds users to the user room within a transaction.
func (s *roomService) CreateAndAddUsers(createdUserId string, userId2 string) (string, error) {
	// Begin a new database transaction.
	tx := s.RoomRepository.GetDB().Begin()
	if tx.Error != nil {
		return "", tx.Error
	}

	// Create a new room object.
	roomObj := &models.Room{
		CreatedUserID: createdUserId, // Set the creator of the room.
	}

	// Create the room in the database.
	room, err := s.Create(tx, roomObj)
	if err != nil {
		tx.Rollback() // Roll back the transaction on error.
		return "", err
	}

	// Add users to the newly created room.
	for _, userId := range []string{createdUserId, userId2} {
		userRoom := &models.UserRoom{
			UserID: userId,      // Assign the user ID.
			RoomID: room.RoomID, // Assign the room ID.
		}
		// Create the user-room association.
		if createErr := s.UserRoomService.Create(tx, userRoom); createErr != nil {
			tx.Rollback() // Roll back the transaction on error.
			return "", createErr
		}
	}

	// Commit the transaction.
	if commitErr := tx.Commit().Error; commitErr != nil {
		return "", commitErr // Return an error if committing fails.
	}

	return room.RoomID.String(), nil
}

//endregion

// region "GetChatList" retrieves the list of chat rooms for a user, including last message details
func (s *roomService) GetChatList(userId, userEmail string) ([]*repository.ChatList, error) {
	return s.RoomRepository.GetChatList(userId, userEmail)
}

// endregion
