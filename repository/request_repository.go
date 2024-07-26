package repository

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"gorm.io/gorm"
)

type RequestRepository struct {
	DB *gorm.DB
}

func (r *RequestRepository) Insert(request *models.Request) error {
	if err := r.DB.Create(&request).Error; err != nil {
		return err
	}
	return nil
}


func (r *RequestRepository) GetComingRequests(receiverMail string) ([]*models.Request, error) {
	var requests []*models.Request

	if err := r.DB.
		Where("receiver_mail = ? AND status = ?", receiverMail, "pending").
		Preload("User").
		Find(&requests).Error; err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *RequestRepository) Accept(request *models.Request) error {
	if err := r.DB.
		Where("receiver_mail = ? AND sender_mail = ?", request.ReceiverMail, request.SenderMail).
		Update("status", "accepted").Error; err != nil {
		return err
	}

	return nil
}

func (r *RequestRepository) Reject(request *models.Request) error {
	if err := r.DB.
		Where("receiver_mail = ? AND sender_mail = ?", request.ReceiverMail, request.SenderMail).
		Update("status", "rejected").Error; err != nil {
		return err
	}

	return nil
}
