package repository

import "gorm.io/gorm"

type AuthRepository struct {
	DB *gorm.DB
}


func (r *AuthRepository) IsIdUnique(id string) bool {
	var count int64
	r.DB.Table("USER").Where("user_id = ?", id).Count(&count)
	return count == 0
}