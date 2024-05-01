package no_right_callback_handler

import (
	"context"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/adapters/telegram"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/models"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/repositories"
	services "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/services/translator"
	logModels "github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"time"
)

type NoRightCallbackHandler struct {
	ts  *services.TranslatorService
	tga *telegram.TelegramBotAdapter
	uor repositories.IUserStateRepository
}

func NewNoRightCallbackHandler(
	ts *services.TranslatorService,
	tga *telegram.TelegramBotAdapter,
	uor repositories.IUserStateRepository,
) *NoRightCallbackHandler {
	return &NoRightCallbackHandler{
		ts,
		tga,
		uor,
	}
}

// Handle send message that user no right to perform action
func (h *NoRightCallbackHandler) Handle(ctx context.Context, query *tgbotapi.CallbackQuery) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var lang string
	uopts, err := h.uor.GetUserOptions(ctx, query.From.ID)
	if uopts != nil {
		lang = uopts.Lang
	}

	err = h.tga.SendSimpleChatMessage(query.Message.Chat.ID, h.ts.GetMessageWithoutPrefix(lang, models.NoRightMessageKey))
	if err != nil {
		slog.Error("cannot no right send callback", logModels.LogEntryAttr(&logModels.LogEntry{
			Props: map[string]interface{}{
				"m_tid": query.Message.MessageID,
			},
		}))
	}

}
