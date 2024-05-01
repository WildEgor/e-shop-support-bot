package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

func ParseSupportTicketFormData(data string) *TelegramHelpFormButtonData {
	values := strings.Split(data, "-")
	if len(values) != 3 {
		return nil
	}

	topicId, _ := strconv.ParseInt(values[1], 10, 64)
	userId, _ := strconv.ParseInt(values[2], 10, 64)

	return &TelegramHelpFormButtonData{
		Action:  values[0],
		TopicID: topicId,
		UserTID: userId,
	}
}

func ParseUserRatingFormData(data string) *TelegramRatingFormButtonData {
	values := strings.Split(data, "-")
	if len(values) != 3 {
		return nil
	}

	topicId, _ := strconv.ParseInt(values[1], 10, 64)
	value, _ := strconv.ParseInt(values[2], 10, 64)

	return &TelegramRatingFormButtonData{
		Action:  values[0],
		TopicID: topicId,
		Value:   uint8(value),
	}
}

func ToCallbackKeyboard(buttons ...[]TelegramCallbackButton) [][]tgbotapi.InlineKeyboardButton {
	keyboardRows := make([][]tgbotapi.InlineKeyboardButton, 0)

	for _, button := range buttons {
		row := make([]tgbotapi.InlineKeyboardButton, 0)

		for _, b := range button {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(b.Text, b.Data))
		}

		keyboardRows = append(keyboardRows, row)
	}

	return keyboardRows
}
