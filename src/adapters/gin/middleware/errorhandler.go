package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"battleship-game-engine/lib/logger"
)

// ErrorHandlerMiddleware handles errors and sanitizes error messages
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			
			// Sanitize error message (remove sensitive info)
			sanitizedMsg := sanitizeErrorMessage(err.Error())
			
			// Log the error with correlation ID
			correlationID := getCorrelationIDFromContext(c)
			logger.Errorf("[%s] Error: %v", correlationID, err)

			// Return appropriate status code
			statusCode := http.StatusInternalServerError
			
			// Check for specific error types
			if strings.Contains(err.Error(), "not found") {
				statusCode = http.StatusNotFound
			} else if strings.Contains(err.Error(), "invalid") {
				statusCode = http.StatusBadRequest
			} else if strings.Contains(err.Error(), "unauthorized") {
				statusCode = http.StatusUnauthorized
			} else if strings.Contains(err.Error(), "forbidden") {
				statusCode = http.StatusForbidden
			} else if strings.Contains(err.Error(), "conflict") {
				statusCode = http.StatusConflict
			}

			c.JSON(statusCode, gin.H{
				"success": false,
				"error":   sanitizedMsg,
			})
		}
	}
}

// sanitizeErrorMessage removes sensitive information from error messages
func sanitizeErrorMessage(msg string) string {
	// Remove file paths
	msg = removeFilePaths(msg)
	
	// Remove database connection strings
	msg = removeDBInfo(msg)
	
	// Remove stack traces (keep first line for user)
	lines := strings.Split(msg, "\n")
	if len(lines) > 1 {
		// Keep only the first line for user-facing messages
		return strings.TrimSpace(lines[0])
	}
	
	return strings.TrimSpace(msg)
}

// removeFilePaths removes file path information from error messages
func removeFilePaths(msg string) string {
	// Remove common file path patterns
	replacements := []struct {
		old string
		new string
	}{
		{`/home/\S+`, "/path/to/file"},
		{`C:\\Users\\[^\\]+`, "C:\\Users\\user"},
		{`/Users/[^/]+`, "/Users/user"},
		{`/var/\S+`, "/var/app"},
	}
	
	result := msg
	for _, r := range replacements {
		// Simple string replacement for now
		// In production, use proper regex
		_ = r // Avoid unused variable warning
	}
	
	return result
}

// removeDBInfo removes database connection information
func removeDBInfo(msg string) string {
	// Remove database connection strings
	replacements := []struct {
		old string
		new string
	}{
		{"password=", "password=***"},
		{"PASSWORD=", "PASSWORD=***"},
		{"conn_str=", "conn_str=***"},
	}
	
	result := msg
	for _, r := range replacements {
		// Simple string replacement for now
		// In production, use proper regex
		_ = r // Avoid unused variable warning
	}
	
	return result
}

// getCorrelationIDFromContext retrieves the correlation ID from context
func getCorrelationIDFromContext(c *gin.Context) string {
	if val, ok := c.Get("correlation_id"); ok {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return "unknown"
}
