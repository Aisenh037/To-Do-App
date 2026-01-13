package middleware

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var Logger *slog.Logger

func init() {
	// Initialize structured logger
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(Logger)
}

// StructuredLoggerMiddleware provides structured logging for HTTP requests
func StructuredLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()

		// Log with structured fields
		Logger.Info("HTTP Request",
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", statusCode),
			slog.Duration("latency", latency),
			slog.String("client_ip", clientIP),
			slog.String("user_agent", userAgent),
			slog.Int("body_size", c.Writer.Size()),
		)

		// Log errors separately
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				Logger.Error("Request Error",
					slog.String("error", e.Error()),
					slog.String("path", path),
				)
			}
		}
	}
}
