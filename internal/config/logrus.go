package config

import (
	"github.com/alfianyulianto/pds-service/internal/hooks"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
)

var (
	Logger      *logrus.Entry
	InternalLog *logrus.Logger
)

func NewLogger(config *viper.Viper) *logrus.Entry {
	// main logger setup (hooks, file output, dll)
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.Level(config.GetInt("log.level")))

	file, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Warn("Failed to open application.log, fallback to stdout:", err)
		file = os.Stdout
	}
	iw := io.MultiWriter(file, os.Stdout)
	logger.SetOutput(iw)

	// internal logger setup (no hooks, only file output)
	il := logrus.New()
	il.SetFormatter(&logrus.JSONFormatter{})
	il.SetLevel(logrus.ErrorLevel)
	il.SetOutput(file)
	// DO NOT add hooks to internal logger
	InternalLog = il

	telegramHook := hooks.NewTelegramHook(
		config.GetString("telegram.bot_token"),
		config.GetString("telegram.chat_id"),
		InternalLog,
	)

	logger.AddHook(telegramHook)

	Logger = logger.WithFields(logrus.Fields{
		"app_name": config.GetString("app.name"),
	})

	return Logger
}
