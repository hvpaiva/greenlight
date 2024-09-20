package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/hvpaiva/greenlight/cmd/api/app"
	"github.com/hvpaiva/greenlight/cmd/api/handler"
)

const version = "1.0.0"

func main() {
	cfg := initConfig(version)
	logger := app.NewLogger(cfg.debug)

	db, err := openDB(cfg.db)
	if err != nil {
		logger.Error("database failed to open", slog.String("erro", err.Error()))
		os.Exit(1)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error("database failed to close", slog.String("erro", err.Error()))
			os.Exit(1)
		}
	}(db)

	a := app.New(logger, cfg.env, cfg.version, cfg.cors.trustedOrigins)
	h := handler.New(a, db, &cfg.limiter)

	if err := serve(cfg, a, h); err != nil {
		logger.Error("server failed to start", slog.String("erro", err.Error()))
		os.Exit(1)
	}
}

func openDB(dbConfig dbConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbConfig.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(dbConfig.maxOpenConns)
	db.SetMaxIdleConns(dbConfig.maxIdleConns)
	db.SetConnMaxIdleTime(dbConfig.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
