package repository

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"gorm.io/gorm"
)

type RequestRepository struct {
	DB *gorm.DB
}

// region INSERT NEW REQUEST REPOSITORY
func (r *RequestRepository) Insert(request *models.Request) error {
	if err := r.DB.Create(&request).Error; err != nil {
		return err
	}
	return nil
}

//endregion

// region GET COMING REQUEST BY RECEIVER EMAIL REPOSITORY
func (r *RequestRepository) GetComingRequests(receiverMail string) ([]*models.Request, error) {
	var requests []*models.Request

	if err := r.DB.Debug().
		Where("receiver_mail = ? AND request_status = ?", receiverMail, "pending").
		Preload("User").
		Find(&requests).Error; err != nil {
		return nil, err
	}

	return requests, nil
}

//endregion

// region UPDATE REQUEST STATUS AND DELETE REPOSITORY
func (r *RequestRepository) UpdateStatusAndDelete(request *models.Request) (bool, error) {
	tx := r.DB.Begin()

	if tx.Error != nil {
		return false, tx.Error
	}

	updateResult := tx.Debug().Model(&models.Request{}).
		Where("receiver_mail = ? AND sender_mail = ?", request.ReceiverMail, request.SenderMail).
		Update("request_status", request.RequestStatus)
	if updateResult.Error != nil {
		tx.Rollback()
		return false, updateResult.Error
	}
	if updateResult.RowsAffected == 0 {
		tx.Rollback()
		return false, gorm.ErrRecordNotFound
	}

	deleteResult := tx.Debug().
		Where("receiver_mail = ? AND sender_mail = ?", request.ReceiverMail, request.SenderMail).
		Delete(&models.Request{})
	if deleteResult.Error != nil {
		tx.Rollback()
		return false, deleteResult.Error
	}
	if deleteResult.RowsAffected == 0 {
		tx.Rollback()
		return false, gorm.ErrRecordNotFound
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, err
	}

	return true, nil
}

//endregion

// region UPDATE FRIENDSHIP REQUEST REPOSITORY
func (r *RequestRepository) UpdateFriendshipRequest(request *models.Request) (bool, error) {
	tx := r.DB.Begin()

	if tx.Error != nil {
		return false, tx.Error
	}

	success, err := r.UpdateStatusAndDelete(request)
	if !success || err != nil {
		return false, err
	}

	if request.RequestStatus == "accepted" {
		friend := models.Friend{
			UserMail:     request.SenderMail,
			UserMail2:    request.ReceiverMail,
			FriendStatus: "friend",
		}

		if err := tx.Create(&friend).Error; err != nil {
			tx.Rollback()
			return false, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, err
	}

	return true, err
}

//endregion
