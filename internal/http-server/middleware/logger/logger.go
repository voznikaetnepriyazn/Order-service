package logger

import (
	"Order/internal/http-server/middleware"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func New(log *slog.Logger) gin.HandlerFunc {
	log = log.With(
		slog.String("component", "middleware/logger"),
	)

	log.Info("logger middleware enabled")

	fn := func(c *gin.Context) {
		entry := log.With(
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("remote_addr", c.ClientIP()),
			slog.String("user_agent", c.GetHeader("User-Agent")),
		)

		c.Set("logger", entry)

		t1 := time.Now()
		defer func() {
			entry.Info("request completed",
				slog.Int("status", c.Writer.Status()),
				slog.Int("bytes", c.Writer.Size()),
				slog.Any("duration", time.Since(t1)),
			)
		}()

		c.Next()
	}
	return fn
}

func LogQuery(log *slog.Logger, op string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(c.Request.Context())),
		)

		c.Set("logger", log)
		c.Next()
	}
}

func FromCtx(c *gin.Context) *slog.Logger {
	if v, ok := c.Get("logger"); ok {
		if log, ok := v.(*slog.Logger); ok {
			return log
		}
	}
	return slog.Default()
}
