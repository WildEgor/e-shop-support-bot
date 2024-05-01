package rating_callback_handler

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

type RatingCallbackHandler struct {
	tga *telegram.TelegramBotAdapter
	ts  *services.TranslatorService
	uor repositories.IUserStateRepository
	tr  repositories.ITopicsRepository
}

func NewRatingCallbackHandler(
	tga *telegram.TelegramBotAdapter,
	ts *services.TranslatorService,
	uor repositories.IUserStateRepository,
	tr repositories.ITopicsRepository,
) *RatingCallbackHandler {
	return &RatingCallbackHandler{
		tga,
		ts,
		uor,
		tr,
	}
}

// Handle send feedback
func (h *RatingCallbackHandler) Handle(ctx context.Context, update *tgbotapi.CallbackQuery) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data := telegram.ParseUserRatingFormData(update.Data)
	if data == nil {
		slog.Error("error parse rating data")
		return
	}

	topic, err := h.tr.FindById(data.TopicID)
	if err != nil {
		slog.Error("error find topic", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
		}))
		return
	}

	if data.Value != 0 && topic != nil && topic.IsNeedFeedback() {
		_, err = h.tr.LeaveFeedback(&models.CreateTopicFeedbackAttrs{
			Topic:  topic,
			Rating: data.Value,
		})
		if err != nil {
			slog.Error("error leave feedback", logModels.LogEntryAttr(&logModels.LogEntry{
				Err: err,
			}))
			return
		}
	}

	uopts, err := h.uor.GetUserOptions(ctx, update.From.ID)
	if err != nil {
		slog.Error("cannot get user uopts", logModels.LogEntryAttr(&logModels.LogEntry{
			Props: map[string]interface{}{
				"t_id": update.From.ID,
			},
		}))
		return
	}

	reply := h.ts.GetLocalizedMessage(uopts.Lang, models.FeedbackSendMessageKey, &models.TicketPayload{
		TicketId: 0,
	})
	if topic != nil {
		reply = h.ts.GetLocalizedMessage(uopts.Lang, models.FeedbackSendMessageKey, &models.TicketPayload{
			TicketId: topic.Id,
		})
	}

	if err := h.tga.SendEditChatMessage(
		update.Message.Chat.ID,
		update.Message.MessageID,
		reply,
	); err != nil {
		slog.Error("error update message", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
		}))
	}

	// TODO: we can update message in group with rating. Need save message id in cache
}
