package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware adds security headers to responses
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")
		
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		
		// Enable XSS filter
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Strict transport security (HSTS)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'")
		
		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin")
		
		// Permissions Policy
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		c.Next()
	}
}
