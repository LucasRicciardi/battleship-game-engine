package db

import (
	"encoding/json"
	"battleship-game-engine/models"
)

// ToGame converts GameDB to Game entity
func ToGame(db GameDB) models.Game {
	// Parse ships JSON
	var ships []models.Ship
	if err := json.Unmarshal([]byte(db.ShipsJSON), &ships); err != nil {
		ships = []models.Ship{}
	}

	// Parse players JSON
	var players []models.Player
	if err := json.Unmarshal([]byte(db.PlayersJSON), &players); err != nil {
		players = []models.Player{}
	}

	// Parse boards JSON
	var boards []models.Board
	if err := json.Unmarshal([]byte(db.BoardsJSON), &boards); err != nil {
		boards = []models.Board{}
	}

	return models.Game{
		ID:            db.ID,
		BoardRows:     db.BoardRows,
		BoardColumns:  db.BoardColumns,
		NumPlayers:    db.NumPlayers,
		Turn:          db.Turn,
		CurrentPlayer: db.CurrentPlayer,
		Status:        db.Status,
		Winner:        db.Winner,
		Ships:         ships,
		Boards:        boards,
		Players:       players,
		CreatedAt:     db.CreatedAt,
		UpdatedAt:     db.UpdatedAt,
	}
}

// ToGameDB converts Game entity to GameDB
func ToGameDB(game models.Game) GameDB {
	// Marshal ships to JSON
	shipsJSON, _ := json.Marshal(game.Ships)

	// Marshal players to JSON
	playersJSON, _ := json.Marshal(game.Players)

	// Marshal boards to JSON
	boardsJSON, _ := json.Marshal(game.Boards)

	return GameDB{
		ID:            game.ID,
		BoardRows:     game.BoardRows,
		BoardColumns:  game.BoardColumns,
		NumPlayers:    game.NumPlayers,
		Turn:          game.Turn,
		CurrentPlayer: game.CurrentPlayer,
		Status:        game.Status,
		Winner:        game.Winner,
		ShipsJSON:     string(shipsJSON),
		PlayersJSON:   string(playersJSON),
		BoardsJSON:    string(boardsJSON),
		CreatedAt:     game.CreatedAt,
		UpdatedAt:     game.UpdatedAt,
	}
}

// ToBoard converts BoardDB to Board entity
func ToBoard(db BoardDB) models.Board {
	// Parse cells JSON
	var cells [][]models.CellState
	if err := json.Unmarshal([]byte(db.CellsJSON), &cells); err != nil {
		cells = [][]models.CellState{}
	}

	// Parse ships JSON
	var ships []models.Ship
	if err := json.Unmarshal([]byte(db.ShipsJSON), &ships); err != nil {
		ships = []models.Ship{}
	}

	return models.Board{
		Rows:      db.Rows,
		Columns:   db.Columns,
		Cells:     cells,
		Ships:     ships,
		OwnerID:   db.OwnerID,
	}
}

// ToBoardDB converts Board entity to BoardDB
func ToBoardDB(board models.Board) BoardDB {
	// Marshal cells to JSON
	cellsJSON, _ := json.Marshal(board.Cells)

	// Marshal ships to JSON
	shipsJSON, _ := json.Marshal(board.Ships)

	return BoardDB{
		Rows:      board.Rows,
		Columns:   board.Columns,
		CellsJSON: string(cellsJSON),
		ShipsJSON: string(shipsJSON),
		OwnerID:   board.OwnerID,
	}
}

// ToShip converts ShipDB to Ship entity
func ToShip(db ShipDB) models.Ship {
	// Parse positions JSON
	var positions []models.Cell
	if err := json.Unmarshal([]byte(db.Positions), &positions); err != nil {
		positions = []models.Cell{}
	}

	return models.Ship{
		ID:        db.ID,
		Type:      db.Type,
		Length:    db.Length,
		Positions: positions,
		Hits:      db.Hits,
		Sunk:      db.Sunk,
		OwnerID:   db.OwnerID,
	}
}

// ToShipDB converts Ship entity to ShipDB
func ToShipDB(ship models.Ship) ShipDB {
	// Marshal positions to JSON
	positionsJSON, _ := json.Marshal(ship.Positions)

	return ShipDB{
		ID:        ship.ID,
		Type:      ship.Type,
		Length:    ship.Length,
		Positions: string(positionsJSON),
		Hits:      ship.Hits,
		Sunk:      ship.Sunk,
		OwnerID:   ship.OwnerID,
	}
}

// ToPlayer converts PlayerDB to Player entity
func ToPlayer(db PlayerDB) models.Player {
	return models.Player{
		ID:         db.ID,
		Number:     db.Number,
		Name:       db.Name,
		ShotsFired: db.ShotsFired,
		ShotsHit:   db.ShotsHit,
		ShipsSunk:  db.ShipsSunk,
	}
}

// ToPlayerDB converts Player entity to PlayerDB
func ToPlayerDB(player models.Player) PlayerDB {
	return PlayerDB{
		ID:         player.ID,
		Number:     player.Number,
		Name:       player.Name,
		ShotsFired: player.ShotsFired,
		ShotsHit:   player.ShotsHit,
		ShipsSunk:  player.ShipsSunk,
	}
}
