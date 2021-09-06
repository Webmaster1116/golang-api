package webpush

import "github.com/SherClockHolmes/webpush-go"

type SherClockPush struct {
	VapidDetails
}

func (s *SherClockPush) SetVapidDetails(d *VapidDetails) {
	s.VapidDetails = *d
}

func (s *SherClockPush) Notify(subscription *Subscription, msg string, options *Options) error {
	_, err := webpush.SendNotification(
		[]byte(msg),
		&webpush.Subscription{
			Endpoint: subscription.Endpoint,
			Keys: webpush.Keys{
				Auth:   subscription.Keys.Auth,
				P256dh: subscription.Keys.P256dh,
			},
		},
		&webpush.Options{
			Topic:           s.Subject,
			TTL:             int(options.TTL),
			VAPIDPublicKey:  s.PublicKey,
			VAPIDPrivateKey: s.PrivateKey,
		},
	)
	return err
}

func NewSherClockPush(d *VapidDetails) *SherClockPush {
	return &SherClockPush{VapidDetails: *d}
}
