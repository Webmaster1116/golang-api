package mailer

import (
	"github.com/med8bra/moni-api-go/config"
)

// A mailer service
type Mailer interface {
	Send(*Email) error
}

func Impl(c *config.MAILER) (Mailer, error) {
	return NewGmailMailer(c.Sender, c.ConfigPath, c.TokenPath)
}
