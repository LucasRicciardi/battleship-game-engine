package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"battleship-game-engine/lib/logger"
)

// PlayerAuthorizationMiddleware checks if the player is authorized for the game
func PlayerAuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get game ID from path parameter
		gameID := c.Param("game_id")
		if gameID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Game ID required",
			})
			c.Abort()
			return
		}

		// Get player ID from context (set by JWT middleware)
		playerID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Player authentication required",
			})
			c.Abort()
			return
		}

		// TODO: Verify player is authorized for this game
		// This would typically involve:
		// 1. Querying the database for the game
		// 2. Checking if the player is in the game's player list
		// 3. Checking if it's this player's turn (for multiplayer games)

		// For now, just log the authorization check
		correlationID := getCorrelationID(c)
		logger.Infof("[%s] Authorization check: player %v accessing game %s", 
			correlationID, playerID, gameID)

		c.Next()
	}
}

// IsPlayerTurnMiddleware checks if it's the player's turn
func IsPlayerTurnMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get game ID and player ID
		gameID := c.Param("game_id")
		playerID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Player authentication required",
			})
			c.Abort()
			return
		}

		// TODO: Check if it's this player's turn
		// This would involve querying the database for the current turn

		// For now, just log the turn check
		correlationID := getCorrelationID(c)
		logger.Infof("[%s] Turn check: player %v for game %s", 
			correlationID, playerID, gameID)

		c.Next()
	}
}

// getCorrelationID retrieves the correlation ID from context
func getCorrelationID(c *gin.Context) string {
	if val, ok := c.Get("correlation_id"); ok {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return "unknown"
}
