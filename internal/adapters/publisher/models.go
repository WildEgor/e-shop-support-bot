package publisher

import "context"

type PublisherType string

const (
	PublisherTypeRabbitMQ PublisherType = "rabbitmq"
)

type PublisherConfig struct {
	Type  PublisherType
	Topic string
	Addr  string
}

type IPublisherConfigFactory interface {
	Config() PublisherConfig
}

type Event struct {
	Data any
}

type IEventPublisher interface {
	Publish(context.Context, string, *Event) error
	Close() error
}

type TopicFeedbackEvent struct {
	Pattern string                 `json:"pattern"`
	Data    TopicFeedbackEventData `json:"data"`
}

type TopicFeedbackEventData struct {
	ID                      string `json:"id"`
	TopicID                 int64  `json:"t_id"`
	SupportTelegramID       int64  `json:"s_tid"`
	SupportTelegramUsername string `json:"s_tun"`
	SupportUserID           string `json:"s_uid"`
	CreatorTelegramID       int64  `json:"u_tid"`
	CreatorTelegramUsername string `json:"u_tun"`
	Rating                  uint8  `json:"rating"`
}
