package configs

import (
	"github.com/WildEgor/e-shop-support-bot/internal/adapters/publisher"
	"github.com/google/wire"
)

var ConfigsSet = wire.NewSet(
	NewAppConfig,
	NewConfigurator,
	NewRedisConfig,
	NewPostgresConfig,
	NewTelegramConfig,
	NewTranslatorConfig,
	NewPublisherConfig,
	wire.Bind(new(publisher.IPublisherConfigFactory), new(*PublisherConfig)),
)
