package accept_callback_handler

import (
	"context"
	logModels "github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/WildEgor/e-shop-support-bot/internal/adapters/telegram"
	"github.com/WildEgor/e-shop-support-bot/internal/models"
	"github.com/WildEgor/e-shop-support-bot/internal/repositories"
	services "github.com/WildEgor/e-shop-support-bot/internal/services/translator"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"strings"
	"time"
)

// AcceptCallbackHandler handle "accept" button from support form
type AcceptCallbackHandler struct {
	tga *telegram.TelegramBotAdapter
	ts  *services.TranslatorService
	uor repositories.IUserStateRepository
	tr  repositories.ITopicsRepository
}

func NewAcceptCallbackHandler(
	tga *telegram.TelegramBotAdapter,
	ts *services.TranslatorService,
	uor repositories.IUserStateRepository,
	tr repositories.ITopicsRepository,
) *AcceptCallbackHandler {
	return &AcceptCallbackHandler{
		tga,
		ts,
		uor,
		tr,
	}
}

// Handle "accept" button from support form
func (h *AcceptCallbackHandler) Handle(ctx context.Context, update *tgbotapi.CallbackQuery) {
	data := telegram.ParseSupportTicketFormData(update.Data)
	if data == nil {
		slog.Error("error parse ticket form data")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existedTopic, _ := h.tr.FindById(data.TopicID)
	if existedTopic == nil {
		return
	}

	sopts := h.checkAndGetSupportOpts(ctx, update.ID, update.From.ID)
	uopts := h.checkAndGetUserOpts(ctx, existedTopic, update.ID)

	existedTopic.Creator.TelegramChatId = uopts.ChatId
	existedTopic.Support = models.TopicSupport{
		TelegramId:       update.From.ID,
		TelegramUsername: update.From.UserName,
		TelegramChatId:   sopts.ChatId,
	}

	existedTopic, err := h.tr.AssignTopicSupport(existedTopic)
	if err != nil {
		slog.Error("error assign support to topic", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
			Props: map[string]interface{}{
				"user_tid": update.From.ID,
			},
		}))
		return
	}

	err = h.uor.DeleteUserFromQueue(ctx, existedTopic.Creator.TelegramId)
	if err != nil {
		slog.Error("error delete user from queue", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
			Props: map[string]interface{}{
				"user_tid": existedTopic.Creator.TelegramId,
			},
		}))
		return
	}

	if err := h.tr.SaveTopicRoom(ctx, existedTopic); err != nil {
		slog.Error("error save topic room", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
			Props: map[string]interface{}{
				"user_tid": update.From.ID,
			},
		}))
		return
	}

	if err := h.tga.SendSimpleChatMessage(
		uopts.ChatId,
		h.ts.GetLocalizedMessage(uopts.Lang, models.TicketAcceptedMessageKey, &models.TicketDecisionPayload{
			TicketId:     data.TopicID,
			FromUsername: update.From.FirstName,
		}),
	); err != nil {
		slog.Error("error notify user", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
			Props: map[string]interface{}{
				"user_tid": update.From.ID,
			},
		}))
	}

	// HINT: update support request task with new text (remove accept/decline buttons)
	if err := h.tga.SendEditChatMessage(
		update.Message.Chat.ID,
		update.Message.MessageID,
		h.ts.GetMessageWithoutPrefix("", models.SuccessAcceptTicketMessageKey, &models.TicketDecisionPayload{
			TicketId:     data.TopicID,
			FromUsername: update.From.UserName,
		}),
	); err != nil {
		slog.Error("error update message", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
		}))
	}

	// TODO: check that bot not blocked from user

	h.processBufferedMessages(ctx, sopts, existedTopic, uopts.TelegramId)
}

// checkAndGetSupportOpts checks if support occupied by another user
func (h *AcceptCallbackHandler) checkAndGetSupportOpts(ctx context.Context, callbackId string, supportId int64) *models.UserOptions {
	sstate := h.uor.CheckUserState(ctx, supportId)

	sopts, err := h.uor.GetUserOptions(ctx, supportId)
	if err != nil {
		slog.Error("error get user sopts", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
			Props: map[string]interface{}{
				"user_tid": supportId,
			},
		}))
		return nil
	}

	// HINT: if support already chatting with user
	if sstate == models.RoomState {
		if err := h.tga.SendCallback(
			callbackId,
			h.ts.GetMessageWithoutPrefix(sopts.Lang, models.CompletePrevTicketMessageKey),
		); err != nil {
			slog.Error("error send accept reject callback", logModels.LogEntryAttr(&logModels.LogEntry{
				Err: err,
				Props: map[string]interface{}{
					"user_tid": supportId,
				},
			}))
		}
		return nil
	}

	return sopts
}

// checkAndGetUserOpts checks that user's topic still need help
func (h *AcceptCallbackHandler) checkAndGetUserOpts(ctx context.Context, topic *models.Topic, callbackId string) *models.UserOptions {
	ustate := h.uor.CheckUserState(ctx, topic.Creator.TelegramId)

	uopts, err := h.uor.GetUserOptions(ctx, topic.Creator.TelegramId)
	if err != nil {
		slog.Error("error get user uopts", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
			Props: map[string]interface{}{
				"user_tid": topic.Creator.TelegramId,
			},
		}))
		return nil
	}

	// HINT: when user already in another room or topic closed
	if ustate != models.QueueState || topic.IsClosed() {
		reply := h.ts.GetMessageWithoutPrefix("", models.UserNoExpectHelpMessageKey)
		if err := h.tga.SendCallback(callbackId, reply); err != nil {
			slog.Error("error send no expect callback", logModels.LogEntryAttr(&logModels.LogEntry{
				Err: err,
				Props: map[string]interface{}{
					"user_tid": topic.Creator.TelegramId,
				},
			}))
		}
		return nil
	}

	return uopts
}

func (h *AcceptCallbackHandler) processBufferedMessages(ctx context.Context, supportOpts *models.UserOptions, topic *models.Topic, userId int64) {
	messages := h.uor.GetUserMessagesFromBuffer(ctx, userId)

	var sb strings.Builder
	for _, message := range messages {
		sb.WriteString(message.Content)
		sb.WriteString("\n")
	}

	if err := h.tga.SendSimpleChatMessage(
		supportOpts.ChatId,
		h.ts.GetLocalizedMessage(supportOpts.Lang, models.NewTicketMessageKey, &models.TicketPayload{
			TicketId:      topic.Id,
			FromFirstName: topic.Creator.TelegramUsername,
			Text:          sb.String(),
		}),
	); err != nil {
		slog.Error("error send user messages", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
		}))
	}

	if err := h.uor.CleanUserMessagesFromBuffer(ctx, userId); err != nil {
		slog.Error("error clean user messages", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
		}))
	}
}
