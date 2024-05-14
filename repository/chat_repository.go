package repository

import "gorm.io/gorm"

type ChatRepository struct {
	DB *gorm.DB
}