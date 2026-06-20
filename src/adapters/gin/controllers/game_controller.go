package controllers

import (
	"fmt"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"battleship-game-engine/src/adapters/gin/helpers"
	"battleship-game-engine/src/services"
)

// GameController handles game-related HTTP requests
type GameController struct {
	gameService *services.GameService
}

// NewGameController creates a new game controller
func NewGameController(gameService *services.GameService) *GameController {
	return &GameController{gameService: gameService}
}

// StartGame handles POST /api/v1/games
func (c *GameController) StartGame(ctx *gin.Context) {
	// Parse request body
	var req struct {
		BoardRows    int `json:"board_rows" binding:"required,min=5,max=100"`
		BoardColumns int `json:"board_columns" binding:"required,min=5,max=100"`
		NumPlayers   int `json:"num_players" binding:"required,min=1,max=2"`
	}
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}
	
	// Start the game
	game, err := c.gameService.StartGame(ctx, req.BoardRows, req.BoardColumns, req.NumPlayers)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	// Return game state
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    game,
	})
}

// Shoot handles POST /api/v1/games/:game_id/shoot
func (c *GameController) Shoot(ctx *gin.Context) {
	gameID := ctx.Param("game_id")
	
	// Parse request body
	var req struct {
		PlayerID string `json:"player_id" binding:"required,uuid4"`
		Row      int    `json:"row" binding:"required,min=0"`
		Column   int    `json:"column" binding:"required,min=0"`
	}
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}
	
	// Get game state first to get board dimensions
	game, err := c.gameService.GetGameState(ctx, gameID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Game not found",
		})
		return
	}
	
	// Validate coordinates against board size
	if req.Row < 0 || req.Row >= game.BoardRows {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   fmt.Sprintf("row must be between 0 and %d", game.BoardRows-1),
		})
		return
	}
	if req.Column < 0 || req.Column >= game.BoardColumns {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   fmt.Sprintf("column must be between 0 and %d", game.BoardColumns-1),
		})
		return
	}
	
	// Fire the shot
	game, hit, err := c.gameService.Shoot(ctx, gameID, req.PlayerID, req.Row, req.Column)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	// Prepare response
	response := gin.H{
		"success": true,
		"data": gin.H{
			"hit":            hit,
			"ships_remaining": game.GetActiveShipsCount(1) + game.GetActiveShipsCount(2),
			"board":          game.GetCurrentBoard().Cells,
			"current_turn":   game.CurrentPlayer,
		},
	}
	
	// Add winner if game is over
	if game.Status == "complete" {
		response["data"].(gin.H)["winner"] = game.Winner
	}
	
	ctx.JSON(http.StatusOK, response)
}

// GetGameState handles GET /api/v1/games/:game_id
func (c *GameController) GetGameState(ctx *gin.Context) {
	gameID := ctx.Param("game_id")
	
	game, err := c.gameService.GetGameState(ctx, gameID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Game not found",
		})
		return
	}
	
	// Create board visualizer
	visualizer := helpers.NewBoardVisualizer(nil, game.BoardRows, game.BoardColumns)
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    game,
	})
}

// GetGameStats handles GET /api/v1/games/:game_id/stats
func (c *GameController) GetGameStats(ctx *gin.Context) {
	gameID := ctx.Param("game_id")
	
	stats, err := c.gameService.GetGameStats(ctx, gameID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Game not found",
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}
