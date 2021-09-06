package mailer

import (
	"fmt"
	"net/smtp"

	"golang.org/x/oauth2"
)

// XOAUTH2 implementation for SMTP (https://developers.google.com/gmail/imap/xoauth2-protocol#smtp_protocol_exchange)
type XOAuth2 struct {
	sender string
	ts     oauth2.TokenSource
}

func (a *XOAuth2) Start(serverInfo *smtp.ServerInfo) (string, []byte, error) {
	if !serverInfo.TLS {
		return "", nil, fmt.Errorf("unencrypted connection: %v", serverInfo)
	}
	token, err := a.ts.Token()
	if err != nil {
		return "", nil, fmt.Errorf("failed to get token: %v", err)
	}
	toServer := fmt.Sprintf("user=%v\001auth=%v %v\001\001", a.sender, token.Type(), token.AccessToken)
	return "XOAUTH2", []byte(toServer), nil
}

func (a *XOAuth2) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		return nil, fmt.Errorf("unexpected challenge: %v", string(fromServer))
	}
	return nil, nil
}

func NewXOAuth2(sender string, ts oauth2.TokenSource) smtp.Auth {
	return &XOAuth2{sender: sender, ts: ts}
}
