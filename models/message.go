package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"

	"github.com/kwa0x2/swiftchat-backend/types"
)

type Message struct {
	MessageID         uuid.UUID         `json:"message_id" gorm:"not null;type:uuid;primaryKey;autoIncrement"`
	Message           string            `json:"message" gorm:"not null;size:10000"`
	SenderID          string            `json:"sender_id" gorm:"not null"`
	RoomID            uuid.UUID         `json:"room_id" gorm:"not null;type:uuid"`
	MessageReadStatus types.ReadStatus  `json:"message_read_status" gorm:"type:read_status;not null;default:unread"`
	MessageType       types.MessageType `json:"message_type" gorm:"type:message_type;not null;default:text"`
	MessageStarred    bool              `json:"message_starred" gorm:"not null;default:false"`

	CreatedAt time.Time      `json:"createdAt" gorm:"not null;column:createdAt;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"not null;column:updatedAt;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"column:deletedAt"`
}

func (Message) TableName() string {
	return "MESSAGE"
}
