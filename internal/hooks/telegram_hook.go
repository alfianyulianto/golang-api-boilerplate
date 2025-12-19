package hooks

import (
	"fmt"
	"github.com/alfianyulianto/pds-service/pkg/telegram"
	"github.com/sirupsen/logrus"
)

type TelegramHook struct {
	BotToken    string
	ChatID      string
	Client      *telegram.TelegramClient
	InternalLog *logrus.Logger
}

func NewTelegramHook(botToken, chatID string, internalLog *logrus.Logger) *TelegramHook {
	return &TelegramHook{
		BotToken:    botToken,
		ChatID:      chatID,
		Client:      telegram.NewTelegramClient(botToken, chatID),
		InternalLog: internalLog,
	}
}

func (h *TelegramHook) Fire(entry *logrus.Entry) error {
	go func() {
		// Format message untuk Telegram dengan HTML
		message := fmt.Sprintf(
			"<b>🚨 %s - Log Alert</b>\n\n"+
				"<b>Level:</b> %s\n"+
				"<b>Message:</b> %s\n"+
				"<b>Time:</b> %s\n"+
				"<b>Data:</b> %v",
			entry.Data["app_name"],
			entry.Level.String(),
			entry.Message,
			entry.Time.Format("2006-01-02 15:04:05"),
			entry.Data,
		)

		err := h.Client.SendMessage(message)
		if err != nil {
			h.InternalLog.WithField("action", "telegram hook").WithError(err).Error("Telegram hook failed")
		}
	}()
	return nil
}

func (h *TelegramHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}
}
