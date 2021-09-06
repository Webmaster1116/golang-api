package webpush

import (
	"time"

	"github.com/med8bra/moni-api-go/config"
)

// VapidDetails are the global config of Push API
type VapidDetails struct {
	Subject    string `json:"subject"`    // mailto Address or URL
	PublicKey  string `json:"publicKey"`  // URL Safe Base64 Encoded Public Key
	PrivateKey string `json:"privateKey"` // URL Safe Base64 Encoded Private Key
}

// Keys are the base64 encoded values from PushSubscription.getKey()
type Keys struct {
	Auth   string `json:"auth"`
	P256dh string `json:"p256dh"`
}

// Subscription represents a PushSubscription object from the Push API
type Subscription struct {
	Endpoint string `json:"endpoint"`
	Keys     Keys   `json:"keys"`
}

// Notification options
type Options struct {
	GcmAPIKey string // GCM API Key
	Timeout   time.Duration
	TTL       time.Duration
	Headers   map[string]string
}

// A WebPush service
type WebPush interface {
	SetVapidDetails(*VapidDetails)
	Notify(subscription *Subscription, message string, options *Options) error
}

func Impl(c *config.WEBPUSH) (WebPush, error) {
	return NewSherClockPush(&VapidDetails{
		Subject:    c.Subject,
		PublicKey:  c.PublicKey,
		PrivateKey: c.PrivateKey,
	}), nil
}
