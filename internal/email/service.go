package email

import (
	"context"
	"github.com/resend/resend-go/v3"
)

type Service interface {
	Send(ctx context.Context, to, from, subject, body string) error
}

type service struct {
	resend *resend.Client
}

func (s *service) Send(_ context.Context, to, from, subject, body string) error {
	params := &resend.SendEmailRequest{ // <- note: from the package, not s.resend
		To:      []string{to},
		From:    from,
		Subject: subject,
		Html:    body,
	}

	_, err := s.resend.Emails.Send(params)
	return err
}
