package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/yourusername/connect-four-backend/internal/api/routes"
	"github.com/yourusername/connect-four-backend/internal/db"
	"github.com/yourusername/connect-four-backend/internal/logger"
)

func main() {
	logger.Init()
	logger.Info("main.go >>>> main >>>>> Initializing application")

	if err := godotenv.Load(); err != nil {
		logger.Warn("main.go >>>> main >>>>> No .env file found")
	}

	if err := db.Connect(); err != nil {
		logger.Error("main.go >>>> main >>>>> Database connection failed", err)
	}
	defer db.Close()

	handler := routes.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("main.go >>>> main >>>>> Server starting on port " + port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		logger.Error("main.go >>>> main >>>>> Server failed to start", err)
	}
}
