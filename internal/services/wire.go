package services

import (
	services "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/services/translator"
	"github.com/google/wire"
)

var ServicesSet = wire.NewSet(
	services.NewTranslatorService,
)
