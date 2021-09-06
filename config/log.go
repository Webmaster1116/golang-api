package config

import "github.com/sirupsen/logrus"

type LOG struct {
	Level logrus.Level
}

var log_configs = Configs{
	"LOG_LEVEL": func(i interface{}, s string) error {
		if l, err := logrus.ParseLevel(s); err != nil {
			return err
		} else {
			i.(*LOG).Level = l
			return err
		}
	},
}

func (*LOG) Configs() Configs { return log_configs }
