package models

import (
	"github.com/kwa0x2/realtime-chat-backend/types"
	"time"
)

type Request struct {
	Id           int64         `json:"id" gorm:"primaryKey;not null"`
	SenderMail   string        `json:"sender_mail" gorm:"not null"`
	ReceiverMail string        `json:"receiver_mail" gorm:"not null"`
	Status       types.RequestStatus `json:"status" gorm:"not null;type:request_status;default:pending"`
	CreatedAt    time.Time     `json:"createdAt" gorm:"column:createdAt;not null"`
	User         User          `json:"user" gorm:"foreignKey:SenderMail;references:UserEmail"`
}

func (Request) TableName() string {
	return "REQUEST"
}
