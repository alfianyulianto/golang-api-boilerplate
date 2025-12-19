package config

import (
	"github.com/alfianyulianto/pds-service/internal/delivery/http/middleware"
	"github.com/alfianyulianto/pds-service/pkg/auth"
	"github.com/alfianyulianto/pds-service/pkg/email"
	storage2 "github.com/alfianyulianto/pds-service/pkg/storage"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/alfianyulianto/pds-service/internal/delivery/http"
	"github.com/alfianyulianto/pds-service/internal/delivery/http/router"
	"github.com/alfianyulianto/pds-service/internal/repository"
	"github.com/alfianyulianto/pds-service/internal/usecase"
)

type BootstrapConfig struct {
	DB        *gorm.DB
	App       *fiber.App
	Validator *validator.Validate
	Config    *viper.Viper
	Log       *logrus.Entry
	Redis     *redis.Client
}

func Boostrap(config *BootstrapConfig) {
	// storage
	var storageProvider storage2.StorageProvider
	if config.Config.GetString("storage.driver") == "local" {
		storageProvider = storage2.NewLocalStorage()
	}

	// token
	jwtConfig := &auth.JWTConfig{
		AppName:               config.Config.GetString("app.name"),
		ExpireDuration:        config.Config.GetInt("jwt.expire_duration"),
		RefreshExpireDuration: config.Config.GetInt("jwt.refresh_expire_duration"),
		SecreteKey:            config.Config.GetString("jwt.secret_key"),
	}
	jwtService := auth.NewJWTService(jwtConfig, config.Redis)

	// smtp
	smtpConfig := email.SMTPConfig{
		Host:     config.Config.GetString("mail.host"),
		Port:     config.Config.GetInt("mail.port"),
		Username: config.Config.GetString("mail.username"),
		Password: config.Config.GetString("mail.password"),
		From:     config.Config.GetString("mail.from_address"),
	}
	emailService := email.NewEmailService(&smtpConfig)

	// repositories
	userRepository := repository.NewUserRepository(config.Log)

	// useCases (service)
	baseUseCase := usecase.NewBaseUseCase(config.DB, config.Validator, storageProvider, config.Config, config.Log)
	authUseCase := usecase.NewAuthUseCase(baseUseCase, userRepository, jwtService, emailService, config.Redis)
	accountUseCase := usecase.NewAccountUseCase(baseUseCase, userRepository, emailService)
	userUseCase := usecase.NewUserUseCase(baseUseCase, userRepository)

	// controller
	authController := http.NewAuthController(authUseCase, config.Log)
	accountController := http.NewAccountController(accountUseCase, config.Log, jwtService)
	userController := http.NewUserController(userUseCase, config.Log)

	// middleware
	httpMiddleware := middleware.NewMiddleware(config.Log, jwtService)

	routerConfig := router.RouterConfig{
		App:               config.App,
		Middleware:        httpMiddleware,
		AuthController:    authController,
		AccountController: accountController,
		UserController:    userController,
	}

	routerConfig.Setup()
}
