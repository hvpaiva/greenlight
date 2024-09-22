package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hvpaiva/greenlight/cmd/api/middleware"
)

type config struct {
	port    int
	env     string
	version string
	debug   bool
	db      dbConfig
	limiter middleware.Limiter
	cors    corsConfig
}

type corsConfig struct {
	trustedOrigins []string
}

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

func initConfig(version string) config {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.BoolVar(&cfg.debug, "debug", false, "Enable debug mode")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL dsn")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.Rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	cfg.version = version

	return cfg
}
