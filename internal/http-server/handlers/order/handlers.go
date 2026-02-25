package handlers

import (
	"errors"
	"log/slog"

	"Order/internal/http-server/middleware/logger"
	resp "Order/internal/lib/api/response"
	"Order/internal/lib/logger/sl"
	"Order/internal/models/order"
	"Order/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

		alias := c.Param("id")
		if alias == "" {
			log.Info("id is empty")

			c.JSON(400, gin.H{
				"error": "invalid request",
			})
			return
		}

		id, err := uuid.Parse(alias)
		if err != nil {
			log.Info("invalid uuid format", slog.String("id", alias), slog.Any("error", err))
			c.JSON(400, gin.H{
				"error": "invalid id format",
			})
			return
		}

		err1 := deleter.DeleteURL(id)
		if errors.Is(err1, storage.ErrUrlNotFound) {
			log.Info("url not found", "id", alias)

			c.JSON(400, gin.H{
				"error": "not found",
			})

			return
		}

		if err1 != nil {
			log.Error("failed to delete url", sl.Err(err))

			c.JSON(500, gin.H{
				"error": "internal error",
			})

			return
		}

		log.Info("deleted url", slog.String("deleted", alias))

		responseOK(c)
	}
}

func NewGetAll(log *slog.Logger, get storage.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.url.getById.New"

		log := logger.FromCtx(c)

		log.Info("handling request")

		var req Request

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			c.JSON(400, gin.H{
				"error":   "validation failed",
				"details": formatValidationError(validateErr),
			})

			return
		}

		ids := c.Param("ids")
		if ids == "" {
			log.Info("ids is empty")

			c.JSON(400, gin.H{
				"error": "invalid request",
			})
			return
		}

		resURL, err := get.GetAllURL()
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("urls not found", "ids", ids)

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

		var req RequestFullStruct

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			c.JSON(400, gin.H{
				"error":   "validation failed",
				"details": formatValidationError(validateErr),
			})

			return
		}

		alias := c.Param("id")
		if alias == "" {
			log.Info("id is empty")

			c.JSON(400, gin.H{
				"error": "invalid request",
			})
			return
		}

		id, err := uuid.Parse(alias)
		if err != nil {
			log.Info("invalid uuid format", slog.String("id", alias), slog.Any("error", err))
			c.JSON(400, gin.H{
				"error": "invalid id format",
			})
			return
		}

		resURL, err := get.GetByIdURL(id)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", "id", id)

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

		alias := c.Param("id")
		if alias == "" {
			log.Info("id is empty")

			c.JSON(400, gin.H{
				"error": "invalid request",
			})
			return
		}

		id, err := uuid.Parse(alias)
		if err != nil {
			log.Info("invalid uuid format", slog.String("id", alias), slog.Any("error", err))
			c.JSON(400, gin.H{
				"error": "invalid id format",
			})
			return
		}

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

		resId, err := ord.IsOrderCreatedURL(id)
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

		log.Info("url exist", slog.Any("url", id))

		responseOK(c)
	}
}
