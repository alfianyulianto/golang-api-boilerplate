package usecase

import (
	"github.com/alfianyulianto/pds-service/pkg/storage"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BaseUseCase struct {
	DB       *gorm.DB
	Validate *validator.Validate
	Storage  storage.StorageProvider
	Config   *viper.Viper
	Log      *logrus.Entry
}

func NewBaseUseCase(DB *gorm.DB, validate *validator.Validate, storage storage.StorageProvider, config *viper.Viper, log *logrus.Entry) *BaseUseCase {
	return &BaseUseCase{DB: DB, Validate: validate, Storage: storage, Config: config, Log: log}
}
