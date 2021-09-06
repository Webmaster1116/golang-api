package mailer

import (
	"bytes"
	"net/smtp"

	"github.com/pkg/errors"
)

// SMTP mailer
type SMTP struct {
	sender string
	addr   string
	auth   smtp.Auth
}

func (s *SMTP) Send(email *Email) error {
	var content bytes.Buffer
	// -- build email
	if err := email.Bytes(&content); err != nil {
		return errors.Wrap(err, "failed to build email")
	}
	return smtp.SendMail(s.addr, s.auth, s.sender, email.Tos, content.Bytes())
}

func NewSMTP(addr string, sender string, auth smtp.Auth) *SMTP {
	return &SMTP{sender: sender, addr: addr, auth: auth}
}
