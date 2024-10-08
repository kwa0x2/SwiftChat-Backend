package models

import (
	"github.com/google/uuid"
	"github.com/kwa0x2/swiftchat-backend/types"
	"gorm.io/gorm"
	"time"
)

type Room struct {
	RoomID        uuid.UUID      `json:"room_id" gorm:"not null;type:uuid;default:gen_random_uuid()"`
	CreatedUserID string         `json:"created_user_id" gorm:"not null;primaryKey"`
	MessageCount  int64          `json:"message_count"`
	LastMessage   string         `json:"last_message"`
	RoomType      types.RoomType `json:"room_type" gorm:"not null;type:room_type;default:private"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"not null;column:createdAt;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `json:"updatedAt" gorm:"not null;column:updatedAt;default:CURRENT_TIMESTAMP"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt" gorm:"column:deletedAt"`
	LastMessageID uuid.UUID      `json:"message_id" gorm:"type:uuid"`

	UserRoom UserRoom `json:"user_room" gorm:"foreignKey:RoomID;references:RoomID"`
}

func (Room) TableName() string {
	return "ROOM"
}
