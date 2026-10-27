package postgres

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	URL             string        `envconfig:"DATABASE_URL"`
	Host            string        `envconfig:"DB_HOST" default:"localhost"`
	Port            int           `envconfig:"DB_PORT" default:"5432"`
	User            string        `envconfig:"DB_USER"`
	Password        string        `envconfig:"DB_PASSWORD"`
	Database        string        `envconfig:"DB_NAME"`
	SSLMode         string        `envconfig:"DB_SSLMODE" default:"disable"`
	MaxConns        int32         `envconfig:"DB_MAX_CONNS" default:"10"`
	MinConns        int32         `envconfig:"DB_MIN_CONNS" default:"0"`
	MaxConnLifetime time.Duration `envconfig:"DB_MAX_CONN_LIFETIME" default:"30m"`
	ConnTimeout     time.Duration `envconfig:"DB_CONN_TIMEOUT" default:"5s"`
}

func LoadFromEnv() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		// Do not include the original parsing error: it can contain configuration values.
		return Config{}, fmt.Errorf("invalid PostgreSQL configuration: check environment variable types")
	}
	if cfg.URL == "" && (cfg.User == "" || cfg.Password == "" || cfg.Database == "") {
		return Config{}, fmt.Errorf("set DATABASE_URL or DB_USER, DB_PASSWORD and DB_NAME")
	}
	if cfg.Port < 1 || cfg.Port > 65535 {
		return Config{}, fmt.Errorf("DB_PORT must be between 1 and 65535")
	}
	if cfg.MaxConns < 1 || cfg.MinConns < 0 || cfg.MinConns > cfg.MaxConns {
		return Config{}, fmt.Errorf("DB_MAX_CONNS must be positive and DB_MIN_CONNS must be between 0 and DB_MAX_CONNS")
	}
	if cfg.ConnTimeout <= 0 || cfg.MaxConnLifetime <= 0 {
		return Config{}, fmt.Errorf("DB_CONN_TIMEOUT and DB_MAX_CONN_LIFETIME must be positive")
	}
	return cfg, nil
}
