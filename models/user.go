package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserID    string         `json:"user_id" gorm:"primaryKey;not null"`
	UserEmail string         `json:"user_email" gorm:"type:varchar(50);not null"`
	UserName  string         `json:"user_name" gorm:"type:varchar(10);not null"`
	UserPhoto string         `json:"user_photo" gorm:"not null"`
	CreatedAt time.Time      `json:"createdAt" gorm:"not null;column:createdAt;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"not null;column:updatedAt;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"column:deletedAt"`

	Friend *Friend `json:"friend" gorm:"foreignKey:UserEmail;references:UserEmail"`
}

func (User) TableName() string {
	return "USER"
}
