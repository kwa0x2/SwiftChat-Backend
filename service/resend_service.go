package service

import (
	"fmt"
	"github.com/kwa0x2/realtime-chat-backend/utils"
	"github.com/resend/resend-go/v2"
)

type ResendService struct {
	ResendClient *resend.Client
}
type EmailData struct {
	Email string
}

func (s *ResendService) SendMail(to, subject, template string) (string, error) {
	htmlContent, err := utils.LoadTemplate(template, to)
	if err != nil {
		return "", err
	}

	params := &resend.SendEmailRequest{
		From:    "NettaSec Solutions <admin@nettasec.com>",
		To:      []string{to},
		Html:    string(htmlContent),
		Subject: subject,
	}

	sent, err := s.ResendClient.Emails.Send(params)
	if err != nil {
		return "", fmt.Errorf("error sending email: %w", err)
	}
	return sent.Id, nil
}
