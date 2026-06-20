package models

import "context"

// GameRepository defines the interface for game data access
type GameRepository interface {
	// Create saves a new game to the repository
	Create(ctx context.Context, game *Game) error
	
	// FindByID retrieves a game by its ID
	FindByID(ctx context.Context, id string) (*Game, error)
	
	// FindAll retrieves all games
	FindAll(ctx context.Context) ([]Game, error)
	
	// Update saves changes to an existing game
	Update(ctx context.Context, game *Game) error
	
	// Delete removes a game from the repository
	Delete(ctx context.Context, id string) error
	
	// FindByPlayer retrieves games for a specific player
	FindByPlayer(ctx context.Context, playerID string) ([]Game, error)
	
	// FindActiveGames returns all active games
	FindActiveGames(ctx context.Context) ([]Game, error)
}
