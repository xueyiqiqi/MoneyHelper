package middleware

import (
	"time"

	"life-financial-assistant-backend/internal/logger"
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := int(time.Since(start).Milliseconds())
		logger.LogRequest(c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)
	}
}