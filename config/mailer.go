package config

type MAILER struct {
	Sender     string
	ConfigPath string
	TokenPath  string
}

var mailer_configs = Configs{
	"MAILER_SENDER": Consumer_Set(func(i interface{}) *string { return &i.(*MAILER).Sender }),
	"MAILER_CONFIG": Consumer_Set(func(i interface{}) *string { return &i.(*MAILER).ConfigPath }),
	"MAILER_TOKEN":  Consumer_Set(func(i interface{}) *string { return &i.(*MAILER).TokenPath }),
}

func (*MAILER) Configs() Configs { return mailer_configs }
