package models

import "time"

type User struct {
	UserID    string    `json:"user_id"`
	UserEmail string    `json:"user_email"`
	UserName  string    `json:"user_name"`
	UserPhoto string    `json:"user_photo"`
	UserRole  string    `json:"user_role" gorm:"default:standard"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updatedAt"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deletedAt"`
}
