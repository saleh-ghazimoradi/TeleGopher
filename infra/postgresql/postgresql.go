package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type Postgresql struct {
	Host           string
	Port           string
	User           string
	Password       string
	Name           string
	TimeZone       string
	MaxOpenConn    int
	MaxIdleConn    int
	MaxIdleTime    time.Duration
	MaxLifetime    time.Duration
	SSLMode        string
	ConnectTimeout time.Duration
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

func WithMaxOpenConn(maxConn int) Options {
	return func(p *Postgresql) {
		p.MaxOpenConn = maxConn
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

func WithConnectTimeout(timeout time.Duration) Options {
	return func(p *Postgresql) {
		p.ConnectTimeout = timeout
	}
}

func WithTimeZone(tz string) Options {
	return func(p *Postgresql) {
		p.TimeZone = tz
	}
}

func WithMaxLifetime(maxLifetime time.Duration) Options {
	return func(p *Postgresql) {
		p.MaxLifetime = maxLifetime
	}
}

func (p *Postgresql) uri() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=%s connect_timeout=%d", p.Host, p.Port, p.User, p.Password, p.Name, p.SSLMode, p.TimeZone, p.ConnectTimeout)
}

func (p *Postgresql) Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", p.uri())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(p.MaxOpenConn)
	db.SetMaxIdleConns(p.MaxIdleConn)
	db.SetConnMaxIdleTime(p.MaxIdleTime)
	db.SetConnMaxLifetime(p.MaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), p.ConnectTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			return nil, fmt.Errorf("ping failed: %w (and failed to close: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return db, nil
}

func NewPostgresql(opts ...Options) *Postgresql {
	p := &Postgresql{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}
