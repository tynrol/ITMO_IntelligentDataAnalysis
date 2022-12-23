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
	Port          string        `env:"HTTP_SERVER_PORT" envDefault:":10001"`
	RemoteTimeout time.Duration `env:"HTTP_REMOTE_TIMEOUT" envDefault:"10s"`
}

type PathConfig struct {
	ModelPath string `env:"MODEL_FILE_PATH" envDefault:"/home/tynrol/Code/GoLandProjects/ITMO_IntelligentDataAnalysis/detection-service/models/model.pt"`
}

type Metrics struct {
	Namespace string `env:"HTTP_METRICS_NAMESPACE" envDefault:"local"`
	Port      uint   `env:"HTTP_METRICS_PORT" envDefault:"9090"`
}

type Config struct {
	App
	HTTPConfig
	PathConfig
	Metrics
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

	if err = env.Parse(&c.Metrics); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
