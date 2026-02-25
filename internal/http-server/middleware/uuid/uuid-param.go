package uuidparam

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UUIDParam(paramName string, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		alias := c.Param("paramName")
		if alias == "" {
			slog.Info("id is empty")

			c.JSON(400, gin.H{
				"error": "invalid request",
			})
			c.Abort()
			return
		}

		id, err := uuid.Parse(alias)
		if err != nil {
			slog.Info("invalid uuid format", slog.String("id", alias), slog.Any("error", err))
			c.JSON(400, gin.H{
				"error": "invalid id format",
			})
			c.Abort()
			return
		}

		c.Set("uuid_"+paramName, id)
		c.Next()
	}
}

func UUIDFromCtx(c *gin.Context, paramName string) (uuid.UUID, bool) {
	if v, ok := c.Get("uuid_" + paramName); ok {
		if id, ok := v.(uuid.UUID); ok {
			return id, true
		}
	}
	return uuid.Nil, false
}
