package configs

import (
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/caarlos0/env/v7"
	"log/slog"
)

type RedisConfig struct {
	URI string `env:"REDIS_URI,required"`
}

func NewRedisConfig(c *Configurator) *RedisConfig {
	cfg := RedisConfig{}

	if err := env.Parse(&cfg); err != nil {
		slog.Error("redis config error", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
		panic(err)
	}

	return &cfg
}
