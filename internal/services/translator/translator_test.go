package services_test

import (
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/configs"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/models"
	services "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/services/translator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslatorService_GetLocaleMessage(t *testing.T) {
	cfg := &configs.TranslatorConfig{
		DefaultLocale: "ru-RU",
		LocalesPath:   "../../locales/",
	}

	service, err := services.NewTranslatorService(cfg)
	assert.Nil(t, err)

	answer := service.GetMessageWithoutPrefix("", models.HelloMessageKey, map[string]interface{}{
		"Username": "TEST",
	})

	assert.NotEqual(t, 0, len(answer))
}
