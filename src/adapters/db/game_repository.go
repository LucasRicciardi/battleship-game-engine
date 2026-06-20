package db

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"battleship-game-engine/models"
)

// GameRepositoryImpl implements the GameRepository interface using GORM
type GameRepositoryImpl struct {
	db *gorm.DB
}

// NewGameRepository creates a new game repository
func NewGameRepository(db *gorm.DB) *GameRepositoryImpl {
	return &GameRepositoryImpl{db: db}
}

// Create saves a new game to the repository
func (r *GameRepositoryImpl) Create(ctx context.Context, game *models.Game) error {
	// Convert entity to DB model
	dbGame := ToGameDB(*game)
	
	// Begin transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	
	defer func() {
		if tx.Error != nil {
			tx.Rollback()
		}
	}()
	
	// Save the game
	if err := tx.Create(&dbGame).Error; err != nil {
		return fmt.Errorf("failed to create game: %w", err)
	}
	
	// Save players
	for _, player := range game.Players {
		dbPlayer := ToPlayerDB(player)
		dbPlayer.GameID = dbGame.ID
		if err := tx.Create(&dbPlayer).Error; err != nil {
			return fmt.Errorf("failed to create player: %w", err)
		}
	}
	
	// Save boards
	for _, board := range game.Boards {
		dbBoard := ToBoardDB(board)
		dbBoard.GameID = dbGame.ID
		if err := tx.Create(&dbBoard).Error; err != nil {
			return fmt.Errorf("failed to create board: %w", err)
		}
		
		// Save ships for this board
		for _, ship := range board.Ships {
			dbShip := ToShipDB(ship)
			dbShip.GameID = dbGame.ID
			dbShip.BoardID = dbBoard.ID
			if err := tx.Create(&dbShip).Error; err != nil {
				return fmt.Errorf("failed to create ship: %w", err)
			}
		}
	}
	
	// Commit transaction
	return tx.Commit().Error
}

// FindByID retrieves a game by its ID
func (r *GameRepositoryImpl) FindByID(ctx context.Context, id string) (*models.Game, error) {
	var dbGame GameDB
	if err := r.db.WithContext(ctx).First(&dbGame, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find game: %w", err)
	}
	
	// Convert to entity
	game := ToGame(dbGame)
	
	// Load players
	var dbPlayers []PlayerDB
	if err := r.db.WithContext(ctx).Where("game_id = ?", id).Find(&dbPlayers).Error; err != nil {
		return nil, fmt.Errorf("failed to load players: %w", err)
	}
	
	for _, dbPlayer := range dbPlayers {
		game.Players = append(game.Players, ToPlayer(dbPlayer))
	}
	
	// Load boards
	var dbBoards []BoardDB
	if err := r.db.WithContext(ctx).Where("game_id = ?", id).Find(&dbBoards).Error; err != nil {
		return nil, fmt.Errorf("failed to load boards: %w", err)
	}
	
	for _, dbBoard := range dbBoards {
		game.Boards = append(game.Boards, ToBoard(dbBoard))
	}
	
	// Load ships for each board
	for i := range game.Boards {
		var dbShips []ShipDB
		if err := r.db.WithContext(ctx).Where("board_id = ?", dbBoards[i].ID).Find(&dbShips).Error; err != nil {
			return nil, fmt.Errorf("failed to load ships for board %s: %w", dbBoards[i].ID, err)
		}
		
		for _, dbShip := range dbShips {
			game.Boards[i].Ships = append(game.Boards[i].Ships, ToShip(dbShip))
		}
	}
	
	return &game, nil
}

// FindAll retrieves all games
func (r *GameRepositoryImpl) FindAll(ctx context.Context) ([]models.Game, error) {
	var dbGames []GameDB
	if err := r.db.WithContext(ctx).Find(&dbGames).Error; err != nil {
		return nil, fmt.Errorf("failed to find games: %w", err)
	}
	
	games := make([]models.Game, len(dbGames))
	for i, dbGame := range dbGames {
		games[i] = ToGame(dbGame)
	}
	
	return games, nil
}

