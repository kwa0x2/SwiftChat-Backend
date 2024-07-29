package repository

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"gorm.io/gorm"
)

type UserRoomRepository struct {
	DB *gorm.DB
}

// region GET PRIVATE ROOM REPOSITORY
func (r *UserRoomRepository) GetPrivateRoom(userId1, userId2 string) (string, error) {
	var roomId string

	if err := r.DB.Model(&models.UserRoom{}).
		Select("\"USER_ROOM\".room_id").
		Joins("JOIN \"ROOM\" ON \"USER_ROOM\".room_id = \"ROOM\".room_id").
		Where("\"USER_ROOM\".user_id IN (?,?) AND \"ROOM\".room_type = ?", userId1, userId2, "private").
		Group("\"USER_ROOM\".room_id").
		Having("COUNT(DISTINCT user_id) = 2").
		Find(&roomId).Error; err != nil {
		return "", err
	}

	return roomId, nil
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
