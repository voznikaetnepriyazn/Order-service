package bindjson

import (
	"Order/internal/lib/logger/sl"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func BindJSON(c *gin.Context, req interface{}, log *slog.Logger) bool {
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to decode request body", sl.Err(err))

		c.JSON(400, gin.H{
			"error": "failed to decode request",
		})

		return false
	}

	slog.Info("request body decoded", slog.Any("request", req))

	return true
}
