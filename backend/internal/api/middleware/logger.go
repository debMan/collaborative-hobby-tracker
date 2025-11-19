package middleware

import (
	"time"

	"github.com/debMan/collaborative-hobby-tracker/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Logger returns a gin middleware for logging HTTP requests
func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Log request
		log.Infow("HTTP Request",
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"status", statusCode,
			"latency", latency,
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)
	}
}
