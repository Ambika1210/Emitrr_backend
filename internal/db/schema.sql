-- Create players table
CREATE TABLE IF NOT EXISTS players (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    wins INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create games table
CREATE TABLE IF NOT EXISTS games (
    id SERIAL PRIMARY KEY,
    player1_id VARCHAR(50) NOT NULL, -- using username for simplicity in this MVP
    player2_id VARCHAR(50) NOT NULL,
    winner_id VARCHAR(50), -- NULL if draw
    game_state JSONB,
    duration_seconds INTEGER,
    played_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Optional: Index for leaderboard
CREATE INDEX IF NOT EXISTS idx_players_wins ON players(wins DESC);
