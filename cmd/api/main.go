// @title WireGuard Manager API
// @version 1.0
// @description API для управления конфигурациями WireGuard.
// @termsOfService http://example.com/terms/
// @contact.name Support Team
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost
// @BasePath /api/v1

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"time"
	"wg_api/config"
	_ "wg_api/docs"
	"wg_api/internal/controllers"
	"wg_api/internal/models"
	"wg_api/internal/repository/postgres"
	"wg_api/internal/scheduler"
	"wg_api/internal/services"
	"wg_api/pkg/database"
	"wg_api/pkg/shell"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Подключение к базе данных
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Автомиграция моделей
	err = db.AutoMigrate(&models.User{}, &models.Server{}, &models.Configuration{})
	if err != nil {
		log.Fatalf("Failed to migrate database models: %v", err)
	}

	// Инициализация репозиториев
	userRepo := postgres.NewUserRepository(db)
	serverRepo := postgres.NewServerRepository(db)
	configRepo := postgres.NewConfigurationRepository(db)

	// Инициализация сервисов
	userService := services.NewUserService(userRepo)
	serverService := services.NewServerService(serverRepo)
	configService := services.NewConfigurationService(configRepo)

	// Инициализация исполнителя shell-команд
	shellExecutor := shell.NewExecutor()

	// Инициализация сервиса WireGuard
	wireguardService := services.NewWireGuardService(shellExecutor, cfg, configService)

	// Инициализация контроллеров
	userController := controllers.NewUserController(userService)
	serverController := controllers.NewServerController(serverService)
	configController := controllers.NewConfigurationController(configService, wireguardService)

	// Инициализация планировщика задач
	schedulerInterval := time.Hour // Раз в час
	taskScheduler := scheduler.NewScheduler(configService, wireguardService, schedulerInterval)
	taskScheduler.Start()
	defer taskScheduler.Stop()

	// Инициализация маршрутизатора Gin
	router := gin.Default()

	// Генерация документации swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Группа маршрутов API
	api := router.Group("/api/v1")

	// Регистрация маршрутов
	userController.RegisterRoutes(api)
	serverController.RegisterRoutes(api)
	configController.RegisterRoutes(api)

	// Запуск сервера
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
