package gin

import (
	"github.com/gin-gonic/gin"
	"battleship-game-engine/src/adapters/gin/handlers"
	"battleship-game-engine/src/adapters/gin/middleware"
)

// NewRouter creates and configures the Gin router
func NewRouter() *gin.Engine {
	r := gin.New()
	
	// Global middleware
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.RateLimitMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())
	r.Use(middleware.CorrelationIDMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	
	// Health check routes (public)
	healthController := handlers.NewHealthController()
	r.GET("/health/live", healthController.Liveness)
	r.GET("/health/ready", healthController.Readiness)
	
	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			// Login endpoint
			auth.POST("/login", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Login endpoint - to be implemented"})
			})
			// Register endpoint
			auth.POST("/register", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Register endpoint - to be implemented"})
			})
		}
		
		// Protected routes (auth required)
		protected := v1.Group("/")
		protected.Use(middleware.JWTMiddleware{SigningKey: []byte("your-secret-key")}.Middleware())
		{
			// Games routes
			games := protected.Group("/games")
			{
				// Start a new game
				games.POST("", middleware.PlayerAuthorizationMiddleware(), func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "Start game endpoint - to be implemented"})
				})
				
				// Get game state
				games.GET("/:game_id", middleware.PlayerAuthorizationMiddleware(), func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "Get game state endpoint - to be implemented"})
				})
				
				// Fire a shot
				games.POST("/:game_id/shoot", middleware.PlayerAuthorizationMiddleware(), middleware.IsPlayerTurnMiddleware(), func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "Shoot endpoint - to be implemented"})
				})
				
				// Get game statistics
				games.GET("/:game_id/stats", middleware.PlayerAuthorizationMiddleware(), func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "Get game stats endpoint - to be implemented"})
				})
			}
		}
	}
	
	return r
}