// Update saves changes to an existing game
func (r *GameRepositoryImpl) Update(ctx context.Context, game *models.Game) error {
	// Convert entity to DB model
	dbGame := ToGameDB(*game)
	
	// Begin transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	
	defer func() {
		if tx.Error != nil {
			tx.Rollback()
		}
	}()
	
	// Update the game
	if err := tx.Save(&dbGame).Error; err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}
	
	// Update players
	for _, player := range game.Players {
		dbPlayer := ToPlayerDB(player)
		if err := tx.Save(&dbPlayer).Error; err != nil {
			return fmt.Errorf("failed to update player: %w", err)
		}
	}
	
	// Update boards
	for _, board := range game.Boards {
		dbBoard := ToBoardDB(board)
		if err := tx.Save(&dbBoard).Error; err != nil {
			return fmt.Errorf("failed to update board: %w", err)
		}
		
		// Update ships for this board
		for _, ship := range board.Ships {
			dbShip := ToShipDB(ship)
			if err := tx.Save(&dbShip).Error; err != nil {
				return fmt.Errorf("failed to update ship: %w", err)
			}
		}
	}
	
	// Commit transaction
	return tx.Commit().Error
}

// Delete removes a game from the repository
func (r *GameRepositoryImpl) Delete(ctx context.Context, id string) error {
	// Begin transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	
	defer func() {
		if tx.Error != nil {
			tx.Rollback()
		}
	}()
	
	// Delete ships
	if err := tx.Where("game_id = ?", id).Delete(&ShipDB{}).Error; err != nil {
		return fmt.Errorf("failed to delete ships: %w", err)
	}
	
	// Delete boards
	if err := tx.Where("game_id = ?", id).Delete(&BoardDB{}).Error; err != nil {
		return fmt.Errorf("failed to delete boards: %w", err)
	}
	
	// Delete players
	if err := tx.Where("game_id = ?", id).Delete(&PlayerDB{}).Error; err != nil {
		return fmt.Errorf("failed to delete players: %w", err)
	}
	
	// Delete game
	if err := tx.Delete(&GameDB{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete game: %w", err)
	}
	
	// Commit transaction
	return tx.Commit().Error
}

// FindByPlayer retrieves games for a specific player
func (r *GameRepositoryImpl) FindByPlayer(ctx context.Context, playerID string) ([]models.Game, error) {
	// Find player IDs for this player
	var playerDBs []PlayerDB
	if err := r.db.WithContext(ctx).Where("id = ?", playerID).Find(&playerDBs).Error; err != nil {
		return nil, fmt.Errorf("failed to find player: %w", err)
	}
	
	if len(playerDBs) == 0 {
		return []models.Game{}, nil
	}
	
	// Find games for this player
	var gameIDs []string
	for _, p := range playerDBs {
		gameIDs = append(gameIDs, p.GameID)
	}
	
	var dbGames []GameDB
	if err := r.db.WithContext(ctx).Where("id IN ?", gameIDs).Find(&dbGames).Error; err != nil {
		return nil, fmt.Errorf("failed to find games: %w", err)
	}
	
	games := make([]models.Game, len(dbGames))
	for i, dbGame := range dbGames {
		games[i] = ToGame(dbGame)
	}
	
	return games, nil
}

// FindActiveGames returns all active games
func (r *GameRepositoryImpl) FindActiveGames(ctx context.Context) ([]models.Game, error) {
	var dbGames []GameDB
	if err := r.db.WithContext(ctx).Where("status = ?", "active").Find(&dbGames).Error; err != nil {
		return nil, fmt.Errorf("failed to find active games: %w", err)
	}
	
	games := make([]models.Game, len(dbGames))
	for i, dbGame := range dbGames {
		games[i] = ToGame(dbGame)
	}
	
	return games, nil
}
