package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// config consumer
type Consumer func(interface{}, string) error

// set pointer value
func Consumer_Set(pv func(interface{}) *string) Consumer {
	return func(i interface{}, v string) error { *pv(i) = v; return nil }
}

// A config
type Configs map[string]Consumer
type IConfig interface {
	Configs() Configs
}

type AppConfig struct {
	DB
	AUTH
	MAILER
	LOG
	WEBPUSH
}

var Config AppConfig

func load_vars(c IConfig) error {
	for k, consume := range c.Configs() {
		v, found := os.LookupEnv(k)
		if !found {
			return fmt.Errorf("variable '%s' isn't set", k)
		}
		if err := consume(c, v); err != nil {
			return errors.Wrapf(err, "Failed to set %s to '%s'", k, v)
		}
	}
	return nil
}

func Init() error {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("no configuration file found: ", err.Error())
	}
	configValue := reflect.ValueOf(&Config).Elem()

	for idx := 0; idx < configValue.NumField(); idx++ {
		config := configValue.Field(idx).Addr().Interface().(IConfig)
		if err := load_vars(config); err != nil {
			return err
		}
	}
	return nil
}
