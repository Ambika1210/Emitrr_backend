package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/connect-four-backend/internal/logger"
)

const fileName = "health.go"

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("health.go >>>> HealthHandler >>>>> Processing health check request")

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "up",
		"message": "Connect Four Backend is running with Clean Architecture & Color Logger",
	})
}
