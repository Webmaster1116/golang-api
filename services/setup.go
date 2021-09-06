package services

import (
	"github.com/med8bra/moni-api-go/config"
	"github.com/med8bra/moni-api-go/services/auth"
	"github.com/med8bra/moni-api-go/services/mailer"
	"github.com/med8bra/moni-api-go/services/template"
	"github.com/med8bra/moni-api-go/services/webpush"
	"github.com/pkg/errors"
)

var Authenticator auth.Authenticator
var Mailer mailer.Mailer
var TemplateManager template.TemplateManager
var WebPush webpush.WebPush

func Init(c *config.AppConfig) error {
	// setup services
	// -- Auth
	if authenticator, err := auth.Impl(&c.AUTH); err != nil {
		return errors.Wrap(err, "failed to create authenticator service")
	} else {
		Authenticator = authenticator
	}

	// -- Mailer
	if mailer, err := mailer.Impl(&c.MAILER); err != nil {
		return errors.Wrap(err, "failed to create mailer service")
	} else {
		Mailer = mailer
	}

	// -- template Manager
	TemplateManager = template.Impl()

	// -- web push service
	if webpush, err := webpush.Impl(&c.WEBPUSH); err != nil {
		return errors.Wrap(err, "failed to create webpush service")
	} else {
		WebPush = webpush
	}

	return nil
}
