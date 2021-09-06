package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/med8bra/moni-api-go/models"
	"github.com/sirupsen/logrus"
)

type TokenClaims struct {
	User models.AuthUser `json:"user"`
	jwt.StandardClaims
}

type JwtAuthenticator struct {
	secret []byte
}

func NewJwtAuthenticator(secret []byte) *JwtAuthenticator {
	return &JwtAuthenticator{secret}
}

func (j *JwtAuthenticator) validate(token *jwt.Token) ([]byte, error) {
	if _, valid := token.Method.(*jwt.SigningMethodHMAC); !valid {
		return nil, fmt.Errorf("invalid token : %s", token.Header["alg"])
	}
	return j.secret, nil
}

func (j *JwtAuthenticator) Verify(tokenStr string) (user *models.AuthUser, err error) {
	// verify
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.validate(token)
	})

	if err != nil || !token.Valid {
		return nil, ErrTokenInvalid
	}
	claims, valid := token.Claims.(*TokenClaims)
	if !valid {
		logrus.Warn("invalid claims: ", token.Claims)
		return nil, ErrTokenClaimsInvalid
	}
	return &claims.User, nil
}

func (j *JwtAuthenticator) Sign(ID uint, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		TokenClaims{
			User: models.AuthUser{ID: ID},
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(duration).Unix(),
				Issuer:    "moni",
			},
		})
	return token.SignedString(j.secret)
}
