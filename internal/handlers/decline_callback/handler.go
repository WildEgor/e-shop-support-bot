package decline_callback_handler

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

type DeclineCallbackHandler struct {
	tga *telegram.TelegramBotAdapter
	ts  *services.TranslatorService
	tr  repositories.ITopicsRepository
	uor repositories.IUserStateRepository
}

func NewDeclineCallbackHandler(
	tga *telegram.TelegramBotAdapter,
	ts *services.TranslatorService,
	tr repositories.ITopicsRepository,
	uor repositories.IUserStateRepository,
) *DeclineCallbackHandler {
	return &DeclineCallbackHandler{
		tga,
		ts,
		tr,
		uor,
	}
}

// Handle "decline" button from support form - close topic, delete user queue and messages, update callback message
func (h *DeclineCallbackHandler) Handle(ctx context.Context, update *tgbotapi.CallbackQuery) {
	data := telegram.ParseSupportTicketFormData(update.Data)
	if data == nil {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	topic, err := h.tr.FindById(data.TopicID)
	if err != nil {
		slog.Error("error find topic", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
		}))
		return
	}

	if topic != nil {
		if err := h.tr.CloseTopic(topic.Id); err != nil {
			slog.Error("error close topic", logModels.LogEntryAttr(&logModels.LogEntry{
				Err: err,
			}))
		}

		if err := h.uor.DeleteUserFromQueue(ctx, data.UserTID); err != nil {
			slog.Error("error delete user from queue", logModels.LogEntryAttr(&logModels.LogEntry{
				Err: err,
			}))
		}

		if err := h.uor.CleanUserMessagesFromBuffer(ctx, data.UserTID); err != nil {
			slog.Error("error delete user messages", logModels.LogEntryAttr(&logModels.LogEntry{
				Err: err,
			}))
		}
	}

	if err := h.tga.SendEditChatMessage(
		update.Message.Chat.ID,
		update.Message.MessageID,
		h.ts.GetMessageWithoutPrefix("", models.SuccessDeclineTicketMessageKey, &models.TicketDecisionPayload{
			TicketId:     data.TopicID,
			FromUsername: update.From.UserName,
		}),
	); err != nil {
		slog.Error("error update message", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
		}))
	}

	uopts, err := h.uor.GetUserOptions(ctx, data.UserTID)
	if err != nil {
		slog.Error("error get user callback", "data", &logModels.LogEntry{
			Err: err,
		})
		return
	}

	if err := h.tga.SendSimpleChatMessage(
		uopts.ChatId,
		h.ts.GetMessageWithoutPrefix("", models.TicketDeclinedMessageKey, &models.TicketDecisionPayload{
			TicketId: data.TopicID,
		}),
	); err != nil {
		slog.Error("error send decline reply", "data", &logModels.LogEntry{
			Err: err,
		})
		return
	}
}
