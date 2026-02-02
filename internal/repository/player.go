package repository

import (
	"context"

	"github.com/yourusername/connect-four-backend/internal/db"
	"github.com/yourusername/connect-four-backend/internal/logger"
)

type PlayerStats struct {
	Username string `json:"username"`
	Wins     int    `json:"wins"`
}

func UpdatePlayerWin(username string) error {
	const query = `
		INSERT INTO players (username, wins) 
		VALUES ($1, 1) 
		ON CONFLICT (username) DO UPDATE SET wins = players.wins + 1`

	_, err := db.Pool.Exec(context.Background(), query, username)
	if err != nil {
		logger.Error("player.go >>>> UpdatePlayerWin >>>>> Failed to update win", err)
	}
	return err
}

func GetLeaderboard() ([]PlayerStats, error) {
	const query = `SELECT username, wins FROM players ORDER BY wins DESC LIMIT 10`

	rows, err := db.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []PlayerStats
	for rows.Next() {
		var s PlayerStats
		if err := rows.Scan(&s.Username, &s.Wins); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}
