package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewRedis(config *viper.Viper, logrus *logrus.Entry) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.GetString("redis.host"), config.GetInt("redis.port")),
		DB:   config.GetInt("redis.db"),
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		logrus.WithField("action", "connect redis").WithError(err).Fatal("Failed to connect to redis")
	}

	return client
}
