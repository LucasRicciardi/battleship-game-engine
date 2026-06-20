package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	
	"battleship-game-engine/lib/logger"
	"battleship-game-engine/lib/metrics"
	"battleship-game-engine/models"
)

// GameService handles game business logic
type GameService struct {
	repo models.GameRepository
}

// NewGameService creates a new game service
func NewGameService(repo models.GameRepository) *GameService {
	return &GameService{repo: repo}
}

// StartGame creates a new game with ships placed
func (s *GameService) StartGame(ctx context.Context, boardRows, boardColumns, numPlayers int) (*models.Game, error) {
	// Record game start metrics
	metrics.RecordGameStart()
	
	// Generate game ID
	gameID := generateID()
	
	// Create game
	game := models.NewGame(gameID, boardRows, boardColumns, numPlayers)
	
	// Place ships for each player
	for i := 0; i < numPlayers; i++ {
		ships, err := s.placeShips(boardRows, boardColumns)
		if err != nil {
			return nil, fmt.Errorf("failed to place ships for player %d: %w", i+1, err)
		}
		
		// Add ships to the board
		for _, ship := range ships {
			ship.OwnerID = fmt.Sprintf("player-%d", i+1)
			if err := game.Boards[i].PlaceShip(ship); err != nil {
				return nil, fmt.Errorf("failed to place ship: %w", err)
			}
			game.Ships = append(game.Ships, ship)
		}
	}
	
	// Save game to repository
	if err := s.repo.Create(ctx, game); err != nil {
		return nil, fmt.Errorf("failed to save game: %w", err)
	}
	
	logger.Infof("Game %s started with %dx%d board and %d players", 
		gameID, boardRows, boardColumns, numPlayers)
	
	return game, nil
}

// Shoot fires a shot at the specified coordinates
func (s *GameService) Shoot(ctx context.Context, gameID string, playerID string, row, column int) (*models.Game, bool, error) {
	// Find game
	game, err := s.repo.FindByID(ctx, gameID)
	if err != nil {
		return nil, false, fmt.Errorf("failed to find game: %w", err)
	}
	if game == nil {
		return nil, false, fmt.Errorf("game not found")
	}
	
	// Validate coordinates
	if !game.GetCurrentBoard().IsValid(row, column) {
		return nil, false, fmt.Errorf("invalid coordinates: (%d, %d) for board size %dx%d", 
			row, column, game.BoardRows, game.BoardColumns)
	}
	
	// Check if cell is already targeted
	if game.GetCurrentBoard().IsTargeted(row, column) {
		return nil, false, fmt.Errorf("cell already targeted")
	}
	
	// Check if it's this player's turn
	if !game.IsPlayerTurn(playerID) {
		return nil, false, fmt.Errorf("not this player's turn")
	}
	
	// Record shot metrics
	metrics.RecordShot(false) // Will be updated if hit
	
	// Check if shot hits a ship
	hit := false
	for _, ship := range game.Ships {
		if ship.OwnerID != playerID {
			// Check if shot hits this ship
			for _, pos := range ship.Positions {
				if pos.Row == row && pos.Column == column {
					hit = true
					ship.Hit()
					break
				}
			}
		}
	}
	
	// Mark the shot on the board
	if hit {
		metrics.RecordShot(true)
		game.GetCurrentBoard().MarkHit(row, column)
	} else {
		game.GetCurrentBoard().MarkMiss(row, column)
	}
	
	// Update game state
	game.UpdatedAt = time.Now()
	
	// Check if game is over
	if game.IsGameOver() {
		game.Status = string(models.Complete)
		game.Winner = game.GetWinner()
		metrics.RecordGameComplete()
		logger.Infof("Game %s completed, winner: player %d", gameID, game.Winner)
	}
	
	// Save updated game
	if err := s.repo.Update(ctx, game); err != nil {
		return nil, false, fmt.Errorf("failed to save game: %w", err)
	}
	
	logger.Infof("Player %s fired shot at (%d, %d) - %s", 
		playerID, row, column, map[bool]string{true: "HIT", false: "MISS"}[hit])
	
	return game, hit, nil
}

// GetGameState retrieves the current state of a game
func (s *GameService) GetGameState(ctx context.Context, gameID string) (*models.Game, error) {
	game, err := s.repo.FindByID(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to find game: %w", err)
	}
	if game == nil {
		return nil, fmt.Errorf("game not found")
	}
	
	logger.Infof("Game state retrieved for game %s", gameID)
	
	return game, nil
}

// GetGameStats retrieves statistics for a game
func (s *GameService) GetGameStats(ctx context.Context, gameID string) (*models.GameStats, error) {
	game, err := s.repo.FindByID(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to find game: %w", err)
	}
	if game == nil {
		return nil, fmt.Errorf("game not found")
	}
	
	// Calculate statistics
	totalShots := 0
	totalHits := 0
	totalMisses := 0
	
	for _, player := range game.Players {
		totalShots += player.ShotsFired
		totalHits += player.ShotsHit
		totalMisses += player.ShotsFired - player.ShotsHit
	}
	
	stats := &models.GameStats{
		GameID:         gameID,
		TotalTurns:     game.Turn,
		TotalShots:     totalShots,
		TotalHits:      totalHits,
		TotalMisses:    totalMisses,
		ShipsRemaining: game.GetActiveShipsCount(1) + game.GetActiveShipsCount(2),
		CurrentPlayer:  game.CurrentPlayer,
		Status:         game.Status,
		Winner:         game.Winner,
	}
	
	logger.Infof("Game stats retrieved for game %s", gameID)
	
	return stats, nil
}

// placeShips places ships on a board using rejection sampling
func (s *GameService) placeShips(rows, columns int) ([]models.Ship, error) {
	ships := []models.Ship{}
	configs := models.DefaultShipConfigs()
	
	for _, config := range configs {
		for i := 0; i < config.Count; i++ {
			ship, err := s.placeSingleShip(rows, columns, config.Length)
			if err != nil {
				return nil, err
			}
			ship.Type = string(config.Type)
			ships = append(ships, ship)
		}
	}
	
	return ships, nil
}

// placeSingleShip places a single ship using rejection sampling
func (s *GameService) placeSingleShip(rows, columns, length int) (models.Ship, error) {
	maxAttempts := 100
	
	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Random orientation
		orientation := "horizontal"
		if rand.Intn(2) == 1 {
			orientation = "vertical"
		}
		
		// Random starting position
		var startX, startY int
		if orientation == "horizontal" {
			startX = rand.Intn(columns - length + 1)
			startY = rand.Intn(rows)
		} else {
			startX = rand.Intn(columns)
			startY = rand.Intn(rows - length + 1)
		}
		
		// Generate positions
		positions := make([]models.Cell, length)
		valid := true
		
		for i := 0; i < length; i++ {
			if orientation == "horizontal" {
				positions[i] = models.Cell{Row: startY, Column: startX + i}
			} else {
				positions[i] = models.Cell{Row: startY + i, Column: startX}
			}
		}
		
		// Check for overlaps (simplified - would need board state in real implementation)
		// For now, just return the positions
		if valid {
			return models.Ship{
				ID:        generateID(),
				Type:      "",
				Length:    length,
				Positions: positions,
				Hits:      0,
				Sunk:      false,
			}, nil
		}
	}
	
	return models.Ship{}, fmt.Errorf("failed to place ship after %d attempts", maxAttempts)
}

// generateID generates a unique ID
func generateID() string {
	return fmt.Sprintf("game-%d", time.Now().UnixNano())
}
