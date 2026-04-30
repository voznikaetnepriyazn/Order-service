package main

/*TODO:
подробнее обработать 500ы ошибки
тесты
норм названия переменных
поле отмены заказа
*/

import (
	"Order/internal/config"
	handlers "Order/internal/http-server/handlers/order"
	"Order/internal/http-server/middleware"
	"Order/internal/http-server/middleware/logger"
	"Order/internal/lib/logger/sl"

	"Order/internal/storage"
	"Order/internal/storage/postgresql"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error(".env file not found", sl.Err(err))
	}

	cfg := config.MustLoad()

	log := setUpLogger(cfg.Env)

	if log == nil {
		log = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
		slog.SetDefault(log)
	} else {
		slog.SetDefault(log)
	}

	log.Info("starting order servise", slog.String("env", cfg.Env))

	log.Debug("debug messages are enabled")

	db, err := postgresql.New(cfg.DB.DSN())
	if err != nil {
		log.Error("failed to connect to database", sl.Err(err))
	}
	defer db.Close()

	urlService := storage.OrderService(db)

	router := gin.Default()
	router.Use(middleware.RequestID())
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer(log))

	registerHandlers(router, log, urlService)

	addr := cfg.HTTPServer.Address
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
	router.POST("/url/add", handlers.NewAdd(log, service))

	router.GET("/url/getById", handlers.NewGetById(log, service))

	router.GET("/url/getAll", handlers.NewGetAll(log, service))

	router.PUT("/url/update/:oldId/:newId", handlers.NewUpdate(log, service))

	router.DELETE("/url/delete", handlers.NewDelete(log, service))

	router.POST("/url/isOrderCreated", handlers.NewIsOrderCreated(log, service))
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
