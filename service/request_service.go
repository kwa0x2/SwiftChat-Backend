package service

import (
	"github.com/kwa0x2/swiftchat-backend/models"
	"github.com/kwa0x2/swiftchat-backend/repository"
	"github.com/kwa0x2/swiftchat-backend/types"
	"gorm.io/gorm"
)

type IRequestService interface {
	Create(tx *gorm.DB, request *models.Request) error
	Update(tx *gorm.DB, whereRequest *models.Request, updateRequest *models.Request) error
	DeleteByEmail(tx *gorm.DB, receiverEmail, senderEmail string) error
	GetRequests(receiverEmail string) ([]*models.Request, error)
	GetSentRequests(senderEmail string) ([]*models.Request, error)
	UpdateFriendshipRequest(receiverEmail, senderEmail string, requestStatus types.RequestStatus) (map[string]interface{}, error)
	InsertAndReturnUser(request *models.Request) (map[string]interface{}, error)
}

type requestService struct {
	RequestRepository repository.IRequestRepository
	FriendService     IFriendService
	UserService       IUserService
}

func NewRequestService(requestRepository repository.IRequestRepository, friendService IFriendService, userService IUserService) IRequestService {
	return &requestService{
		RequestRepository: requestRepository,
		FriendService:     friendService,
		UserService:       userService,
	}
}

// region "Create" adds a new request to the database
func (s *requestService) Create(tx *gorm.DB, request *models.Request) error {
	return s.RequestRepository.Create(tx, request)
}

//endregion

// region "Update" modifies the fields of a request in the database based on specified conditions
func (s *requestService) Update(tx *gorm.DB, whereRequest *models.Request, updateRequest *models.Request) error {
	return s.RequestRepository.Update(tx, whereRequest, updateRequest)
}

//endregion

// region "DeleteByEmail" removes a request based on the provided email information.
func (s *requestService) DeleteByEmail(tx *gorm.DB, receiverEmail, senderEmail string) error {
	whereRequest := &models.Request{
		ReceiverMail: receiverEmail,
		SenderMail:   senderEmail,
	}

	return s.RequestRepository.Delete(tx, whereRequest)
}

//endregion

// region "GetRequests" retrieves requests for a given receiver email
func (s *requestService) GetRequests(receiverEmail string) ([]*models.Request, error) {
	return s.RequestRepository.GetRequests(receiverEmail)
}

//endregion

// region "GetSentRequests" retrieves sent requests for a given sender email
func (s *requestService) GetSentRequests(senderEmail string) ([]*models.Request, error) {
	return s.RequestRepository.GetSentRequests(senderEmail)
}

//endregion

// region "UpdateFriendshipRequest" updates the status of a friendship request and manages friendship creation
func (s *requestService) UpdateFriendshipRequest(receiverEmail, senderEmail string, requestStatus types.RequestStatus) (map[string]interface{}, error) {
	tx := s.RequestRepository.GetDB().Begin() // Start a new database transaction
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Prepare the request that needs to be updated
	whereRequest := &models.Request{
		ReceiverMail: receiverEmail,
		SenderMail:   senderEmail,
	}

	// Create an updated request with the new status
	updateRequest := &models.Request{
		RequestStatus: requestStatus,
	}

	// Update the friendship request status in the database
	if err := s.Update(tx, whereRequest, updateRequest); err != nil {
		tx.Rollback() // Rollback the transaction on error
		return nil, err
	}

	// Delete the request from the database after handling the response
	if err := s.DeleteByEmail(tx, receiverEmail, senderEmail); err != nil {
		tx.Rollback() // Rollback the transaction on error
		return nil, err
	}

	// Retrieve user data of the receiver for the response
	userData, err := s.UserService.GetByEmail(receiverEmail)
	if err != nil {
		tx.Rollback() // Rollback the transaction on error
		return nil, err
	}

	var result map[string]interface{}
	if requestStatus == "accepted" {
		// If the request is accepted, create a friendship entry in the database
		friend := &models.Friend{
			UserMail:     senderEmail,
			UserMail2:    receiverEmail,
			FriendStatus: "friend",
		}

		if createErr := s.FriendService.Create(tx, friend); createErr != nil {
			tx.Rollback() // Rollback the transaction on error
			return nil, createErr
		}

		// Prepare the response data for accepted request
		result = map[string]interface{}{
			"status": "accepted",
			"user_data": map[string]interface{}{
				"friend_mail": receiverEmail,      // Friend's email
				"user_name":   userData.UserName,  // User's name
				"user_photo":  userData.UserPhoto, // User's photo
			},
		}

	} else {
		// Prepare the response data for declined request
		result = map[string]interface{}{
			"status": requestStatus, // Current status of the request
			"user_data": map[string]interface{}{
				"user_name": userData.UserName, // User's name
			},
		}
	}

	if commitErr := tx.Commit().Error; commitErr != nil {
		tx.Rollback()
		return nil, commitErr
	}

	return result, nil
}

//endregion

// region "InsertAndReturnUser" adds a new request and returns associated user data.
func (s *requestService) InsertAndReturnUser(request *models.Request) (map[string]interface{}, error) {
	tx := s.RequestRepository.GetDB().Begin() // Start a new database transaction
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Create the new request in the database
	if err := s.Create(tx, request); err != nil {
		tx.Rollback() // Rollback the transaction on error
		return nil, err
	}

	// Retrieve user data of the sender for the response
	userData, err := s.UserService.GetByEmail(request.SenderMail)
	if err != nil {
		tx.Rollback() // Rollback the transaction on error
		return nil, err
	}

	// Commit the transaction if everything was successful
	if commitErr := tx.Commit().Error; commitErr != nil {
		return nil, commitErr
	}

	// Prepare the response data with sender's information
	result := map[string]interface{}{
		"sender_mail": request.SenderMail, // Sender's email
		"user_name":   userData.UserName,  // User's name
		"user_photo":  userData.UserPhoto, // User's photo
	}

	return result, nil
}

// endregion
