package repository

import (
	"github.com/google/uuid"
	"github.com/kwa0x2/realtime-chat-backend/models"
	"gorm.io/gorm"
)

type UserRoomRepository struct {
	DB *gorm.DB
}

type PrivateRoom struct {
	RoomID       uuid.UUID `json:"room_id"`
	FriendStatus string    `json:"friend_status"`
}

// region GET PRIVATE ROOM REPOSITORY
func (r *UserRoomRepository) GetPrivateRoom(userId1, userId2 string) (string, error) {
	var userRooms string

	//if err := r.DB.Model(&models.UserRoom{}).Debug().Select("\"USER_ROOM\".room_id").
	//	Preload("Room", func(db *gorm.DB) *gorm.DB {
	//		return db.Where("room_type = ?", "private")
	//	}).Preload("User.Friend", func(db *gorm.DB) *gorm.DB {
	//	return db.Select("friend_status")
	//}).
	//	Where("\"USER_ROOM\".user_id IN (?,?)", userId1, userId2).
	//	Group("\"USER_ROOM\".room_id").
	//	Having("COUNT(DISTINCT user_id) = 2").
	//	Scan(&userRooms).Error; err != nil {
	//	return nil, err
	//}

	//if err := r.DB.Model("USER_ROOM").Debug().
	//	Select("\"USER_ROOM\".room_id, \"FRIEND\".friend_status").
	//	Joins("JOIN \"ROOM\" ON \"USER_ROOM\".room_id = \"ROOM\".room_id").
	//	Joins("JOIN \"USER\" ON \"USER_ROOM\".user_id = \"USER\".user_id").
	//	Joins("LEFT JOIN \"FRIEND\" ON \"USER\".user_email = \"FRIEND\".user_mail").
	//	Where("\"USER_ROOM\".user_id IN (?,?) AND \"ROOM\".room_type = ?", userId1, userId2, "private").
	//	Group("\"USER_ROOM\".room_id, \"FRIEND\".friend_status").
	//	Having("COUNT(DISTINCT \"USER_ROOM\".user_id) = 2").
	//	Scan(&userRooms).Error; err != nil {
	//	return nil, err
	//}

	if err := r.DB.Model(&models.UserRoom{}).
		Select("\"USER_ROOM\".room_id").
		Joins("JOIN \"ROOM\" ON \"USER_ROOM\".room_id = \"ROOM\".room_id").
		Where("\"USER_ROOM\".user_id IN (?,?) AND \"ROOM\".room_type = ?", userId1, userId2, "private").
		Group("\"USER_ROOM\".room_id").
		Having("COUNT(DISTINCT user_id) = 2").
		Find(&userRooms).Error; err != nil {
		return "", err
	}

	return userRooms, nil
}

//endregion

// region CREATE USER ROOM REPOSITORY
func (r *UserRoomRepository) CreateUserRoom(tx *gorm.DB, userRoom *models.UserRoom) error {
	db := r.DB
	if tx != nil {
		db = tx
	}

	if err := db.Create(&userRoom).Error; err != nil {
		return err
	}
	return nil
}

//endregion
