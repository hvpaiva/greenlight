package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hvpaiva/greenlight/cmd/api/app"
	"github.com/hvpaiva/greenlight/cmd/api/handler"
)

func serve(c config, a *app.Application, handler *handler.Handler) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.port),
		Handler:      handler.Router(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(a.Logger.Handler(), slog.LevelError),
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		a.Logger.Info("gracefully shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		a.Wg.Wait()
		shutdownError <- nil
	}()

	a.Logger.Info("starting server",
		slog.String("version", c.version),
		slog.String("env", c.env),
		slog.String("addr", srv.Addr),
		slog.Bool("debug", c.debug),
	)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	return nil
}
