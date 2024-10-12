package repository

import (
	"github.com/getsentry/sentry-go"
	"github.com/kwa0x2/swiftchat-backend/models"
	"github.com/kwa0x2/swiftchat-backend/types"
	"gorm.io/gorm"
)

type IRequestRepository interface {
	Create(tx *gorm.DB, request *models.Request) error
	Update(tx *gorm.DB, whereRequest *models.Request, updateRequest *models.Request) error
	Delete(tx *gorm.DB, whereRequest *models.Request) error
	GetRequests(receiverEmail string) ([]*models.Request, error)
	GetSentRequests(senderEmail string) ([]*models.Request, error)
	InsertAndReturnUser(request *models.Request) (*models.Request, error)
	GetDB() *gorm.DB
}

type requestRepository struct {
	DB *gorm.DB
}

func NewRequestRepository(db *gorm.DB) IRequestRepository {
	return &requestRepository{
		DB: db,
	}
}

// region "Create" adds a new request to the database
func (r *requestRepository) Create(tx *gorm.DB, request *models.Request) error {
	db := r.DB
	if tx != nil {
		db = tx // Use the provided transaction if available
	}

	if err := db.Create(&request).Error; err != nil {
		sentry.CaptureException(err)
		return err
	}
	return nil
}

//endregion

// region "Update" modifies the fields of a request in the database based on specified conditions
func (r *requestRepository) Update(tx *gorm.DB, whereRequest *models.Request, updateRequest *models.Request) error {
	db := r.DB
	if tx != nil {
		db = tx
	}

	result := db.Model(&models.Request{}).Where(whereRequest).Updates(updateRequest)

	if result.Error != nil {
		sentry.CaptureException(result.Error)

		return result.Error // Return error if update fails
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Return an error if no records were updated
	}

	return nil
}

//endregion

// region "Delete" removes a request from the database
func (r *requestRepository) Delete(tx *gorm.DB, whereRequest *models.Request) error {
	db := r.DB
	if tx != nil {
		db = tx
	}

	result := db.Where(whereRequest).
		Delete(&models.Request{})

	if result.Error != nil {
		sentry.CaptureException(result.Error)
		return result.Error // Return error if deletion fails
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Return an error if no records were deleted
	}

	return nil
}

//endregion

// region "GetRequests" retrieves pending requests for a given receiver email.
func (r *requestRepository) GetRequests(receiverEmail string) ([]*models.Request, error) {
	var requests []*models.Request

	if err := r.DB.
		Where(&models.Request{ReceiverEmail: receiverEmail, RequestStatus: types.Pending}).
		Preload("User").
		Find(&requests).Error; err != nil {
		sentry.CaptureException(err)
		return nil, err
	}

	return requests, nil
}

// endregion

// region "GetSentRequests" retrieves sent requests for a given sender email
func (r *requestRepository) GetSentRequests(senderEmail string) ([]*models.Request, error) {
	var requests []*models.Request

	if err := r.DB.
		Where(&models.Request{SenderEmail: senderEmail, RequestStatus: types.Pending}).
		Find(&requests).Error; err != nil {
		sentry.CaptureException(err)
		return nil, err
	}

	return requests, nil
}

// endregion

// region "InsertAndReturnUser" creates a new request and returns the associated user.
func (r *requestRepository) InsertAndReturnUser(request *models.Request) (*models.Request, error) {
	var requestData *models.Request

	if err := r.DB.Create(request).Preload("User").Find(&requestData).Error; err != nil {
		sentry.CaptureException(err)
		return nil, err
	}

	return requestData, nil
}

// endregion

// region "GetDB" returns the underlying gorm.DB instance
func (r *requestRepository) GetDB() *gorm.DB {
	return r.DB // Return the database instance
}

// endregion
