package models

type UserOptions struct {
	TelegramId int64  `json:"t_id"`
	ChatId     int64  `json:"chat_id"`
	Lang       string `json:"lang"`
}

func (uo *UserOptions) SaveLang(lang string) {
	if lang == "ru" {
		uo.Lang = "ru-RU"
	}

	if lang == "en" {
		uo.Lang = "en-EN"
	}
}

type UserMessage struct {
	TelegramMessageId int    `json:"m_tid"`
	TelegramUserId    int64  `json:"u_tid"`
	ChatId            int64  `json:"chat_id"`
	Content           string `json:"content"`
}
