package db

import (
	"time"

	"gorm.io/gorm"
)

// GameDB represents the database model for a game
type GameDB struct {
	ID            string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	BoardRows     int    `gorm:"not null;default:8"`
	BoardColumns  int    `gorm:"not null;default:8"`
	NumPlayers    int    `gorm:"not null;default:1"`
	Turn          int    `gorm:"not null;default:1"`
	CurrentPlayer int    `gorm:"not null;default:1"`
	Status        string `gorm:"not null;default:'active'"`
	Winner        int    `gorm:"default:0"`
	ShipsJSON     string `gorm:"type:text"`   // JSON representation of ships
	PlayersJSON   string `gorm:"type:text"`   // JSON representation of players
	BoardsJSON    string `gorm:"type:text"`   // JSON representation of boards
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// BoardDB represents the database model for a board
type BoardDB struct {
	ID         string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	GameID     string `gorm:"not null;index"`
	OwnerID    string `gorm:"not null"`
	Rows       int    `gorm:"not null"`
	Columns    int    `gorm:"not null"`
	CellsJSON  string `gorm:"type:text"` // JSON representation of cells
	ShipsJSON  string `gorm:"type:text"` // JSON representation of ships
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// ShipDB represents the database model for a ship
type ShipDB struct {
	ID        string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	GameID    string `gorm:"not null;index"`
	BoardID   string `gorm:"not null;index"`
	OwnerID   string `gorm:"not null"`
	Type      string `gorm:"not null"`
	Length    int    `gorm:"not null"`
	Positions string `gorm:"type:text"` // JSON array of positions
	Hits      int    `gorm:"not null;default:0"`
	Sunk      bool   `gorm:"not null;default:false"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// PlayerDB represents the database model for a player
type PlayerDB struct {
	ID         string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	GameID     string `gorm:"not null;index"`
	Number     int    `gorm:"not null"`
	Name       string `gorm:"not null"`
	ShotsFired int    `gorm:"not null;default:0"`
	ShotsHit   int    `gorm:"not null;default:0"`
	ShipsSunk  int    `gorm:"not null;default:0"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// GameStatsDB represents the database model for game statistics
type GameStatsDB struct {
	ID             string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	GameID         string `gorm:"not null;uniqueIndex"`
	TotalTurns     int    `gorm:"not null;default:0"`
	TotalShots     int    `gorm:"not null;default:0"`
	TotalHits      int    `gorm:"not null;default:0"`
	TotalMisses    int    `gorm:"not null;default:0"`
	ShipsRemaining int    `gorm:"not null;default:0"`
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

// TableName returns the table name for GameDB
func (GameDB) TableName() string {
	return "games"
}

// TableName returns the table name for BoardDB
func (BoardDB) TableName() string {
	return "boards"
}

// TableName returns the table name for ShipDB
func (ShipDB) TableName() string {
	return "ships"
}

// TableName returns the table name for PlayerDB
func (PlayerDB) TableName() string {
	return "players"
}

// TableName returns the table name for GameStatsDB
func (GameStatsDB) TableName() string {
	return "game_stats"
}
