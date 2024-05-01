package configs

import (
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/caarlos0/env/v7"
	"log/slog"
	"path"
)

type TranslatorConfig struct {
	DefaultLocale string `env:"DEFAULT_LOCALE" envDefault:"ru-RU"`
	Prefix        string `env:"PREFIX" envDefault:""`
	LocalesPath   string
}

func NewTranslatorConfig(c *Configurator) *TranslatorConfig {
	cfg := TranslatorConfig{
		LocalesPath: "internal/locales/",
	}

	if err := env.Parse(&cfg); err != nil {
		slog.Error("translator config error", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
		panic(err)
	}

	return &cfg
}

func (c *TranslatorConfig) LocalesFullPath() string {
	return c.LocalesPath + "/*/*"
}

func (c *TranslatorConfig) LocalesDir() string {
	return path.Join(c.LocalesPath)
}
