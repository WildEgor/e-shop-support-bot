package mappers

import "github.com/WildEgor/e-shop-support-bot/internal/models"

func FromTopicsTableToModel(table *models.TopicTable) *models.Topic {
	return &models.Topic{
		Id: table.Id,
		Message: models.TopicMessage{
			TelegramMessageId: table.MessageTelegramId,
			Question:          table.MessageContent,
		},
		Creator: models.TopicCreator{
			TelegramId:       table.AuthorTelegramId,
			TelegramUsername: table.AuthorTelegramUsername,
		},
		Support: models.TopicSupport{
			TelegramId:       table.SupportTelegramId,
			TelegramUsername: table.SupportTelegramUsername.String,
		},
		FeedbackId: table.FeedbackId.String,
		Status:     table.Status,
		CreatedAt:  table.CreatedAt,
		UpdatedAt:  table.UpdatedAt,
	}
}
