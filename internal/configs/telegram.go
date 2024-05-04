package configs

import (
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/caarlos0/env/v7"
	"log/slog"
)

// TelegramConfig holds the main app configurations
type TelegramConfig struct {
	Token  string `env:"TELEGRAM_BOT_TOKEN,required"`
	Prefix string `env:"TELEGRAM_BOT_PREFIX,required"`
	Debug  bool   `env:"TELEGRAM_DEBUG" envDefault:"false"`
}

func NewTelegramConfig(c *Configurator) *TelegramConfig {
	cfg := TelegramConfig{}

	if err := env.Parse(&cfg); err != nil {
		slog.Error("telegram config error", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
		panic(err)
	}

	return &cfg
}
