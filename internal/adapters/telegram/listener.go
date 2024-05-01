package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"regexp"
	"strings"
)

// TelegramListener handle updates
type TelegramListener struct {
	adapter *TelegramBotAdapter
	// middlewares run before each handler
	middlewares []Middleware
	// messageHandlers handle messages
	messageHandlers []messageHandler
	// editMessageHandler handle edit message
	editMessageHandler messageHandlerFunc
	// callbackQueryMatcher handle callbacks
	callbackQueryMatcher callbackQueryMatcher
}

func NewTelegramListener(
	adapter *TelegramBotAdapter,
) *TelegramListener {
	return &TelegramListener{
		adapter:              adapter,
		messageHandlers:      make([]messageHandler, 0),
		editMessageHandler:   func(ctx context.Context, message *tgbotapi.Message) {},
		callbackQueryMatcher: make(callbackQueryMatcher),
	}
}

// Use register middleware
func (t *TelegramListener) Use(middleware Middleware) {
	t.middlewares = append(t.middlewares, middleware)
}

// RegisterCallbackHandler Define callback handlers per key, and the key is actually the cq.Data we attach to our buttons
// Note: It only works if you call HandleCallback along this function.
func (t *TelegramListener) RegisterCallbackHandler(key string, handler callbackQueryFunc) {
	t.callbackQueryMatcher[key] = handler
}

// HandleMessage sets handler for incoming messages
func (t *TelegramListener) HandleMessage(pattern string, handler messageHandlerFunc) {
	rx := regexp.MustCompile(pattern)
	t.messageHandlers = append(t.messageHandlers, messageHandler{rx: rx, f: handler})
}

// HandleEditedMessage set handler for incoming edited messages
func (t *TelegramListener) HandleEditedMessage(handler messageHandlerFunc) {
	t.editMessageHandler = handler
}

// Stop listen
func (t *TelegramListener) Stop() {
	t.adapter.Stop()
}

// ListenUpdates handle updates in sep goroutine
func (t *TelegramListener) ListenUpdates(ctx context.Context) {
	slog.Debug("bot is listening")

	h := func(ctx context.Context, update tgbotapi.Update) {
		slog.Debug("Update ID", slog.Int("ID", update.UpdateID))

		var f = t.handleUpdates
		for i := len(t.middlewares) - 1; i >= 0; i-- {
			f = t.middlewares[i](f)
		}

		go f(ctx, update)
	}

	t.adapter.HandleUpdates(ctx, h)
}

// handleUpdates catch and route update to handlers
func (t *TelegramListener) handleUpdates(ctx context.Context, update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		t.handleMessage(ctx, update.Message)
	case update.EditedMessage != nil:
		t.editMessageHandler(ctx, update.EditedMessage)
	case update.CallbackQuery != nil:
		t.handleCallback(ctx, update.CallbackQuery)
	}
}

// handleCallback catch callback updates
func (t *TelegramListener) handleCallback(ctx context.Context, update *tgbotapi.CallbackQuery) {
	values := strings.Split(update.Data, "-")
	if len(values) == 0 {
		return
	}
	handler, ok := t.callbackQueryMatcher[values[0]]
	if !ok {
		return
	}
	handler(ctx, update)
}

// handleMessage catch messages
func (t *TelegramListener) handleMessage(ctx context.Context, message *tgbotapi.Message) {
	for _, handler := range t.messageHandlers {
		if handler.rx.MatchString(message.Text) {
			handler.f(ctx, message)
			return
		}
	}
}
