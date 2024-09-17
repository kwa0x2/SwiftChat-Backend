package repository

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"gorm.io/gorm"
)

type RequestRepository struct {
	DB *gorm.DB
}

// region INSERT NEW REQUEST REPOSITORY
func (r *RequestRepository) Insert(tx *gorm.DB, request *models.Request) error {
	db := r.DB
	if tx != nil {
		db = tx
	}

	if err := db.Create(&request).Error; err != nil {
		return err
	}
	return nil
}

//endregion

// region GET COMING REQUEST BY RECEIVER EMAIL REPOSITORY
func (r *RequestRepository) GetComingRequests(receiverMail string) ([]*models.Request, error) {
	var requests []*models.Request

	if err := r.DB.
		Where("receiver_mail = ? AND request_status = ?", receiverMail, "pending").
		Preload("User").
		Find(&requests).Error; err != nil {
		return nil, err
	}

	return requests, nil
}

//endregion

// region UPDATE BY MAIL REPOSITORY
func (r *RequestRepository) Update(tx *gorm.DB, request *models.Request) error {
	db := r.DB
	if tx != nil {
		db = tx
	}

	result := db.Model(&models.Request{}).
		Where("receiver_mail = ? AND sender_mail = ?", request.ReceiverMail, request.SenderMail).
		Update("request_status", request.RequestStatus)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

//endregion

// region DELETE BY MAIL REPOSITORY
func (r *RequestRepository) Delete(tx *gorm.DB, request *models.Request) error {
	db := r.DB
	if tx != nil {
		db = tx
	}

	result := db.Model(&models.Request{}).
		Where("receiver_mail = ? AND sender_mail = ?", request.ReceiverMail, request.SenderMail).
		Delete(&models.Request{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

//endregion

func (r *RequestRepository) InsertAndReturnUser(request *models.Request) (*models.Request, error) {
	var requestData *models.Request

	if err := r.DB.Create(request).Preload("User").Find(&requestData).Error; err != nil {
		return nil, err
	}

	return requestData, nil
}
