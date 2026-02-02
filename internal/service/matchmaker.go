package service

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/yourusername/connect-four-backend/internal/logger"
)

type Matchmaker struct {
	waitingPlayer *Player
	mu            sync.Mutex
}

var Instance = &Matchmaker{}

func (m *Matchmaker) AddToQueue(player *Player) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.waitingPlayer == nil {
		m.waitingPlayer = player
		logger.Info("matchmaker.go >>>> AddToQueue >>>>> Player " + player.Username + " is waiting for opponent")

		// Start a timer to spawn a bot after 10 seconds
		go m.handleTimeout(player)
	} else {
		// Match found!
		opponent := m.waitingPlayer
		m.waitingPlayer = nil
		logger.Info("matchmaker.go >>>> AddToQueue >>>>> Match found: " + player.Username + " vs " + opponent.Username)

		m.startNewGame(player, opponent)
	}
}

func (m *Matchmaker) handleTimeout(player *Player) {
	time.Sleep(10 * time.Second)

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.waitingPlayer != nil && m.waitingPlayer.Username == player.Username {
		m.waitingPlayer = nil
		logger.Info("matchmaker.go >>>> handleTimeout >>>>> Timeout! Spawning bot for " + player.Username)

		m.startBotGame(player)
	}
}

func (m *Matchmaker) startNewGame(p1, p2 *Player) {
	logger.Info("matchmaker.go >>>> startNewGame >>>>> Game started: " + p1.Username + " vs " + p2.Username)
	room := NewGameRoom(p1, p2, false)

	// Notify players about game start
	msg, _ := json.Marshal(map[string]string{
		"type":     "GAME_START",
		"roomId":   room.ID,
		"opponent": p2.Username,
	})
	room.SendMessage(p1, msg)

	msg2, _ := json.Marshal(map[string]string{
		"type":     "GAME_START",
		"roomId":   room.ID,
		"opponent": p1.Username,
	})
	room.SendMessage(p2, msg2)

	room.BroadcastState()
}

func (m *Matchmaker) startBotGame(p1 *Player) {
	logger.Info("matchmaker.go >>>> startBotGame >>>>> Match vs Bot for " + p1.Username)
	room := NewGameRoom(p1, nil, true)

	msg, _ := json.Marshal(map[string]string{
		"type":     "GAME_START",
		"roomId":   room.ID,
		"opponent": "Competitive Bot",
	})
	room.SendMessage(p1, msg)

	room.BroadcastState()
}
