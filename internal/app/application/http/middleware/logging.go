package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/logger"
)

// loggerKey is the Gin context key used to store/retrieve the request-scoped logger.
const loggerKey = "logger"

// LoggerMiddleware injects a request-scoped logger into the Gin context.
// It enriches the logger with basic request metadata and attaches the request context.
func LoggerMiddleware(base logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = ulid.Make().String()
		}

		l := base.With(
			"request_id", reqID,
			"method", c.Request.Method,
			"path", c.FullPath(),
		).WithContext(c.Request.Context())

		Set(c, l)

		c.Next()
	}
}

// Get returns the request-scoped logger from the Gin context.
func Get(c *gin.Context) (logger.Logger, bool) {
	if v, ok := c.Get(loggerKey); ok {
		if l, ok2 := v.(logger.Logger); ok2 {
			return l, true
		}
	}
	return nil, false
}

// Set stores the supplied logger into the Gin context.
func Set(c *gin.Context, l logger.Logger) {
	c.Set(loggerKey, l)
}
