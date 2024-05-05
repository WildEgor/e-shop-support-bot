package configs

import (
	"github.com/WildEgor/e-shop-support-bot/internal/adapters/publisher"
	"github.com/caarlos0/env/v7"
	"log/slog"
)

var _ publisher.IPublisherConfigFactory = (*PublisherConfig)(nil)

// PublisherConfig holds the main app configurations
type PublisherConfig struct {
	Topic string                  `env:"PUBLISHER_TOPIC,required"`
	Addr  string                  `env:"PUBLISHER_ADDR,required"`
	Type  publisher.PublisherType `env:"PUBLISHER_TYPE,required"`
}

func NewPublisherConfig(c *Configurator) *PublisherConfig {
	cfg := PublisherConfig{}

	if err := env.Parse(&cfg); err != nil {
		slog.Error("app publisher parse error")
	}

	return &cfg
}

func (p *PublisherConfig) Config() publisher.PublisherConfig {
	return publisher.PublisherConfig{
		Topic: p.Topic,
		Addr:  p.Addr,
	}
}
