package usecases

import "github.com/Leonardo-Henrique/decoreagora/app/core/ports"

type EmailUseCase struct {
	emailHandler ports.EmailHandler
}

func NewEmailUseCase(emailHandler ports.EmailHandler) *EmailUseCase {
	return &EmailUseCase{
		emailHandler: emailHandler,
	}
}

func (e *EmailUseCase) SendEmail(to, subject, body string) error {
	return e.emailHandler.SendWithoutPreStyle(to, subject, body)
}
