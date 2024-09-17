package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"

	"github.com/kwa0x2/realtime-chat-backend/types"
)

type Message struct {
	MessageID     uuid.UUID        `json:"message_id" gorm:"not null;type:uuid;primaryKey;autoIncrement"`
	Message       string           `json:"message" gorm:"not null;size:10000"`
	SenderID      string           `json:"sender_id" gorm:"not null"`
	RoomID        uuid.UUID        `json:"room_id" gorm:"not null;type:uuid"`
	MessageStatus types.ReadStatus `json:"message_status" gorm:"type:read_status;not null;default:unread"`
	CreatedAt     time.Time        `json:"createdAt" gorm:"not null;column:createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt" gorm:"not null;column:updatedAt"`
	DeletedAt     gorm.DeletedAt   `json:"deletedAt" gorm:"column:deletedAt"`
}

func (Message) TableName() string {
	return "MESSAGE"
}
