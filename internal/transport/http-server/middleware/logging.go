package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error("request error",
					slog.String("method", method),
					slog.String("path", path),
					slog.Int("status", statusCode),
					slog.Duration("duration", duration),
					slog.String("client_ip", clientIP),
					slog.String("error", e),
				)
			}
		} else {
			logger.Info("request completed",
				slog.String("method", method),
				slog.String("path", path),
				slog.Int("status", statusCode),
				slog.Duration("duration", duration),
				slog.String("client_ip", clientIP),
			)
		}
	}
}
