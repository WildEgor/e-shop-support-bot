package telegram

import (
	"context"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/WildEgor/e-shop-support-bot/internal/configs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

// TelegramBotAdapter wrap requests to telegram api
type TelegramBotAdapter struct {
	bot     *tgbotapi.BotAPI
	handler UpdateHandler
}

func NewTelegramBotAdapter(cfg *configs.TelegramConfig) *TelegramBotAdapter {
	slog.Debug("Create new bot")

	bot, err := tgbotapi.NewBotAPI(cfg.Token)

	logger := &TelegramLogger{}
	if err := tgbotapi.SetLogger(logger); err != nil {
		return nil
	}

	if err != nil {
		slog.Error("error init bot", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
		panic(err)
	}

	bot.Debug = cfg.Debug

	return &TelegramBotAdapter{
		bot:     bot,
		handler: func(ctx context.Context, update tgbotapi.Update) {},
	}
}

func (t *TelegramBotAdapter) SendChatAction(chatId int64, action string) error {
	cmd := tgbotapi.NewChatAction(chatId, action)
	_, err := t.bot.Send(cmd)
	return err
}

func (t *TelegramBotAdapter) SendSimpleChatMessage(chatId int64, message string) error {
	cmd := tgbotapi.NewMessage(chatId, message)
	cmd.ParseMode = tgbotapi.ModeMarkdown
	_, err := t.bot.Send(cmd)
	return err
}

func (t *TelegramBotAdapter) SendEditChatMessage(chatId int64, msgId int, message string) error {
	cmd := tgbotapi.NewEditMessageText(chatId, msgId, message)
	cmd.ParseMode = tgbotapi.ModeMarkdown
	_, err := t.bot.Send(cmd)
	return err
}

func (t *TelegramBotAdapter) SendCallback(chatId string, message string) error {
	callback := tgbotapi.NewCallback(
		chatId,
		message,
	)
	_, err := t.bot.Send(callback)
	return err
}

// TODO: refactor
func (t *TelegramBotAdapter) SendForm(id int64, text string, buttons [][]tgbotapi.InlineKeyboardButton) error {
	msg := tgbotapi.NewMessage(
		id,
		text,
	)

	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}

	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err := t.bot.Send(msg)
	return err
}

func (t *TelegramBotAdapter) HandleUpdates(ctx context.Context, handler UpdateHandler) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for update := range t.bot.GetUpdatesChan(u) {
		handler(ctx, update)
	}
}

func (t *TelegramBotAdapter) Stop() {
	t.bot.StopReceivingUpdates()
}
