package publisher

import (
	"fmt"
	"github.com/wagslane/go-rabbitmq"
	"log/slog"
)

var _ rabbitmq.Logger = (*PublisherLogger)(nil)

type PublisherLogger struct {
}

func (t PublisherLogger) Fatalf(s string, i ...interface{}) {
	slog.Error(fmt.Sprintf(s, i))
	panic(fmt.Sprintf(s, i))
}

func (t PublisherLogger) Errorf(s string, i ...interface{}) {
	slog.Error(fmt.Sprintf(s, i))
}

func (t PublisherLogger) Warnf(s string, i ...interface{}) {
	slog.Warn(fmt.Sprintf(s, i))
}

func (t PublisherLogger) Infof(s string, i ...interface{}) {
	slog.Info(fmt.Sprintf(s, i))
}

func (t PublisherLogger) Debugf(s string, i ...interface{}) {
	slog.Debug(fmt.Sprintf(s, i))
}
