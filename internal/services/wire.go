package services

import (
	services "github.com/WildEgor/e-shop-support-bot/internal/services/translator"
	"github.com/google/wire"
)

var ServicesSet = wire.NewSet(
	services.NewTranslatorService,
)
