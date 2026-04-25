package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type Postgresql struct {
	Host        string
	Port        string
	User        string
	Password    string
	Name        string
	MaxOpenConn int
	MaxIdleConn int
	MaxIdleTime time.Duration
	SSLMode     string
	Timeout     time.Duration
	logger      utils.LoggerStrategy
}

type Options func(*Postgresql)

func WithHost(host string) Options {
	return func(p *Postgresql) {
		p.Host = host
	}
}

func WithPort(port string) Options {
	return func(p *Postgresql) {
		p.Port = port
	}
}

func WithUser(user string) Options {
	return func(p *Postgresql) {
		p.User = user
	}
}

func WithPassword(password string) Options {
	return func(p *Postgresql) {
		p.Password = password
	}
}

func WithName(name string) Options {
	return func(p *Postgresql) {
		p.Name = name
	}
}

func WithMaxOpenConn(maxOpenConn int) Options {
	return func(p *Postgresql) {
		p.MaxOpenConn = maxOpenConn
	}
}

func WithMaxIdleConn(maxIdleConn int) Options {
	return func(p *Postgresql) {
		p.MaxIdleConn = maxIdleConn
	}
}

func WithMaxIdleTime(maxIdleTime time.Duration) Options {
	return func(p *Postgresql) {
		p.MaxIdleTime = maxIdleTime
	}
}

func WithSSLMode(mode string) Options {
	return func(p *Postgresql) {
		p.SSLMode = mode
	}
}

func WithTimeout(timeout time.Duration) Options {
	return func(p *Postgresql) {
		p.Timeout = timeout
	}
}

func WithLogger(logger utils.LoggerStrategy) Options {
	return func(p *Postgresql) {
		p.logger = logger
	}
}

func (p *Postgresql) uri() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC", p.Host, p.User, p.Password, p.Name, p.Port, p.SSLMode)
}

func (p *Postgresql) Connect() (*gorm.DB, *sql.DB, error) {
	db, err := gorm.Open(postgres.Open(p.uri()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		p.logger.Error("failed to connect database", "error", err.Error())
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		p.logger.Error("failed to get sql.db", "error", err.Error())
		return nil, nil, err
	}

	sqlDB.SetMaxOpenConns(p.MaxOpenConn)
	sqlDB.SetMaxIdleConns(p.MaxIdleConn)
	sqlDB.SetConnMaxIdleTime(p.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	if err = sqlDB.PingContext(ctx); err != nil {
		sqlDB.Close()
		p.logger.Error("failed to ping database", "error", err.Error())
	}

	return db, sqlDB, nil
}

func NewPostgresql(opts ...Options) *Postgresql {
	p := &Postgresql{}
	for _, o := range opts {
		o(p)
	}
	return p
}
