package handlers

import (
	"errors"
	"log/slog"

	"Order/internal/http-server/middleware/logger"
	uuidparam "Order/internal/http-server/middleware/uuid"
	resp "Order/internal/lib/api/response"
	"Order/internal/lib/logger/sl"
	"Order/internal/models/order"
	"Order/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	URL string `json:"url" validate:"required, url"`
}

type RequestFullStruct struct {
	Order order.Order
}

type Request struct {
	resp.Response
	URL string `json:"url" validate:"required, url"`
}

type Crud interface {
	NewAdd(log *slog.Logger, adder storage.OrderService) gin.HandlerFunc
	NewDelete(log *slog.Logger, deleter storage.OrderService) gin.HandlerFunc
	NewGetAll(log *slog.Logger, get storage.OrderService) gin.HandlerFunc
	NewGetById(log *slog.Logger, get storage.OrderService) gin.HandlerFunc
	NewUpdate(log *slog.Logger, update storage.OrderService) gin.HandlerFunc
	NewIsOrderCreated(log *slog.Logger, ord storage.OrderService) gin.HandlerFunc
}

func NewAdd(log *slog.Logger, adder storage.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.url.add.New"

		log := logger.FromCtx(c)

		log.Info("handling request")

		var req RequestFullStruct

		//декодирование из джэйсона в структуру данных
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			c.JSON(400, gin.H{
				"error": "failed to decode request",
			})

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		//валидация - просто проверка на то, есть ли значение(в данном случае)
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			c.JSON(400, gin.H{
				"error":   "validation failed",
				"details": formatValidationError(validateErr),
			})

			return
		}

		//проверка на уже существующее значение
		id, err := adder.AddURL(req.Order)
		/*if errors.Is(err, storage.ErrUrlExist) {
			log.Info("url already exists", slog.Any("url", req.order))

			c.JSON(400, gin.H{
				"error": "url already exists",
			})

			return
		}*/

		//прочие ошибки
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			c.JSON(500, gin.H{
				"error": "failed to add url",
			})

			return
		}

		log.Info("url added", slog.Any("id", id))

		responseOK(c)
	}
}

func responseOK(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "OK",
	})
}

func formatValidationError(err validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)
	for _, e := range err {
		errors[e.Field()] = e.Error()
	}
	return errors
}

func NewDelete(log *slog.Logger, deleter storage.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.url.delete.New"

		log := logger.FromCtx(c)

		log.Info("handling request")

		uuidParam, ok := uuidparam.UUIDFromCtx(c, "id")
		if ok {
			return
		}

		err := deleter.DeleteURL(uuidParam)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", "id", uuidParam)

			c.JSON(400, gin.H{
				"error": "not found",
			})

			return
		}

		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			c.JSON(500, gin.H{
				"error": "internal error",
			})

			return
		}

		log.Info("deleted url", slog.Any("deleted", uuidParam))

		responseOK(c)
	}
}

func NewGetAll(log *slog.Logger, get storage.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.url.getById.New"

		log := logger.FromCtx(c)

		log.Info("handling request")

		resURL, err := get.GetAllURL()
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("urls not found")

			c.JSON(400, gin.H{
				"error": "not found",
			})

			return
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			c.JSON(500, gin.H{
				"error": "internal error",
			})

			return
		}

		log.Info("got urls", slog.Any("urls", resURL))

		c.JSON(201, gin.H{
			"urls": resURL,
		})
	}
}

func NewGetById(log *slog.Logger, get storage.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.url.getById.New"

		log := logger.FromCtx(c)

		log.Info("handling request")

		uuidParam, ok := uuidparam.UUIDFromCtx(c, "id")
		if ok {
			return
		}
		resURL, err := get.GetByIdURL(uuidParam)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", slog.String("id", uuidParam.String()))

			c.JSON(400, gin.H{
				"error": "not found",
			})

			return
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			c.JSON(500, gin.H{
				"error": "internal error",
			})

			return
		}

		log.Info("got url", slog.Any("url", resURL))

		c.JSON(201, gin.H{
			"url": resURL,
		})
	}
}

func NewUpdate(log *slog.Logger, update storage.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.url.update.New"

		log := logger.FromCtx(c)

		log.Info("handling request")

		var req RequestFullStruct

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			c.JSON(400, gin.H{
				"error": "failed to decode request",
			})

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			c.JSON(400, gin.H{
				"error":   "validation failed",
				"details": formatValidationError(validateErr),
			})

			return
		}

		err := update.UpdateURL(req.Order)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", "id", req)

			c.JSON(404, gin.H{
				"error": "not found",
			})

			return
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			c.JSON(500, gin.H{
				"error": "internal error",
			})

			return
		}

		log.Info("updated url", slog.Any("url", req))

		responseOK(c)
	}
}

func NewIsOrderCreated(log *slog.Logger, ord storage.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.url.IsOrderCreated.New"

		log := logger.FromCtx(c)

		log.Info("handling request")

		var req Request

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			c.JSON(400, gin.H{
				"error": "failed to decode request",
			})

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			c.JSON(400, gin.H{
				"error":   "validation failed",
				"details": formatValidationError(validateErr),
			})

			return
		}

		uuidParam, ok := uuidparam.UUIDFromCtx(c, "id")
		if ok {
			return
		}

		resId, err := ord.IsOrderCreatedURL(uuidParam)
		if resId == false {
			log.Info("url not found", "id", resId)

			c.JSON(404, gin.H{
				"error": "not found",
			})

			return
		}

		if err != nil {
			log.Error("failed to check url", sl.Err(err))

			c.JSON(500, gin.H{
				"error": "internal error",
			})

			return
		}

		log.Info("url exist", slog.Any("url", uuidParam))

		responseOK(c)
	}
}
