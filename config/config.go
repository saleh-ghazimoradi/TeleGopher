package config

import (
	"github.com/caarlos0/env/v11"
	"sync"
	"time"
)

var (
	cfg     *Config
	once    sync.Once
	initErr error
)

type Config struct {
	Application Application
	JWT         JWT
	Postgresql  Postgresql
	Server      Server
}

type Application struct {
	Version     string `env:"APP_VERSION"`
	Environment string `env:"APP_ENVIRONMENT"`
}

type Postgresql struct {
	Host        string        `env:"POSTGRES_HOST"`
	Port        string        `env:"POSTGRES_PORT"`
	User        string        `env:"POSTGRES_USER"`
	Password    string        `env:"POSTGRES_PASSWORD"`
	Name        string        `env:"POSTGRES_NAME"`
	MaxOpenConn int           `env:"POSTGRES_MAX_OPEN_CONN"`
	MaxIdleConn int           `env:"POSTGRES_MAX_IDLE_CONN"`
	MaxIdleTime time.Duration `env:"POSTGRES_MAX_IDLE_TIME"`
	SSLMode     string        `env:"POSTGRES_SSL_MODE"`
	Timeout     time.Duration `env:"POSTGRES_TIMEOUT"`
}

type JWT struct {
	Secret              string        `env:"JWT_SECRET"`
	ExpiresIn           time.Duration `env:"JWT_EXPIRES_IN"`
	RefreshTokenExpires time.Duration `env:"JWT_REFRESH_TOKEN_EXPIRES"`
}

type Server struct {
	Host         string        `env:"SERVER_HOST"`
	Port         string        `env:"SERVER_PORT"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT"`
}

func GetCfg() (*Config, error) {
	once.Do(func() {
		cfg = &Config{}
		initErr = env.Parse(cfg)
		if initErr != nil {
			cfg = nil
		}
	})
	return cfg, initErr
}
