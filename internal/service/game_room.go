package service

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yourusername/connect-four-backend/internal/game"
	"github.com/yourusername/connect-four-backend/internal/logger"
	"github.com/yourusername/connect-four-backend/internal/repository"
)

type Player struct {
	Username string
	Conn     interface{} // websocket.Conn
	RoomID   string
	IsActive bool
}

type GameRoom struct {
	ID         string
	Player1    *Player
	Player2    *Player // Could be nil if playing vs Bot
	IsVsBot    bool
	GameState  *game.GameState
	Mu         sync.Mutex
	LastActive time.Time
}

type GameMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

var Rooms = make(map[string]*GameRoom)
var Players = make(map[string]*Player) // For reconnection
var GlobalMu sync.Mutex

func NewGameRoom(p1, p2 *Player, isBot bool) *GameRoom {
	room := &GameRoom{
		ID:         p1.Username + "_" + time.Now().Format("150405"),
		Player1:    p1,
		Player2:    p2,
		IsVsBot:    isBot,
		GameState:  game.NewConnectFourGame(),
		LastActive: time.Now(),
	}

	p1.RoomID = room.ID
	if p2 != nil {
		p2.RoomID = room.ID
	}

	GlobalMu.Lock()
	Rooms[room.ID] = room
	GlobalMu.Unlock()

	return room
}

func (r *GameRoom) BroadcastState() {
	stateBytes, _ := json.Marshal(map[string]interface{}{
		"type":    "GAME_STATE",
		"payload": r.GameState,
	})

	r.SendMessage(r.Player1, stateBytes)
	if !r.IsVsBot && r.Player2 != nil {
		r.SendMessage(r.Player2, stateBytes)
	}
}

func (r *GameRoom) SendMessage(p *Player, data []byte) {
	if p == nil || p.Conn == nil || !p.IsActive {
		return
	}
	conn := p.Conn.(*websocket.Conn)
	err := conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		logger.Error("game_room.go >>>> SendMessage >>>>> Failed to send to "+p.Username, err)
	}
}

func (r *GameRoom) HandleMove(username string, col int) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	if r.GameState.IsFinished {
		return
	}

	// 1. Validation
	isP1 := username == r.Player1.Username
	curr := r.GameState.CurrentPlayer

	if (isP1 && curr != game.Player1) || (!isP1 && curr != game.Player2) {
		logger.Warn("game_room.go >>>> HandleMove >>>>> Not player's turn: " + username)
		return
	}

	// 2. Execute Move
	err := r.GameState.ExecutePlayerMove(col)
	if err != nil {
		logger.Warn("game_room.go >>>> HandleMove >>>>> Move failed for " + username + " | " + err.Error())
		return
	}

	r.LastActive = time.Now()
	r.BroadcastState()

	// 3. Finish Logic
	if r.GameState.IsFinished {
		r.finishGame()
	} else if r.IsVsBot && r.GameState.CurrentPlayer == game.Player2 {
		go func() {
			time.Sleep(1 * time.Second)
			botCol := r.GameState.SelectBestMove()
			r.HandleMove("Bot", botCol)
		}()
	}
}

func (r *GameRoom) finishGame() {
	logger.Info("game_room.go >>>> finishGame >>>>> Game finished ID: " + r.ID)

	winnerName := ""
	if r.GameState.Winner == game.Player1 {
		winnerName = r.Player1.Username
		repository.UpdatePlayerWin(winnerName)
	} else if r.GameState.Winner == game.Player2 {
		if r.IsVsBot {
			winnerName = "Bot"
		} else {
			winnerName = r.Player2.Username
			repository.UpdatePlayerWin(winnerName)
		}
	}

	p2Name := "Bot"
	if r.Player2 != nil {
		p2Name = r.Player2.Username
	}

	repository.SaveGameResult(repository.GameRecord{
		Player1: r.Player1.Username,
		Player2: p2Name,
		Winner:  winnerName,
		State:   r.GameState.Board,
	})
}
