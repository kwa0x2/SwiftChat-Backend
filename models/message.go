package models

import (
	"time"

	"github.com/kwa0x2/realtime-chat-backend/types"
)



type Message struct {
	MessageID     uint64     `json:"message_id" gorm:"not null;primaryKey;autoIncrement"`
	Message       string     `json:"message" gorm:"not null;size:10000"`
	SenderID      string     `json:"sender_id" gorm:"not null"`
	RoomID        string     `json:"room_id" gorm:"not null"`
	MessageStatus types.ReadStatus `json:"message_status" gorm:"type:read_status;not null;default:unread"`
	CreatedAt     time.Time  `json:"createdAt" gorm:"not null;column:createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt" gorm:"not null;column:updatedAt"`
	DeletedAt     time.Time  `json:"deletedAt" gorm:"column:deletedAt"`
}

func (Message) TableName() string {
	return "MESSAGE"
}
