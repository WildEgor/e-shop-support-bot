package telegram

import (
	"fmt"
	"log/slog"
)

type TelegramLogger struct {
}

func (t TelegramLogger) Println(v ...interface{}) {
	slog.Debug(fmt.Sprintf("%s", v))
}

func (t TelegramLogger) Printf(format string, v ...interface{}) {
	slog.Debug(fmt.Sprintf(format, v))
}
