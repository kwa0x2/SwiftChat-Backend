package models

import (
	"time"

	"github.com/kwa0x2/realtime-chat-backend/types"
	"gorm.io/gorm"
)


type Friend struct {
	UserMail  string       `json:"user_mail" gorm:"primaryKey;not null"`
	UserMail2 string       `json:"user_mail2" gorm:"primaryKey;not null"`
	Status    types.FriendStatus `json:"friend_status" gorm:"type:friend_status;not null;default:friend"`
	CreatedAt time.Time    `json:"createdAt" gorm:"column:createdAt;not null"`
	UpdatedAt time.Time    `json:"updatedAt" gorm:"column:updatedAt"`
	DeletedAt gorm.DeletedAt    `json:"deletedAt" gorm:"column:deletedAt"`
	User      User         `json:"user" gorm:"foreignKey:UserMail;references:UserEmail"`
}

func (Friend) TableName() string{
	return "FRIEND"
}