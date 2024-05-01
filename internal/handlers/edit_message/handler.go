package edit_message_handler

import (
	"context"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/adapters/telegram"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/mappers"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/models"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/repositories"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type EditMessageHandler struct {
	tga *telegram.TelegramBotAdapter
	uor repositories.IUserStateRepository
}

func NewEditMessageHandler(
	tga *telegram.TelegramBotAdapter,
	uor repositories.IUserStateRepository,
) *EditMessageHandler {
	return &EditMessageHandler{
		tga,
		uor,
	}
}

// Handle edit message action and update message in user buffered messages
func (h *EditMessageHandler) Handle(ctx context.Context, msg *tgbotapi.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ustate := h.uor.CheckUserState(ctx, msg.From.ID)

	switch ustate {
	case models.QueueState:
		var updated *models.UserMessage

		for _, message := range h.uor.GetUserMessagesFromBuffer(ctx, msg.From.ID) {
			if message.TelegramMessageId == msg.MessageID {
				updated = mappers.FromTelegramMessageToUserMessage(msg)
				break
			}
		}

		if updated != nil {
			// TODO: if edited (check in another handler) then find and update in buffer
			// TODO: update message in cache
		}

	case models.RoomState:
		// TODO: get room id (whomSend) and send update message version
	default:
		return
	}
}
