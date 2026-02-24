package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// генерация ид запроса
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := fmt.Sprintf("%d", time.Now().UnixNano())
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "requestID", id))
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

// перехват паники - обертка запроса в дефер и рековер
func Recoverer(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stackBuf := make([]byte, 4096)
				stackSize := runtime.Stack(stackBuf, false)
				stackTrace := string(stackBuf[:stackSize])

				log.Error("panic recovered",
					slog.Any("error", err),
					slog.String("stack_trace", stackTrace),
					slog.String("method", c.Request.Method),
					slog.String("path", c.Request.URL.Path),
				)

				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}

			c.Abort()
		}()

		c.Next()
	}
}

// получение ид из контекста
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if reqID, ok := ctx.Value("requestID").(string); ok {
		return reqID
	}

	return ""
}
