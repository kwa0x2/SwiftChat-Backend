package models

import (
	"github.com/google/uuid"
	"time"
)

type UserRoom struct {
	UserID    string    `json:"user_id" gorm:"primaryKey;not null"`
	RoomID    uuid.UUID `json:"room_id" gorm:"primaryKey;not null;type:uuid"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null;column:createdAt"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updatedAt"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deletedAt"`

	Room Room `json:"room" gorm:"foreignKey:RoomID;references:RoomID"`
}

func (UserRoom) TableName() string {
	return "USER_ROOM"
}
