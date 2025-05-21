package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env       string `envconfig:"ENV" default:"local"` // local, dev, prod
	Server    Server
	Store     Store
	Token     Token
	Scheduler Scheduler
}
type Scheduler struct {
	TimeoutRemoveRefreshTokens time.Duration `envconfig:"SCHEDULER_TIMEOUT_REMOVE_REFRESH_TOKENS" required:"true"`
}
type Server struct {
	BindAddr     string        `envconfig:"SERVER_BIND_ADDR" required:"true"`
	Name         string        `envconfig:"SERVER_NAME" required:"true"`
	WriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" required:"true"`
}

type Store struct {
	Host                string        `envconfig:"STORE_HOST" required:"true"`
	Port                int           `envconfig:"STORE_PORT" required:"true"`
	Name                string        `envconfig:"STORE_NAME" required:"true"`
	User                string        `envconfig:"STORE_USER" required:"true"`
	Password            string        `envconfig:"STORE_PASSWORD" required:"true"`
	SSLMode             string        `envconfig:"STORE_SSL_MODE" default:"disable"`
	PoolMaxConns        int           `envconfig:"STORE_POOL_MAX_CONNS" default:"5"`
	PoolMaxConnLifetime time.Duration `envconfig:"STORE_POOL_MAX_CONN_LIFETIME" default:"180s"`
	PoolMaxConnIdleTime time.Duration `envconfig:"STORE_POOL_MAX_CONN_IDLE_TIME" default:"100s"`
}
type Token struct {
	PrivateKeyPath  string        `envconfig:"TOKEN_PRIVATE_KEY_PATH" required:"true"`
	AccessLifetime  time.Duration `envconfig:"TOKEN_ACCESS_LIFETIME" required:"true"`
	RefreshLifetime time.Duration `envconfig:"TOKEN_REFRESH_LIFETIME" required:"true"`
}

func MustNew() *Config {
	//TODO
	if err := godotenv.Load("./../../.env"); err != nil {
		log.Fatalf("failed to load configuration: %v\n", err)
	}

	cfg := new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		log.Fatalf("failed to load configuration: %v\n", err)
	}
	return cfg
}
