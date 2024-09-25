package models

import (
	"github.com/kwa0x2/realtime-chat-backend/types"
	"time"

	"gorm.io/gorm"
)

type Friend struct {
	FriendId     int64              `json:"friend_id" gorm:"primaryKey;not null"`
	UserMail     string             `json:"user_mail" gorm:"not null"`
	UserMail2    string             `json:"user_mail2" gorm:"not null"`
	CreatedAt    time.Time          `json:"createdAt" gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time          `json:"updatedAt" gorm:"column:updatedAt;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt     `json:"deletedAt" gorm:"column:deletedAt"`
	FriendStatus types.FriendStatus `json:"friend_status" gorm:"type:friend_status;not null;default:friend"`

	User User `json:"user" gorm:"foreignKey:UserMail;references:UserEmail"`
}

func (Friend) TableName() string {
	return "FRIEND"
}
