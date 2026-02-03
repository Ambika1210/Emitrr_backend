package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/connect-four-backend/internal/logger"
)

var Pool *pgxpool.Pool

const fileName = "db.go"

func Connect() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		logger.Error("db.go >>>> Connect >>>>> DATABASE_URL environment variable is not set", nil)
		return os.ErrNotExist
	}

	var err error
	Pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Error("db.go >>>> Connect >>>>> Creating pool failed", err)
		return err
	}

	err = Pool.Ping(context.Background())
	if err != nil {
		logger.Error("db.go >>>> Connect >>>>> Ping failed", err)
		return err
	}

	logger.Info("db.go >>>> Connect >>>>> Successfully connected to Postgres")
	return nil
}

func Close() {
	if Pool != nil {
		Pool.Close()
		logger.Info("db.go >>>> Close >>>>> Database connection closed")
	}
}
