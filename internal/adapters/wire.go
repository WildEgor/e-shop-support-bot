package adapters

import (
	"github.com/WildEgor/e-shop-support-bot/internal/adapters/publisher"
	"github.com/WildEgor/e-shop-support-bot/internal/adapters/telegram"
	"github.com/google/wire"
)

var AdaptersSet = wire.NewSet(
	telegram.NewTelegramBotAdapter,
	telegram.NewTelegramListener,
	publisher.NewRabbitPublisher,
	publisher.NewEventPublisherAdapter,
)
