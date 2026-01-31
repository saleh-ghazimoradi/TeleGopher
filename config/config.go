package config

import (
	"github.com/caarlos0/env/v11"
	"sync"
	"time"
)

var (
	instance *Config
	once     sync.Once
	initErr  error
)

type Config struct {
	Server      Server
	Postgresql  Postgresql
	Application Application
}

type Server struct {
	Host         string        `env:"SERVER_HOST"`
	Port         string        `env:"SERVER_PORT"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT"`
}

type Postgresql struct {
	Host              string        `env:"POSTGRES_HOST"`
	Port              string        `env:"POSTGRES_PORT"`
	User              string        `env:"POSTGRES_USER"`
	Password          string        `env:"POSTGRES_PASSWORD"`
	Name              string        `env:"POSTGRES_NAME"`
	TimeZone          string        `env:"POSTGRES_TIMEZONE"`
	MaxOpenConn       int           `env:"POSTGRES_MAX_OPEN_CONN"`
	MaxIdleConn       int           `env:"POSTGRES_MAX_IDLE_CONN"`
	MaxIdleTime       time.Duration `env:"POSTGRES_MAX_IDLE_TIME"`
	MaxLifetime       time.Duration `env:"POSTGRES_MAX_LIFETIME"`
	SSLMode           string        `env:"POSTGRES_SSL_MODE"`
	ConnectionTimeout time.Duration `env:"POSTGRES_CONNECTION_TIMEOUT"`
}

type Application struct {
	Version     string `env:"VERSION"`
	Environment string `env:"ENVIRONMENT"`
}

func GetConfigInstance() (*Config, error) {
	once.Do(func() {
		instance = &Config{}
		initErr = env.Parse(instance)
		if initErr != nil {
			instance = nil
		}
	})
	return instance, initErr
}
