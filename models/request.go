package models

import (
	"github.com/kwa0x2/realtime-chat-backend/types"
	"gorm.io/gorm"
	"time"
)

type Request struct {
	RequestId     int64               `json:"request_id" gorm:"primaryKey;not null"`
	SenderMail    string              `json:"sender_mail" gorm:"not null"`
	ReceiverMail  string              `json:"receiver_mail" gorm:"not null"`
	RequestStatus types.RequestStatus `json:"request_status" gorm:"not null;type:request_status;default:pending"`
	CreatedAt     time.Time           `json:"createdAt" gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt     gorm.DeletedAt      `json:"deletedAt" gorm:"column:deletedAt"`

	User User `json:"user" gorm:"foreignKey:SenderMail;references:UserEmail"`
}

func (Request) TableName() string {
	return "REQUEST"
}
