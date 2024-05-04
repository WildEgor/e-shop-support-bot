package new_message_handler

import (
	"context"
	"errors"
	logModels "github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/WildEgor/e-shop-support-bot/internal/adapters/telegram"
	"github.com/WildEgor/e-shop-support-bot/internal/mappers"
	"github.com/WildEgor/e-shop-support-bot/internal/models"
	"github.com/WildEgor/e-shop-support-bot/internal/repositories"
	services "github.com/WildEgor/e-shop-support-bot/internal/services/translator"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"time"
)

type NewMessageHandler struct {
	tga *telegram.TelegramBotAdapter
	ts  *services.TranslatorService
	tr  repositories.ITopicsRepository
	uor repositories.IUserStateRepository
	gr  repositories.IGroupRepository
}

func NewNewMessageHandler(
	tga *telegram.TelegramBotAdapter,
	ts *services.TranslatorService,
	tr repositories.ITopicsRepository,
	uor repositories.IUserStateRepository,
	gr repositories.IGroupRepository,
) *NewMessageHandler {
	return &NewMessageHandler{
		tga,
		ts,
		tr,
		uor,
		gr,
	}
}

// Handle any new message that send to bot
func (h *NewMessageHandler) Handle(ctx context.Context, msg *tgbotapi.Message) {
	// HINT: don't process any messages not from 1-to-1 chats
	if !msg.Chat.IsPrivate() {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// HINT: use language reply
	uopts, err := h.uor.GetUserOptions(ctx, msg.From.ID)
	if err != nil {
		slog.Error("cannot get user uopts", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
			Props: map[string]interface{}{
				"user_tid": msg.From.ID,
				"m_tid":    msg.MessageID,
			},
		}))
		return
	}

	ustate := h.uor.CheckUserState(ctx, msg.From.ID)

	userMsg := mappers.FromTelegramMessageToUserMessage(msg)

	// HINT: Admin or support always in default/rooms state
	// User may in any states
	switch ustate {
	case models.DefaultState:
		existedTopic, err := h.tr.FindFirstAuthorActiveTopic(msg.From.ID)
		// TODO: make better error handling
		if errors.Is(err, pgx.ErrNoRows) {
			existedTopic, err = h.tr.CreateUniqueTopic(
				&models.CreateTopicAttrs{
					Message: models.TopicMessage{
						TelegramMessageId: msg.MessageID,
						TelegramChatId:    msg.Chat.ID,
						Question:          msg.Text,
					},
					Creator: models.TopicCreator{
						TelegramId:       msg.From.ID,
						TelegramUsername: msg.From.UserName,
					},
				})
		}

		if existedTopic == nil {
			slog.Error("error create new topic", logModels.LogEntryAttr(&logModels.LogEntry{
				Err: err,
			}))
			return
		}

		if err := h.uor.AddUserToQueue(ctx, msg.From.ID); err != nil {
			slog.Error("error add user to queue", logModels.LogEntryAttr(&logModels.LogEntry{
				Err: err,
				Props: map[string]interface{}{
					"user_tid": msg.From.ID,
				},
			}))
		}

		if err := h.uor.AddUserMessageToBuffer(ctx, userMsg); err != nil {
			slog.Error("error buffer messages", logModels.LogEntryAttr(&logModels.LogEntry{
				Err: err,
				Props: map[string]interface{}{
					"user_tid": msg.From.ID,
				},
			}))
		}

		if err := h.tga.SendSimpleChatMessage(
			msg.Chat.ID,
			h.ts.GetLocalizedMessage(uopts.Lang, models.GotTicketKey, &models.TicketPayload{
				TicketId:     existedTopic.Id,
				FromUsername: msg.From.UserName,
			}),
		); err != nil {
			slog.Error("error uReplyMsg message", logModels.LogEntryAttr(&logModels.LogEntry{
				Err: err,
			}))
			return
		}

		// HINT: send notification to support group
		groupId, err := h.gr.GetGroupId(ctx)
		if err != nil {
			slog.Error("cannot notify support cause empty group id", logModels.LogEntryAttr(&logModels.LogEntry{
				Err: err,
			}))
			return
		}

		// HINT: use default lang (RU)
		newMessage := h.ts.GetLocalizedMessage("", models.NewTicketCreatedMessageKey, &models.TicketDecisionPayload{
			TicketId:      existedTopic.Id,
			FromFirstName: msg.From.FirstName,
			Text:          existedTopic.Message.Question,
		})

		keyboard := telegram.ToCallbackKeyboard([]telegram.TelegramCallbackButton{
			{
				Text: h.ts.GetMessageWithoutPrefix("", models.AcceptTicketMessageKey),
				Data: telegram.TelegramHelpFormButtonData{
					Action:  models.AcceptHelp,
					ChatTID: msg.Chat.ID,
					TopicID: existedTopic.Id,
					UserTID: existedTopic.Creator.TelegramId,
				}.ToData(),
			},
			{
				Text: h.ts.GetMessageWithoutPrefix("", models.DeclineTicketMessageKey),
				Data: telegram.TelegramHelpFormButtonData{
					Action:  models.DeclineHelp,
					ChatTID: msg.Chat.ID,
					TopicID: existedTopic.Id,
					UserTID: existedTopic.Creator.TelegramId,
				}.ToData(),
			},
		})

		if err := h.tga.SendForm(groupId, newMessage, keyboard); err != nil {
			slog.Error("error send help form", logModels.LogEntryAttr(&logModels.LogEntry{
				Err: err,
			}))
		}
	case models.QueueState:
		// TODO: it's okay if we collect all messages?
		if err := h.uor.AddUserMessageToBuffer(ctx, userMsg); err != nil {
			slog.Error("error queue buffer messages", logModels.LogEntryAttr(&logModels.LogEntry{
				Err: err,
			}))
		}
	case models.RoomState:
		room, err := h.tr.GetTopicRoomBySenderId(ctx, msg.From.ID)
		if err != nil {
			slog.Error("error get topic room", logModels.LogEntryAttr(&logModels.LogEntry{
				Err: err,
			}))
			return
		}

		// TODO: handle this
		if room == nil {
			slog.Error("room is empty", logModels.LogEntryAttr(&logModels.LogEntry{
				Err: err,
			}))
			return
		}

		sopts, err := h.uor.GetUserOptions(ctx, msg.From.ID)
		if err != nil {
			slog.Error("error get user sopts", logModels.LogEntryAttr(&logModels.LogEntry{
				Err: err,
				Props: map[string]interface{}{
					"user_tid": msg.From.ID,
				},
			}))
			return
		}
		if sopts.ChatId == 0 {
			return
		}

		// TODO: process if bot blocked
		// TODO: send reply to sender that message/comment to topic send

		if room.IsSupport {
			if err := h.tga.SendSimpleChatMessage(
				sopts.ChatId,
				h.ts.GetLocalizedMessage(uopts.Lang, models.NewTicketCommentMessageKey, &models.TicketPayload{
					TicketId:     room.Id,
					FromUsername: msg.From.UserName,
					Text:         msg.Text,
				}),
			); err != nil {
				slog.Error("error send redirects", logModels.LogEntryAttr(&logModels.LogEntry{
					Err: err,
					Props: map[string]interface{}{
						"user_tid": msg.From.ID,
					},
				}))
				return
			}
		}

		if err := h.tga.SendSimpleChatMessage(
			sopts.ChatId,
			h.ts.GetLocalizedMessage(uopts.Lang, models.NewTicketMessageKey, &models.TicketPayload{
				TicketId:      room.Id,
				FromFirstName: msg.From.FirstName,
				Text:          msg.Text,
			}),
		); err != nil {
			slog.Error("error send redirects", logModels.LogEntryAttr(&logModels.LogEntry{
				Err: err,
				Props: map[string]interface{}{
					"user_tid": msg.From.ID,
				},
			}))
			return
		}
	}
}
