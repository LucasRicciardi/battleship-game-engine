package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
)

// JWTMiddleware handles JWT authentication
type JWTMiddleware struct {
	SigningKey []byte
}

// NewJWTMiddleware creates a new JWT middleware
func NewJWTMiddleware(signingKey string) *JWTMiddleware {
	return &JWTMiddleware{
		SigningKey: []byte(signingKey),
	}
}

// Middleware returns the gin.HandlerFunc for JWT authentication
func (m *JWTMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"success": false,
				"error":   "Authorization header required",
			})
			c.Abort()
			return
		}

		// Parse and validate token
		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return m.SigningKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{
				"success": false,
				"error":   "Invalid token",
			})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{
				"success": false,
				"error":   "Invalid token claims",
			})
			c.Abort()
			return
		}

		// Validate expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				c.JSON(401, gin.H{
					"success": false,
					"error":   "Token expired",
				})
				c.Abort()
				return
			}
		}

		// Store user ID in context
		if userID, ok := claims["user_id"].(string); ok {
			c.Set("user_id", userID)
		}

		c.Next()
	}
}
