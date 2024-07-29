package models

import (
	"github.com/google/uuid"
	"github.com/kwa0x2/realtime-chat-backend/types"
	"time"
)

type Room struct {
	RoomID        uuid.UUID      `json:"room_id" gorm:"not null;default:gen_random_uuid()"`
	CreatedUserID string         `json:"created_user_id" gorm:"not null"`
	MessageCount  int64          `json:"message_count"`
	LastMessage   string         `json:"last_message"`
	RoomType      types.RoomType `json:"room_type" gorm:"not null;type:room_type;default:private"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"not null;column:createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt" gorm:"column:updatedAt"`
	DeletedAt     time.Time      `json:"deletedAt" gorm:"column:deletedAt"`
}

func (Room) TableName() string {
	return "ROOM"
}
