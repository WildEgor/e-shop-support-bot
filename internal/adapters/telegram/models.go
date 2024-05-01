package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
)

type UpdateHandler func(ctx context.Context, update tgbotapi.Update)

type Middleware func(UpdateHandler) UpdateHandler

type messageHandlerFunc func(ctx context.Context, message *tgbotapi.Message)

type callbackQueryFunc func(ctx context.Context, data *tgbotapi.CallbackQuery)

type callbackQueryMatcher map[string]func(ctx context.Context, query *tgbotapi.CallbackQuery)

type messageHandler struct {
	rx *regexp.Regexp
	f  messageHandlerFunc
}

// TelegramHelpFormButtonData payload for buttons on support request form
type TelegramHelpFormButtonData struct {
	Action  string
	TopicID int64
	ChatTID int64
	UserTID int64
}

func (p TelegramHelpFormButtonData) ToData() string {
	return fmt.Sprintf("%s-%d-%d", p.Action, p.TopicID, p.UserTID)
}

type TelegramRatingFormButtonData struct {
	Action  string
	TopicID int64
	Value   uint8
}

func (p TelegramRatingFormButtonData) ToData() string {
	return fmt.Sprintf("%s-%d-%d", p.Action, p.TopicID, p.Value)
}

type TelegramCallbackButton struct {
	Text string
	Data string
}
