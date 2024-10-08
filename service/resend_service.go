package service

import (
	"fmt"
	"github.com/kwa0x2/swiftchat-backend/utils"
	"github.com/resend/resend-go/v2"
)

type IResendService interface {
	SendEmail(to, subject, template string) (string, error)
}

type resendService struct {
	ResendClient *resend.Client
}

func NewResendService(resendClient *resend.Client) IResendService {
	return &resendService{
		ResendClient: resendClient,
	}
}

// region "EmailData" holds email-related information.
type EmailData struct {
	Email string
}

// endregion

// region "SendEmail" sends an email using the Resend client.
func (s *resendService) SendEmail(to, subject, template string) (string, error) {
	// Load the HTML template and replace placeholders with the recipient's information.
	htmlContent, err := utils.LoadTemplate(template, to)
	if err != nil {
		return "", err
	}

	// Prepare the parameters for sending the email.
	params := &resend.SendEmailRequest{
		From:    "NettaSec Solutions <admin@nettasec.com>", // Sender's email information.
		To:      []string{to},                              // Recipient's email address in a slice.
		Html:    string(htmlContent),                       // HTML content of the email.
		Subject: subject,                                   // Subject of the email.
	}

	sent, sentErr := s.ResendClient.Emails.Send(params)
	if sentErr != nil {
		return "", fmt.Errorf("error sending email: %w", sentErr)
	}
	return sent.Id, nil
}

// endregion
