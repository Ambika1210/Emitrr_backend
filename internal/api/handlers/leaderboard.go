package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/connect-four-backend/internal/logger"
	"github.com/yourusername/connect-four-backend/internal/repository"
)

func LeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("leaderboard.go >>>> LeaderboardHandler >>>>> Fetching top players")

	stats, err := repository.GetLeaderboard()
	if err != nil {
		logger.Error("leaderboard.go >>>> LeaderboardHandler >>>>> DB error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}
