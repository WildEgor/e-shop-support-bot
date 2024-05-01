package configs

import (
	"fmt"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/caarlos0/env/v7"
	"log/slog"
)

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST,required"`
	Port     uint16 `env:"POSTGRES_PORT,required"`
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	Name     string `env:"POSTGRES_DB,required"`
	UseSSL   bool   `env:"POSTGRES_USE_SSL" envDefault:"false"`
}

func NewPostgresConfig(c *Configurator) *PostgresConfig {
	cfg := PostgresConfig{}

	if err := env.Parse(&cfg); err != nil {
		slog.Error("postgres config error", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
		panic(err)
	}

	return &cfg
}

func (p *PostgresConfig) URI() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.Name,
	)
}

func (p *PostgresConfig) MigrationURI() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.Name,
	)
}
