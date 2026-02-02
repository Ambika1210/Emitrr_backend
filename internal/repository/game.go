package repository

import (
	"context"
	"encoding/json"

	"github.com/yourusername/connect-four-backend/internal/db"
	"github.com/yourusername/connect-four-backend/internal/logger"
)

type GameRecord struct {
	Player1  string
	Player2  string
	Winner   string
	Duration int
	State    interface{}
}

func SaveGameResult(record GameRecord) error {
	const query = `
		INSERT INTO games (player1_id, player2_id, winner_id, duration_seconds, game_state)
		VALUES ($1, $2, $3, $4, $5)`

	stateJSON, _ := json.Marshal(record.State)

	winner := record.Winner
	if winner == "" {
		// handle draw null if needed, but our table allows null
	}

	_, err := db.Pool.Exec(context.Background(), query,
		record.Player1, record.Player2, winner, record.Duration, stateJSON)

	if err != nil {
		logger.Error("game.go >>>> SaveGameResult >>>>> Failed to save game", err)
	}
	return err
}
