package service

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/repository"
)

type RequestService struct {
	RequestRepository *repository.RequestRepository
}

// region INSERT NEW REQUEST SERVICE
func (s *RequestService) Insert(request *models.Request) error {
	return s.RequestRepository.Insert(request)
}

//endregion

// region GET COMING REQUEST BY RECEIVER EMAIL SERVICE
func (s *RequestService) GetComingRequests(receiverMail string) ([]*models.Request, error) {
	return s.RequestRepository.GetComingRequests(receiverMail)
}

//endregion V

// region UPDATE REQUEST STATUS AND DELETE SERVICE
func (s *RequestService) UpdateStatusAndDelete(request *models.Request) (bool, error) {
	return s.RequestRepository.UpdateStatusAndDelete(request)
}

//endregion

// region UPDATE FRIENDSHIP REQUEST SERVICE
func (s *RequestService) UpdateFriendshipRequest(request *models.Request) (bool, error) {
	return s.RequestRepository.UpdateFriendshipRequest(request)
}

//endregion
