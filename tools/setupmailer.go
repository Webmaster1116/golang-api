package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"
)

var (
	sender     string
	configPath string
	tokenPath  string
)

func init() {
	flag.StringVar(&sender, "sender", "", "Specifies the sender's email address.")
	flag.StringVar(&configPath, "c", "", "Config path")
	flag.StringVar(&tokenPath, "t", "", "Token path")
}

func main() {
	flag.Parse()
	// if atDomain := "@gmail.com"; !strings.HasSuffix(sender, atDomain) {
	// 	log.Fatalf("-sender must specify an %v email address.", atDomain)
	// }
	config := getConfig()
	setUpToken(config, tokenPath)
}

func getConfig() *oauth2.Config {
	configJSON, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v.", err)
	}
	config, err := googleOAuth2.ConfigFromJSON(configJSON, "https://mail.google.com/")
	if err != nil {
		log.Fatalf("Failed to parse config: %v.", err)
	}
	return config
}

func setUpToken(config *oauth2.Config, tokenPath string) {
	fmt.Println()
	fmt.Println("1. Ensure that you are logged in as", sender, "in your browser.")
	fmt.Println()
	fmt.Println("2. Open the following link and authorise sendgmail:")
	fmt.Println(config.AuthCodeURL("state", oauth2.AccessTypeOffline))
	fmt.Println()
	fmt.Println("3. Enter the authorisation code:")
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Failed to read authorisation code: %v.", err)
	}
	fmt.Println()
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Failed to exchange authorisation code for token: %v.", err)
	}
	tokenFile, err := os.OpenFile(tokenPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Failed to open token file for writing: %v.", err)
	}
	defer tokenFile.Close()
	if err := json.NewEncoder(tokenFile).Encode(token); err != nil {
		log.Fatalf("Failed to write token: %v.", err)
	}
}
