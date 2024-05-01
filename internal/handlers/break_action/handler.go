package break_action_handler

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

type BreakActionHandler struct {
	tga *telegram.TelegramBotAdapter
	ts  *services.TranslatorService
	tr  repositories.ITopicsRepository
	uor repositories.IUserStateRepository
}

func NewBreakActionHandler(
	tga *telegram.TelegramBotAdapter,
	ts *services.TranslatorService,
	tr repositories.ITopicsRepository,
	uor repositories.IUserStateRepository,
) *BreakActionHandler {
	return &BreakActionHandler{
		tga,
		ts,
		tr,
		uor,
	}
}

// Handle /break command - close topic, delete room and queue user/messages
func (h *BreakActionHandler) Handle(ctx context.Context, msg *tgbotapi.Message) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	room, err := h.tr.GetTopicRoomBySenderId(ctx, msg.From.ID)
	if room == nil {
		slog.Warn("error room not found", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
		}))
		return
	}

	topic, err := h.tr.FindById(room.Id)
	if err != nil {
		slog.Error("error find topic", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
		}))
		return
	}

	if err := h.tr.CloseTopic(topic.Id); err != nil {
		slog.Error("error close topic", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
		}))
	}

	if err := h.tr.RemoveTopicRoom(ctx, topic); err != nil {
		slog.Error("error remove topic room", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
		}))
	}

	h.notifyAboutFeedback(ctx, topic, msg)
	h.notifyTopicClosed(ctx, msg)
}

func (h *BreakActionHandler) notifyAboutFeedback(ctx context.Context, topic *models.Topic, msg *tgbotapi.Message) {
	uopts, err := h.uor.GetUserOptions(ctx, topic.Creator.TelegramId)
	if err != nil {
		slog.Error("cannot get user uopts", logModels.LogEntryAttr(&logModels.LogEntry{
			Props: map[string]interface{}{
				"m_tid": msg.From.ID,
			},
		}))
		return
	}

	lfdbMsg := h.ts.GetLocalizedMessage(uopts.Lang, models.LeaveFeedbackMessageKey, &models.FeedbackTemplatePayload{
		TicketId: topic.Id,
	})

	keyboard := telegram.ToCallbackKeyboard(
		[]telegram.TelegramCallbackButton{
			{
				Text: h.ts.GetMessageWithoutPrefix(uopts.Lang, models.RatingOk),
				Data: telegram.TelegramRatingFormButtonData{
					Action:  models.Rating,
					TopicID: topic.Id,
					Value:   5,
				}.ToData(),
			},
			{
				Text: h.ts.GetMessageWithoutPrefix(uopts.Lang, models.RatingNormal),
				Data: telegram.TelegramRatingFormButtonData{
					Action:  models.Rating,
					TopicID: topic.Id,
					Value:   3,
				}.ToData(),
			},
			{
				Text: h.ts.GetMessageWithoutPrefix(uopts.Lang, models.RatingBad),
				Data: telegram.TelegramRatingFormButtonData{
					Action:  models.Rating,
					TopicID: topic.Id,
					Value:   1,
				}.ToData(),
			},
		},
		[]telegram.TelegramCallbackButton{
			{
				Text: h.ts.GetMessageWithoutPrefix(uopts.Lang, models.NoRating),
				Data: telegram.TelegramRatingFormButtonData{
					Action:  models.Rating,
					TopicID: topic.Id,
					Value:   0,
				}.ToData(),
			},
		},
	)

	if err := h.tga.SendForm(uopts.ChatId, lfdbMsg, keyboard); err != nil {
		slog.Error("error send form", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
		}))
	}
}

func (h *BreakActionHandler) notifyTopicClosed(ctx context.Context, msg *tgbotapi.Message) {
	topts, err := h.uor.GetUserOptions(ctx, msg.From.ID)
	if err != nil {
		slog.Error("cannot get user topts", "data", &logModels.LogEntry{
			Err: err,
			Props: map[string]interface{}{
				"user_tid": msg.From.ID,
			},
		})
	}

	tendMsg := h.ts.GetLocalizedMessage(topts.Lang, models.TicketClosedMessageKey)
	if err := h.tga.SendSimpleChatMessage(msg.Chat.ID, tendMsg); err != nil {
		slog.Error("error send end ticket", "data", &logModels.LogEntry{
			Err: err,
			Props: map[string]interface{}{
				"user_tid": msg.From.ID,
			},
		})
	}
}
