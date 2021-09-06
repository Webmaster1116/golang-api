package config

type WEBPUSH struct {
	Subject    string `json:"subject"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

var webpush_configs = Configs{
	"WEBPUSH_SUBJECT":    Consumer_Set(func(i interface{}) *string { return &(i).(*WEBPUSH).Subject }),
	"WEBPUSH_PUBLICKEY":  Consumer_Set(func(i interface{}) *string { return &(i).(*WEBPUSH).PublicKey }),
	"WEBPUSH_PRIVATEKEY": Consumer_Set(func(i interface{}) *string { return &(i).(*WEBPUSH).PrivateKey }),
}

func (*WEBPUSH) Configs() Configs { return webpush_configs }
