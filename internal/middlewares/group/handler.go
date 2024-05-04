package middlewares

import (
	"context"
	logModels "github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/WildEgor/e-shop-support-bot/internal/adapters/telegram"
	"github.com/WildEgor/e-shop-support-bot/internal/configs"
	"github.com/WildEgor/e-shop-support-bot/internal/repositories"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"strings"
)

// ExtractGroupMiddleware
type ExtractGroupMiddleware struct {
	uor  repositories.IUserStateRepository
	gr   repositories.IGroupRepository
	tcfg *configs.TelegramConfig
}

func NewExtractGroupMiddleware(
	uor repositories.IUserStateRepository,
	gr repositories.IGroupRepository,
	tcfg *configs.TelegramConfig,
) *ExtractGroupMiddleware {
	return &ExtractGroupMiddleware{
		uor,
		gr,
		tcfg,
	}
}

func (m *ExtractGroupMiddleware) Next(h telegram.UpdateHandler) telegram.UpdateHandler {
	return func(ctx context.Context, u tgbotapi.Update) {
		defer h(ctx, u)

		slog.Debug("", slog.Any("update", u.MyChatMember))

		if u.MyChatMember == nil || u.MyChatMember.NewChatMember.Status != "administrator" {
			slog.Debug("ignore")
			return
		}

		slog.Debug("Try save group")

		chatID := u.MyChatMember.Chat.ID

		if strings.Contains(u.MyChatMember.NewChatMember.User.UserName, m.tcfg.Prefix) {
			existed, err := m.gr.GetGroupId(ctx)
			if err != nil {
				slog.Error("err get group id", logModels.LogEntryAttr(&logModels.LogEntry{
					Err: err,
					Props: map[string]interface{}{
						"chat_id": chatID,
					},
				}))
			}

			if existed == chatID {
				return
			}

			if err := m.gr.SaveGroupId(ctx, chatID); err != nil {
				slog.Error("cannot save group id", "data", &logModels.LogEntry{
					Err: err,
					Props: map[string]interface{}{
						"chat_id": chatID,
					},
				})
			}
		}
	}
}
