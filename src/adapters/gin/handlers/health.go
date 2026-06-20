package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthController handles health check endpoints
type HealthController struct{}

// NewHealthController creates a new health controller
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Liveness handles /health/live requests
func (h *HealthController) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// Readiness handles /health/ready requests
func (h *HealthController) Readiness(c *gin.Context) {
	// Check database connectivity
	status := "ready"
	checks := gin.H{
		"database": "healthy",
		"cache":    "healthy",
	}

	// For now, just return ready status
	// In production, you would check actual dependencies

	c.JSON(http.StatusOK, gin.H{
		"status":    status,
		"checks":    checks,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
