package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/hvpaiva/greenlight/cmd/api/rest"
)

func main() {
	cfg := rest.InitConfig()
	logger := rest.NewLogger(cfg.Debug)

	db, err := openDB(cfg.DB)
	if err != nil {
		logger.Error("database failed to open", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error("database failed to close", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}(db)

	app := rest.NewApplication(db, cfg, logger)

	if err := serve(app); err != nil {
		logger.Error("server failed to start", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func openDB(dbConfig rest.DB) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbConfig.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetConnMaxIdleTime(dbConfig.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
