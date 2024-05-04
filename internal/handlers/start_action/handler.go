package start_action_handler

import (
	"context"
	logModels "github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/WildEgor/e-shop-support-bot/internal/adapters/telegram"
	"github.com/WildEgor/e-shop-support-bot/internal/models"
	"github.com/WildEgor/e-shop-support-bot/internal/repositories"
	services "github.com/WildEgor/e-shop-support-bot/internal/services/translator"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"time"
)

type StartActionHandler struct {
	tga *telegram.TelegramBotAdapter
	uor repositories.IUserStateRepository
	ts  *services.TranslatorService
}

func NewStartActionHandler(
	tga *telegram.TelegramBotAdapter,
	uor repositories.IUserStateRepository,
	ts *services.TranslatorService,
) *StartActionHandler {
	return &StartActionHandler{
		tga,
		uor,
		ts,
	}
}

// Handle /start command
func (h *StartActionHandler) Handle(ctx context.Context, msg *tgbotapi.Message) {
	if !msg.Chat.IsPrivate() {
		slog.Warn("chat not private")
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	uopts := &models.UserOptions{
		TelegramId: msg.From.ID,
		ChatId:     msg.Chat.ID,
	}
	uopts.SaveLang(msg.From.LanguageCode)

	if err := h.uor.SaveUserOptions(ctx, uopts); err != nil {
		slog.Error("cannot save user opts", logModels.LogEntryAttr(&logModels.LogEntry{
			Props: map[string]interface{}{
				"u_tid": msg.From.ID,
			},
		}))
		return
	}

	answer := h.ts.GetLocalizedMessage(uopts.Lang, models.HelloMessageKey, &models.UserInfoPayload{
		Username: msg.From.UserName,
	})

	ustate := h.uor.CheckUserState(ctx, msg.From.ID)

	switch ustate {
	case models.QueueState:
		answer = h.ts.GetLocalizedMessage(uopts.Lang, models.QueueStartMessageKey)
	case models.RoomState:
		answer = h.ts.GetLocalizedMessage(uopts.Lang, models.RoomStartMessageKey)
	default:
	}

	if err := h.tga.SendSimpleChatMessage(msg.Chat.ID, answer); err != nil {
		slog.Error("cannot send telegram answer", logModels.LogEntryAttr(&logModels.LogEntry{
			Props: map[string]interface{}{
				"chat_id": msg.Chat.ID,
			},
		}))
		return
	}
}
