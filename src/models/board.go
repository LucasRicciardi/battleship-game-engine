package models

import (
	"fmt"
)

// Board represents the game board
type Board struct {
	Rows    int           `json:"rows"`
	Columns int           `json:"columns"`
	Cells   [][]CellState `json:"cells"`
	Ships   []Ship        `json:"ships,omitempty"`
	OwnerID string        `json:"owner_id,omitempty"`
}

// CellState represents the state of a cell on the board
type CellState string

const (
	// Cell states
	Empty    CellState = " " // Untargeted
	Miss     CellState = "O" // Missed shot
	Hit      CellState = "X" // Hit shot
	Ship     CellState = "S" // Ship position (hidden in opponent view)
	ShipHit  CellState = "*" // Ship position that was hit
	ShipMiss CellState = "O" // Ship position that was missed (same as miss)
)

// NewBoard creates a new board with the specified dimensions
func NewBoard(rows, columns int) *Board {
	cells := make([][]CellState, rows)
	for i := range cells {
		cells[i] = make([]CellState, columns)
		for j := range cells[i] {
			cells[i][j] = Empty
		}
	}

	return &Board{
		Rows:    rows,
		Columns: columns,
		Cells:   cells,
		Ships:   []Ship{},
	}
}

// IsValid checks if the given coordinates are valid for this board
func (b *Board) IsValid(row, column int) bool {
	return row >= 0 && row < b.Rows && column >= 0 && column < b.Columns
}

// IsTargeted checks if a cell has already been targeted
func (b *Board) IsTargeted(row, column int) bool {
	if !b.IsValid(row, column) {
		return false
	}
	return b.Cells[row][column] != Empty
}

// MarkMiss marks a cell as a miss
func (b *Board) MarkMiss(row, column int) {
	if b.IsValid(row, column) {
		b.Cells[row][column] = Miss
	}
}

// MarkHit marks a cell as a hit
func (b *Board) MarkHit(row, column int) {
	if b.IsValid(row, column) {
		b.Cells[row][column] = Hit
	}
}

// PlaceShip places a ship on the board
func (b *Board) PlaceShip(ship Ship) error {
	// Check if all positions are valid
	for _, pos := range ship.Positions {
		if !b.IsValid(pos.Row, pos.Column) {
			return fmt.Errorf("invalid position: (%d, %d)", pos.Row, pos.Column)
		}

		// Check if position is already occupied
		if b.Cells[pos.Row][pos.Column] != Empty {
			return fmt.Errorf("position (%d, %d) is already occupied", pos.Row, pos.Column)
		}
	}

	// Place the ship
	for _, pos := range ship.Positions {
		b.Cells[pos.Row][pos.Column] = Ship
	}

	// Add ship to the board's ship list
	b.Ships = append(b.Ships, ship)

	return nil
}

// GetCell returns the state of a cell
func (b *Board) GetCell(row, column int) CellState {
	if b.IsValid(row, column) {
		return b.Cells[row][column]
	}
	return Empty
}

// GetShips returns all ships on the board
func (b *Board) GetShips() []Ship {
	return b.Ships
}

// GetActiveShips returns only ships that are not sunk
func (b *Board) GetActiveShips() []Ship {
	active := []Ship{}
	for _, ship := range b.Ships {
		if !ship.Sunk {
			active = append(active, ship)
		}
	}
	return active
}

// GetActiveShipsCount returns the number of active ships
func (b *Board) GetActiveShipsCount() int {
	return len(b.GetActiveShips())
}

// Display returns a string representation of the board for debugging
func (b *Board) Display() string {
	var result string

	// Column headers
	result += "   "
	for i := 0; i < b.Columns; i++ {
		result += fmt.Sprintf("%2d ", i)
	}
	result += "\n"

	// Board rows
	for i := 0; i < b.Rows; i++ {
		result += fmt.Sprintf("%2d ", i)
		for j := 0; j < b.Columns; j++ {
			result += fmt.Sprintf(" %s ", b.Cells[i][j])
		}
		result += "\n"
	}

	return result
}
