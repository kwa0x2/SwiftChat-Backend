package service

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/repository"
	"gorm.io/gorm"
)

type RequestService struct {
	RequestRepository *repository.RequestRepository
	FriendService     *FriendService
	UserService       *UserService
}

// region INSERT NEW REQUEST SERVICE
func (s *RequestService) Insert(tx *gorm.DB, request *models.Request) error {
	return s.RequestRepository.Insert(tx, request)
}

//endregion

// region GET COMING REQUEST BY RECEIVER EMAIL SERVICE
func (s *RequestService) GetComingRequests(receiverMail string) ([]*models.Request, error) {
	return s.RequestRepository.GetComingRequests(receiverMail)
}

//endregion

// region UPDATE BY MAIL SERVICE
func (s *RequestService) Update(tx *gorm.DB, request *models.Request) error {
	return s.RequestRepository.Update(tx, request)
}

//endregion

// region DELETE BY MAIL SERVICE
func (s *RequestService) Delete(tx *gorm.DB, request *models.Request) error {
	return s.RequestRepository.Delete(tx, request)
}

//endregion

// region UPDATE REQUEST STATUS and DELETE AND IF STATUS ACCEPTED INSERT NEW FRIENDSHIP IN FRIEND WITH TRANSACTION SERVICE
func (s *RequestService) UpdateFriendshipRequest(request *models.Request) (bool, error) {
	tx := s.RequestRepository.DB.Begin()
	if tx.Error != nil {
		return false, tx.Error
	}

	if err := s.Update(tx, request); err != nil {
		tx.Rollback()
		return false, err
	}

	if err := s.Delete(tx, request); err != nil {
		tx.Rollback()
		return false, err
	}

	if request.RequestStatus == "accepted" {
		friend := &models.Friend{
			UserMail:     request.SenderMail,
			UserMail2:    request.ReceiverMail,
			FriendStatus: "friend",
		}

		if err := s.FriendService.Insert(tx, friend); err != nil {
			tx.Rollback()
			return false, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, err
	}

	return true, nil
}

//endregion

func (s *RequestService) InsertAndReturnUser(request *models.Request) (map[string]interface{}, error) {
	tx := s.RequestRepository.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := s.Insert(tx, request); err != nil {
		tx.Rollback()
		return nil, err
	}

	userData, err := s.UserService.GetByEmail(request.SenderMail)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"sender_mail": request.SenderMail,
		"user_name":   userData.UserName,
		"user_photo":  userData.UserPhoto,
	}

	return result, nil
}
