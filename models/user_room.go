package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type UserRoom struct {
	UserID    string         `json:"user_id" gorm:"primaryKey;not null"`
	RoomID    uuid.UUID      `json:"room_id" gorm:"primaryKey;not null;type:uuid"`
	CreatedAt time.Time      `json:"createdAt" gorm:"not null;column:createdAt;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"not null;column:updatedAt;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"column:deletedAt"`

	User User  `json:"user" gorm:"foreignKey:UserID;references:UserID"`
	Room *Room `json:"room" gorm:"foreignKey:RoomID;references:RoomID"`
}

func (UserRoom) TableName() string {
	return "USER_ROOM"
}
