package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"time"
)

func NewDatabase(config *viper.Viper, logrus *logrus.Entry) *gorm.DB {
	username := config.GetString("database.username")
	password := config.GetString("database.password")
	host := config.GetString("database.host")
	port := config.GetInt("database.port")
	database := config.GetString("database.name")
	maxIdleConns := config.GetInt("database.pool.max_idle_conn")
	maxOpenConns := config.GetInt("database.pool.max_open_conn")
	connMaxLifetime := config.GetDuration("database.pool.conn_max_lifetime")
	connMaxIdleTime := config.GetDuration("database.pool.conn_max_idle_time")

	file, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.WithError(err).Warn("Failed to open application.log, fallback to stdout")
		file = os.Stdout
	}
	iw := io.MultiWriter(file, os.Stdout)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(log.New(iw, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             time.Second * 5,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		}),
	})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connection database")
	}

	connection, err := db.DB()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connection database")
	}

	connection.SetMaxIdleConns(maxIdleConns)
	connection.SetMaxOpenConns(maxOpenConns)
	connection.SetConnMaxLifetime(connMaxLifetime * time.Minute)
	connection.SetConnMaxIdleTime(connMaxIdleTime * time.Minute)

	return db
}
