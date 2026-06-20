package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

// correlationIDKey is the context key for correlation ID
type correlationIDKey struct{}

// CorrelationIDMiddleware creates and propagates correlation IDs
func CorrelationIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a new correlation ID
		correlationID := uuid.New().String()
		
		// Store in context
		ctx := context.WithValue(c.Request.Context(), correlationIDKey{}, correlationID)
		c.Request = c.Request.WithContext(ctx)
		
		// Set in response header
		c.Writer.Header().Set("X-Correlation-ID", correlationID)
		
		// Store in context for access in handlers
		c.Set("correlation_id", correlationID)
		
		c.Next()
	}
}

// GetCorrelationID retrieves the correlation ID from context
func GetCorrelationID(c *gin.Context) string {
	if val, ok := c.Get("correlation_id"); ok {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}

// GetCorrelationIDFromContext retrieves the correlation ID from a context
func GetCorrelationIDFromContext(ctx context.Context) string {
	if val := ctx.Value(correlationIDKey{}); val != nil {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}
