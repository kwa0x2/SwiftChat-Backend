package models

import (
	"github.com/kwa0x2/swiftchat-backend/types"
	"gorm.io/gorm"
	"time"
)

type Request struct {
	RequestID     int64               `json:"request_id" gorm:"primaryKey;not null"`
	SenderEmail   string              `json:"sender_email" gorm:"not null"`
	ReceiverEmail string              `json:"receiver_email" gorm:"not null"`
	RequestStatus types.RequestStatus `json:"request_status" gorm:"not null;type:request_status;default:pending"`
	CreatedAt     time.Time           `json:"createdAt" gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt     gorm.DeletedAt      `json:"deletedAt" gorm:"column:deletedAt"`

	User User `json:"user" gorm:"foreignKey:SenderEmail;references:UserEmail"`
}

func (Request) TableName() string {
	return "REQUEST"
}
