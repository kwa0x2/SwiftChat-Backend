package service

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/repository"
)

type RequestService struct {
	RequestRepository *repository.RequestRepository
}

func (s *RequestService) Insert(request *models.Request) error {
	return s.RequestRepository.Insert(request)
}

func (s *RequestService) GetComingRequests(receiverMail string) ([]*models.Request, error) {
	return s.RequestRepository.GetComingRequests(receiverMail)
}

func (s *RequestService) Accept(request *models.Request) error {
	return s.RequestRepository.Accept(request)
}

func (s *RequestService) Reject(request *models.Request) error {
	return s.RequestRepository.Reject(request)
}
