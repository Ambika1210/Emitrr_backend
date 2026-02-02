package game

// BotPlayer represents the AI opponent
type BotPlayer struct {
	Level int // For future scalability
}

// SelectBestMove analyzes the board and selects the most strategic column
func (gs *GameState) SelectBestMove() int {
	// Priority 1: Can I win in one move?
	for col := 0; col < Cols; col++ {
		if gs.canWinInMove(col, gs.CurrentPlayer) {
			return col
		}
	}

	// Priority 2: Can my opponent win in their next move? (Block them)
	opponent := Player1
	if gs.CurrentPlayer == Player1 {
		opponent = Player2
	}
	for col := 0; col < Cols; col++ {
		if gs.canWinInMove(col, opponent) {
			return col
		}
	}

	// Priority 3: Try to take the center column
	center := Cols / 2
	if gs.Board[0][center] == Empty {
		return center
	}

	// Priority 4: Pick the first available column
	for col := 0; col < Cols; col++ {
		if gs.Board[0][col] == Empty {
			return col
		}
	}

	return -1
}

// canWinInMove simulates a move to see if it results in a win
func (gs *GameState) canWinInMove(col int, player int) bool {
	if col < 0 || col >= Cols || gs.Board[0][col] != Empty {
		return false
	}

	// Find the row where the disc would fall
	row := -1
	for r := Rows - 1; r >= 0; r-- {
		if gs.Board[r][col] == Empty {
			row = r
			break
		}
	}

	if row == -1 {
		return false
	}

	// Temporarily make the move
	originalPlayer := gs.CurrentPlayer
	gs.Board[row][col] = player

	isWin := gs.DetectWinner(row, col)

	// Undo the move
	gs.Board[row][col] = Empty
	gs.CurrentPlayer = originalPlayer

	return isWin
}
