package repository

import (
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/kwa0x2/swiftchat-backend/models"
	"gorm.io/gorm"
)

type IUserRoomRepository interface {
	Create(tx *gorm.DB, userRoom *models.UserRoom) error
	GetRoom(userId1, userId2 string) (string, error)
}

type userRoomRepository struct {
	DB *gorm.DB
}

func NewUserRoomRepository(db *gorm.DB) IUserRoomRepository {
	return &userRoomRepository{
		DB: db,
	}
}

// region "Create" adds a new user room to the database
func (r *userRoomRepository) Create(tx *gorm.DB, userRoom *models.UserRoom) error {
	db := r.DB
	if tx != nil {
		db = tx // Use the provided transaction if available
	}

	if err := db.Create(&userRoom).Error; err != nil {
		sentry.CaptureException(err)
		return err
	}
	return nil
}

//endregion

// region GetRoom DTO represents the structure of a room.
type GetRoom struct {
	RoomID       uuid.UUID `json:"room_id"`
	FriendStatus string    `json:"friend_status"`
}

// endregion

// region "GetRoom" fetches the room ID of a room for the specified user IDs
func (r *userRoomRepository) GetRoom(userId1, userId2 string) (string, error) {
	var userRooms string

	if err := r.DB.Model(&models.UserRoom{}).
		Select(`"USER_ROOM".room_id`).
		Joins(`JOIN "ROOM" ON "USER_ROOM".room_id = "ROOM".room_id`).
		Where(`"USER_ROOM".user_id IN (?,?)`, userId1, userId2).
		Group(`"USER_ROOM".room_id`).
		Having("COUNT(DISTINCT user_id) = 2").
		Find(&userRooms).Error; err != nil {
		sentry.CaptureException(err)
		return "", err
	}

	return userRooms, nil
}

//endregion
