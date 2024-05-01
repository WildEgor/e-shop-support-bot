package adapters

import (
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/adapters/telegram"
	"github.com/google/wire"
)

var AdaptersSet = wire.NewSet(
	telegram.NewTelegramBotAdapter,
	telegram.NewTelegramListener,
)
