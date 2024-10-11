package models

import (
	"github.com/kwa0x2/swiftchat-backend/types"
	"time"

	"gorm.io/gorm"
)

type Friend struct {
	FriendID     int64              `json:"friend_id" gorm:"primaryKey;not null;autoIncrement"`
	UserEmail    string             `json:"user_email" gorm:"not null"`
	UserEmail2   string             `json:"user_email2" gorm:"not null"`
	FriendStatus types.FriendStatus `json:"friend_status" gorm:"type:friend_status;not null;default:friend"`
	CreatedAt    time.Time          `json:"createdAt" gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time          `json:"updatedAt" gorm:"column:updatedAt;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt    *gorm.DeletedAt    `json:"deletedAt" gorm:"column:deletedAt"`

	User User `json:"user" gorm:"foreignKey:UserEmail;references:UserEmail"`
}

func (Friend) TableName() string {
	return "FRIEND"
}
