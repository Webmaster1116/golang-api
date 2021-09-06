package auth

import (
	"time"

	"github.com/med8bra/moni-api-go/config"
	"github.com/med8bra/moni-api-go/models"
)

type Authenticator interface {
	Verify(token string) (user *models.AuthUser, err error)
	Sign(ID uint, duration time.Duration) (token string, err error)
}

func Impl(c *config.AUTH) (Authenticator, error) {
	return NewJwtAuthenticator(c.JWT_SECRET), nil
}
