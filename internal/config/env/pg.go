package config

import (
	"errors"
	"os"

	"github.com/MariaPopova99/microservices/internal/config"
)

const (
	dsnEnvName = "PG_DSN"
)

var _ config.PGConfig = (*pgConfig)(nil)

type pgConfig struct {
	dsn string
}

func NewPGConfig() (*pgConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
