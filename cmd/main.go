package main

/*TODO:
подробнее обработать 500ы ошибки
перенести некоторые проверки хэндлеров в мидлвейр
контекст!!!!
тесты
норм названия переменных
поле отмены заказа
запросы с учетом sql иньекций
env локально для паролей
*/

import (
	"Order/internal/config"
	handlers "Order/internal/http-server/handlers/order"
	"Order/internal/http-server/middleware"
	"Order/internal/http-server/middleware/logger"

	"Order/internal/storage"
	"Order/internal/storage/postgresql"
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg) //это надо хранить в секретах

	log := setUpLogger(cfg.Env)

	log.Info("starting order servise", slog.String("env", cfg.Env))

	log.Debug("debug messages are enabled")

	dbStorage, err := postgresql.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slog.Any("error", err))
		os.Exit(1)
	}

	urlService := storage.OrderService(dbStorage)

	router := gin.Default()

	router.Use(middleware.RequestID())
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer(log))

	registerHandlers(router, log, urlService)

	addr := cfg.HttpServer.Address
	if addr == "" {
		addr = ":8080"
	}

	log.Info("starting server", slog.String("address", addr))
	if err := router.Run(addr); err != nil {
		log.Error("failed to start server", slog.Any("error", err))
		os.Exit(1)
	}

}

func registerHandlers(router *gin.Engine, log *slog.Logger, service storage.OrderService) {
	router.POST("url/add", handlers.NewAdd(log, service))

	router.GET("url/getById", handlers.NewGetById(log, service))

	router.GET("url/getAll", handlers.NewGetAll(log, service))

	router.PUT("url/update/:oldId/:newId", handlers.NewUpdate(log, service))

	router.DELETE("url/delete", handlers.NewDelete(log, service))

	router.POST("url/isOrderCreated", handlers.NewIsOrderCreated(log, service))
}

func setUpLogger(env string) *slog.Logger { //конфигурация логгера
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
