package models

// Ship represents a ship in the Battleship game
type Ship struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Length    int    `json:"length"`
	Positions []Cell `json:"positions"`
	Hits      int    `json:"hits"`
	Sunk      bool   `json:"sunk"`
	OwnerID   string `json:"owner_id,omitempty"`
}

// Cell represents a cell position on the board
type Cell struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}

// ShipType represents the type of ship
type ShipType string

const (
	// Ship types
	Destroyer  ShipType = "destroyer"
	Cruiser    ShipType = "cruiser"
	Battleship ShipType = "battleship"
	Aircraft   ShipType = "aircraft"
)

// ShipConfig represents the configuration for a ship type
type ShipConfig struct {
	Type   ShipType
	Length int
	Count  int
}

// DefaultShipConfigs returns the default ship configurations for a standard game
func DefaultShipConfigs() []ShipConfig {
	return []ShipConfig{
		{Type: Destroyer, Length: 2, Count: 1},
		{Type: Cruiser, Length: 3, Count: 1},
		{Type: Battleship, Length: 4, Count: 1},
		{Type: Aircraft, Length: 5, Count: 1},
	}
}

// IsSunk returns true if the ship is sunk (all positions hit)
func (s *Ship) IsSunk() bool {
	return s.Hits >= s.Length
}

// Hit marks a position on the ship as hit
func (s *Ship) Hit() {
	s.Hits++
	s.Sunk = s.IsSunk()
}
