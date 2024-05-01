package mappers

import (
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func FromTelegramMessageToUserMessage(msg *tgbotapi.Message) *models.UserMessage {
	return &models.UserMessage{
		TelegramMessageId: msg.MessageID,
		TelegramUserId:    msg.From.ID,
		ChatId:            msg.Chat.ID,
		Content:           msg.Text,
	}
}
