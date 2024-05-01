package middlewares

import (
	"context"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/adapters/telegram"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AuthMiddleware check if user admin/support for callback actions
type AuthMiddleware struct {
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Next(h telegram.UpdateHandler) telegram.UpdateHandler {
	return func(ctx context.Context, u tgbotapi.Update) {
		defer h(ctx, u)

		// FIXME: dirty solution route user to no_right handler
		// If callback called with non-support privileges then redirect to no_right handler
		if u.CallbackQuery != nil && (u.CallbackQuery.Data == models.AcceptHelp || u.CallbackQuery.Data == models.DeclineHelp) {
			if !m.hasRight(u.CallbackQuery.From.ID) {
				u.CallbackQuery.Data = "no_right"
			}
		}

		// /break allow only for support redirect to no_right
		if u.Message != nil && u.Message.Text == models.BreakTopicCommand {
			if !m.hasRight(u.Message.From.ID) {
				u.CallbackQuery = &tgbotapi.CallbackQuery{
					From:    u.Message.From,
					Message: u.Message,
					Data:    models.NoRight,
				}

				u.Message = nil
			}
		}
	}
}

func (m *AuthMiddleware) hasRight(userID int64) bool {
	// TODO: check userID is admin or support
	// make request to gAuth via gRPC to fetch user by telegram id

	return true
}
