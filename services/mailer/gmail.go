package mailer

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"golang.org/x/oauth2"
	goauth2 "golang.org/x/oauth2/google"
)

func loadGmailConfig(path string) (*oauth2.Config, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	configData, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}
	return goauth2.ConfigFromJSON(configData)
}

func loadToken(path string) (*oauth2.Token, error) {
	tokenFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer tokenFile.Close()
	var token oauth2.Token
	if err := json.NewDecoder(tokenFile).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

func NewGmailMailer(sender string, configPath string, tokenPath string) (Mailer, error) {
	ctx := context.Background()
	// load gmail config
	config, err := loadGmailConfig(configPath)
	if err != nil {
		return nil, err
	}
	// load token
	token, err := loadToken(tokenPath)
	if err != nil {
		return nil, err
	}

	tokenSource := config.TokenSource(ctx, token)
	auth := NewXOAuth2(sender, tokenSource)
	return NewSMTP("smtp.gmail.com:587", sender, auth), nil
}
