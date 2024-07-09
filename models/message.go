package models

import (
	"database/sql"
	"time"
)

type ReadStatus string

const (
	Unread ReadStatus = "unread"
	Readed ReadStatus = "readed"
)

type Message struct {
	MessageID      uint64 `json:"message_id" gorm:"not null;primaryKey;autoIncrement"`
	MessageContent string `json:"message_content" gorm:"not null;size:10000"`
	CreatedAt      time.Time `json:"createdAt" gorm:"not null;column:createdAt"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"not null;column:updatedAt"`
	DeletedAt sql.NullTime `json:"deletedAt" gorm:"column:deletedAt"`
	MessageSenderID string `json:"message_sender_id" gorm:"not null"`
	MessageReceiverID string `json:"message_receiver_id" gorm:"not null"`
	MessageReadStatus ReadStatus `json:"message_read_status" gorm:"type:read_status;not null;default:unread"`
}
