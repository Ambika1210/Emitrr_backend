package game

import "errors"

const (
	Rows    = 6
	Cols    = 7
	Empty   = 0
	Player1 = 1
	Player2 = 2
)

var (
	ErrInvalidColumn = errors.New("invalid column")
	ErrColumnFull    = errors.New("column is full")
	ErrGameOver      = errors.New("game is over")
)

type Board [Rows][Cols]int

type GameState struct {
	Board         Board `json:"board"`
	CurrentPlayer int   `json:"currentPlayer"`
	Winner        int   `json:"winner"`
	IsDraw        bool  `json:"isDraw"`
	IsFinished    bool  `json:"isFinished"`
}

// NewConnectFourGame initializes a new 7x6 board
func NewConnectFourGame() *GameState {
	return &GameState{
		CurrentPlayer: Player1,
	}
}

// ExecutePlayerMove handles dropping a disc into a column
func (gs *GameState) ExecutePlayerMove(col int) error {
	if gs.IsFinished {
		return ErrGameOver
	}

	if col < 0 || col >= Cols {
		return ErrInvalidColumn
	}

	// Logic to drop disc to the bottom
	rowFound := -1
	for r := Rows - 1; r >= 0; r-- {
		if gs.Board[r][col] == Empty {
			rowFound = r
			break
		}
	}

	if rowFound == -1 {
		return ErrColumnFull
	}

	gs.Board[rowFound][col] = gs.CurrentPlayer

	// Post-move checks
	if gs.DetectWinner(rowFound, col) {
		gs.Winner = gs.CurrentPlayer
		gs.IsFinished = true
	} else if gs.CheckForDraw() {
		gs.IsDraw = true
		gs.IsFinished = true
	} else {
		gs.SwitchTurns()
	}

	return nil
}

// DetectWinner checks all 4 directions for 4-in-a-row
func (gs *GameState) DetectWinner(row, col int) bool {
	player := gs.Board[row][col]
	directions := [][2]int{
		{0, 1},  // Horizontal
		{1, 0},  // Vertical
		{1, 1},  // Diagonal \
		{1, -1}, // Diagonal /
	}

	for _, d := range directions {
		count := 1
		// Check forward
		for i := 1; i < 4; i++ {
			r, c := row+d[0]*i, col+d[1]*i
			if r >= 0 && r < Rows && c >= 0 && c < Cols && gs.Board[r][c] == player {
				count++
			} else {
				break
			}
		}
		// Check backward
		for i := 1; i < 4; i++ {
			r, c := row-d[0]*i, col-d[1]*i
			if r >= 0 && r < Rows && c >= 0 && c < Cols && gs.Board[r][c] == player {
				count++
			} else {
				break
			}
		}
		if count >= 4 {
			return true
		}
	}
	return false
}

// CheckForDraw returns true if the board is full
func (gs *GameState) CheckForDraw() bool {
	for c := 0; c < Cols; c++ {
		if gs.Board[0][c] == Empty {
			return false
		}
	}
	return true
}

// SwitchTurns changes CurrentPlayer from 1 to 2 or 2 to 1
func (gs *GameState) SwitchTurns() {
	if gs.CurrentPlayer == Player1 {
		gs.CurrentPlayer = Player2
	} else {
		gs.CurrentPlayer = Player1
	}
}
