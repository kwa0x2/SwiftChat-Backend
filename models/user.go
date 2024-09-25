package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserID    string         `json:"user_id" gorm:"primaryKey"`
	UserEmail string         `json:"user_email"`
	UserName  string         `json:"user_name"`
	UserPhoto string         `json:"user_photo"`
	UserRole  string         `json:"user_role" gorm:"default:standard"`
	CreatedAt time.Time      `json:"createdAt" gorm:"not null;column:createdAt;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"not null;column:updatedAt;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"column:deletedAt"`

	Friend *Friend `json:"friend" gorm:"foreignKey:UserMail;references:UserEmail"`
}

func (User) TableName() string {
	return "USER"
}
