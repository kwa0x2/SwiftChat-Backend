package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Room struct {
	RoomID        uuid.UUID      `json:"room_id" gorm:"not null;type:uuid;default:gen_random_uuid()"`
	CreatedUserID string         `json:"created_user_id" gorm:"not null;primaryKey"`
	LastMessageID uuid.UUID      `json:"message_id" gorm:"type:uuid"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"not null;column:createdAt;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `json:"updatedAt" gorm:"not null;column:updatedAt;default:CURRENT_TIMESTAMP"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt" gorm:"column:deletedAt"`

	UserRoom UserRoom `json:"user_room" gorm:"foreignKey:RoomID;references:RoomID"`
}

func (Room) TableName() string {
	return "ROOM"
}
