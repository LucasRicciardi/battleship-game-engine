package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/memory"
)

// RateLimitMiddleware creates a rate limiter middleware
func RateLimitMiddleware() gin.HandlerFunc {
	// Create memory store for rate limiting
	store := memory.NewStore()

	// Create rate limiter configuration
	rate := limiter.Rate{
		Period: 60 * time.Second,
		Limit:  100, // 100 requests per minute
	}

	// Create limiter instance
	limit := limiter.New(store, rate)

	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		// Check rate limit
		if err := limit.Allow(clientIP); err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
