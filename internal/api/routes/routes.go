package routes

import (
	"net/http"

	"github.com/yourusername/connect-four-backend/internal/api/handlers"
	"github.com/yourusername/connect-four-backend/internal/api/middleware"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Health Route
	mux.HandleFunc("/health", handlers.HealthHandler)

	// Leaderboard Route
	mux.HandleFunc("/leaderboard", handlers.LeaderboardHandler)

	// WebSocket Route (Matchmaking & Gameplay)
	mux.HandleFunc("/ws", handlers.WebSocketHandler)

	// Wrap mux with middlewares (Node.js style)
	handler := middleware.JSONMiddleware(mux)
	return middleware.CORSMiddleware(handler)
}
