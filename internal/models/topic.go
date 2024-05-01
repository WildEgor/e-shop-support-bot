package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

var (
	ActiveTopic = 1
	ClosedTopic = 0
)

type TopicCreator struct {
	TelegramId       int64  `json:"a_tid"`
	TelegramUsername string `json:"a_tun"`
	TelegramChatId   int64  `json:"c_tid" db:"-"`
}

type TopicSupport struct {
	UserId           string `json:"s_uid"`
	TelegramId       int64  `json:"s_tid"`
	TelegramUsername string `json:"s_tun"`
	TelegramChatId   int64  `json:"c_tid" db:"-"`
}

type CreateTopicFeedbackAttrs struct {
	*Topic
	Rating  uint8  `json:"rating"`
	Content string `json:"content"`
}

type TopicFeedback struct {
	CreateTopicFeedbackAttrs
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type TopicMessage struct {
	TelegramChatId    int64  `json:"c_tid" db:"-"`
	TelegramMessageId int    `json:"m_tid" db:"msg_id"`
	Question          string `json:"content" db:"content"`
}

type CreateTopicAttrs struct {
	Message TopicMessage `json:"message" db:"-"`
	Creator TopicCreator `json:"creator" db:"-"`
}

type Topic struct {
	Id         int64        `json:"id"`
	Message    TopicMessage `json:"message" db:"-"`
	Creator    TopicCreator `json:"creator" db:"-"`
	Support    TopicSupport `json:"support"`
	FeedbackId string       `json:"f_id"`
	Status     int          `json:"status"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

func (t *Topic) IsNeedFeedback() bool {
	return len(t.FeedbackId) == 0
}

func (t *Topic) IsClosed() bool {
	return t.Status == 0
}

type TopicTable struct {
	Id                      int64       `db:"id"`
	FeedbackId              pgtype.Text `db:"feedback_id"`
	AuthorTelegramId        int64       `db:"author_tid"`
	AuthorTelegramUsername  string      `db:"author_tun"`
	SupportTelegramId       int64       `db:"support_tid"`
	SupportTelegramUsername pgtype.Text `db:"support_tun"`
	MessageTelegramId       int         `db:"msg_tid"`
	MessageContent          string      `db:"content"`
	Status                  int         `db:"status"`
	CreatedAt               time.Time   `db:"created_at"`
	UpdatedAt               time.Time   `db:"updated_at"`
}

type TopicRoom struct {
	Id        int64 `json:"topic_id"`
	From      int64 `json:"from_tid"`
	To        int64 `json:"to_tid"`
	IsSupport bool  `json:"is_s"`
}
