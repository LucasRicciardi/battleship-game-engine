package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// SimpleRateLimiter implements a basic in-memory rate limiter
type SimpleRateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

// NewSimpleRateLimiter creates a new rate limiter
func NewSimpleRateLimiter(limit int, window time.Duration) *SimpleRateLimiter {
	return &SimpleRateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow checks if a request is allowed for the given client IP
func (rl *SimpleRateLimiter) Allow(clientIP string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Get or create request history for this client
	times, exists := rl.requests[clientIP]
	if !exists {
		times = []time.Time{}
	}

	// Filter out old requests outside the window
	var validTimes []time.Time
	for _, t := range times {
		if t.After(windowStart) {
			validTimes = append(validTimes, t)
		}
	}

	// Check if limit exceeded
	if len(validTimes) >= rl.limit {
		rl.requests[clientIP] = validTimes
		return &RateLimitError{RetryAfter: rl.window}
	}

	// Add current request
	validTimes = append(validTimes, now)
	rl.requests[clientIP] = validTimes
	return nil
}

// RateLimitError represents a rate limit exceeded error
type RateLimitError struct {
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	return "Rate limit exceeded"
}

// RateLimitMiddleware creates a rate limiter middleware
func RateLimitMiddleware() gin.HandlerFunc {
	// Create rate limiter: 100 requests per minute
	limiter := NewSimpleRateLimiter(100, 60*time.Second)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if err := limiter.Allow(clientIP); err != nil {
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
