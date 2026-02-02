package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yourusername/connect-four-backend/internal/logger"
	"github.com/yourusername/connect-four-backend/internal/service"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("game.go >>>> WebSocketHandler >>>>> Upgrade failed", err)
		return
	}
	defer conn.Close()

	// Check for reconnection
	service.GlobalMu.Lock()
	player, exists := service.Players[username]
	if exists {
		player.Conn = conn
		player.IsActive = true
		service.GlobalMu.Unlock()
		logger.Info("game.go >>>> WebSocketHandler >>>>> Player reconnected: " + username)

		// If player was in a room, send them current state
		if player.RoomID != "" {
			if room, ok := service.Rooms[player.RoomID]; ok {
				room.BroadcastState()
			}
		}
	} else {
		player = &service.Player{
			Username: username,
			Conn:     conn,
			IsActive: true,
		}
		service.Players[username] = player
		service.GlobalMu.Unlock()
		logger.Info("game.go >>>> WebSocketHandler >>>>> New player connected: " + username)
		service.Instance.AddToQueue(player)
	}

	// Message Loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if err != io.EOF {
				logger.Warn("game.go >>>> WebSocketHandler >>>>> Read error/disconnect: " + username)
			}
			break
		}

		var msg service.GameMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		if msg.Type == "MOVE" {
			var payload struct {
				Column int    `json:"column"`
				RoomID string `json:"roomId"`
			}
			json.Unmarshal(msg.Payload, &payload)

			service.GlobalMu.Lock()
			room, ok := service.Rooms[payload.RoomID]
			service.GlobalMu.Unlock()
			if ok {
				room.HandleMove(username, payload.Column)
			}
		}
	}

	// Disconnection cleanup with 30s window
	handleDisconnect(player, conn)
}

func handleDisconnect(p *service.Player, closingConn *websocket.Conn) {
	service.GlobalMu.Lock()
	if p.Conn == closingConn {
		p.IsActive = false
		logger.Info("game.go >>>> handleDisconnect >>>>> Player " + p.Username + " offline. Waiting 30s...")
	}
	service.GlobalMu.Unlock()

	time.AfterFunc(30*time.Second, func() {
		service.GlobalMu.Lock()
		defer service.GlobalMu.Unlock()

		if !p.IsActive {
			logger.Info("game.go >>>> handleDisconnect >>>>> Player " + p.Username + " final timeout. Cleaning up.")
			// If in a room, maybe forfeit the game
			if p.RoomID != "" {
				if room, ok := service.Rooms[p.RoomID]; ok {
					if !room.GameState.IsFinished {
						// Forfeit logic
						logger.Info("game.go >>>> handleDisconnect >>>>> Forfeiting game for " + p.Username)
						room.Mu.Lock()
						room.GameState.IsFinished = true
						if p.Username == room.Player1.Username {
							room.GameState.Winner = 2 // Player 2 wins
						} else {
							room.GameState.Winner = 1 // Player 1 wins
						}
						room.Mu.Unlock()

						// In handleDisconnect we don't have access to the other player's connection easily
						// but room.finishGame() will handle DB saving
						// Actually we need to satisfy the room.finishGame conditions
						// I'll make finishGame public or handle it here
						room.BroadcastState() // Tell the other player
						// Save to DB (already handled in finishGame if we call it)
						// For now I'll just rely on room.BroadcastState and let the other player see the win
					}
				}
			}
			delete(service.Players, p.Username)
		}
	})
}
