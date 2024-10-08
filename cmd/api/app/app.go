package app

import (
	"log/slog"
	"os"
	"sync"
)

type Application struct {
	Logger      *slog.Logger
	Wg          sync.WaitGroup
	Env         string
	Version     string
	TrustedCors []string
}

func New(logger *slog.Logger, env, version string, trustedCors []string) *Application {
	return &Application{
		Logger:      logger,
		Env:         env,
		Version:     version,
		TrustedCors: trustedCors,
	}
}

func NewLogger(debug bool) *slog.Logger {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: debug,
		Level:     level,
	}))
}
