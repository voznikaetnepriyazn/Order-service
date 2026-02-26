package valid

import (
	"Order/internal/lib/logger/sl"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Validate(req interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			slog.Error("invalid request", sl.Err(err))

			c.JSON(400, gin.H{
				"error":   "validation failed",
				"details": FormatValidationError(validateErr),
			})

			c.Abort()

			return
		}
		c.Next()
	}
}

func FormatValidationError(err validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)
	for _, e := range err {
		errors[e.Field()] = e.Error()
	}
	return errors
}
