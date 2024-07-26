package models

import "time"

type Room struct {
	RoomID        string `json:"room_id" gorm:"not null"`
	CreatedUserID string `json:"created_user_id" gorm:"not null"`
	CreatedAt      time.Time `json:"createdAt" gorm:"not null;column:createdAt"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updatedAt"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deletedAt"`
	MessageCount int64 `json:"message_count"`
}