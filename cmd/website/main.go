package main

import (
	"fmt"
	"github.com/spf13/viper"

	"github.com/alfianyulianto/pds-service/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validator := config.NewValidator(db)
	app := config.NewFiber(viperConfig)
	redis := config.NewRedis(viperConfig, log)

	config.Boostrap(&config.BootstrapConfig{
		DB:        db,
		App:       app,
		Validator: validator,
		Config:    viperConfig,
		Log:       log,
		Redis:     redis,
	})

	baseUrl := viper.GetString("app.base_url")
	webPort := viperConfig.GetInt("app.port")
	err := app.Listen(fmt.Sprintf("%s:%d", baseUrl, webPort))
	if err != nil {
		log.WithError(err).Fatal("Error starting server:")
	}
}
