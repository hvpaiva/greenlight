package rest

import (
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/hvpaiva/greenlight/internal/data"
)

const version = "1.0.0"

type Application struct {
	Config Config
	Logger *slog.Logger
	Models data.Models
}

type Config struct {
	Port    int
	Env     string
	Version string
	Debug   bool
	DB      DB
	Limiter Limiter
}

type DB struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type Limiter struct {
	RPS     float64
	Burst   int
	Enabled bool
}

func NewApplication(db *sql.DB, cfg Config, logger *slog.Logger) *Application {
	return &Application{
		Config: cfg,
		Logger: logger,
		Models: data.NewModel(db),
	}
}

func InitConfig() Config {
	var cfg Config

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.BoolVar(&cfg.Debug, "debug", false, "Enable debug mode")

	flag.StringVar(&cfg.DB.DSN, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.Limiter.RPS, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Parse()

	cfg.Version = version

	return cfg
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
