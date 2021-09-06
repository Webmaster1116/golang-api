package config

type AUTH struct {
	JWT_SECRET []byte
}

var auth_configs = Configs{
	"JWT_SECRET": func(i interface{}, s string) error {
		i.(*AUTH).JWT_SECRET = []byte(s)
		return nil
	},
}

func (*AUTH) Configs() Configs { return auth_configs }
