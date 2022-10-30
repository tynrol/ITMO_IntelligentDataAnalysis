package config

import (
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
	"time"
)

type App struct {
	Name  string `env:"APP_NAME" envDefault:"accessor-service"`
	Debug bool   `env:"APP_DEBUG" envDefault:"false"`
}

type HTTPConfig struct {
	Port          string        `env:"HTTP_SERVER_PORT" envDefault:":10000"`
	RemoteTimeout time.Duration `env:"HTTP_REMOTE_TIMEOUT" envDefault:"10s"`
}

type PathConfig struct {
	DatasetsPath string `env:"DATASETS_PATH" envDefault:"./datasets/dataset"`
	DBPath       string `env:"DB_PATH" envDefault:"./init/db/database.db"`
}

type Token struct {
	UplashToken string `env:"UNSPLASH_TOKEN"`
}

type Config struct {
	App
	HTTPConfig
	PathConfig
	Token
}

func (c *Config) Parse() (err error) {
	const op = "config_Parse"

	if err = env.Parse(&c.App); err != nil {
		return errors.Wrap(err, op)
	}

	if err = env.Parse(&c.HTTPConfig); err != nil {
		return errors.Wrap(err, op)
	}

	if err = env.Parse(&c.PathConfig); err != nil {
		return errors.Wrap(err, op)
	}

	if err = env.Parse(&c.Token); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
