package valid

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(c *gin.Context, req interface{}, log *slog.Logger) bool {
	if err := validate.Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		if !errors.As(err, &validateErr) {
			log.Error("unknown validation error", slog.Any("error", err))
			c.JSON(500, gin.H{
				"error": "internal error"})
			return false
		}

		log.Warn("validation failed", slog.Any("errors", FormatValidationError(validateErr)))

		c.JSON(400, gin.H{
			"error":   "validation failed",
			"details": FormatValidationError(validateErr),
		})
		return false
	}
	return true
}

func FormatValidationError(err validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)
	for _, e := range err {
		errors[e.Field()] = e.Error()
	}
	return errors
}
