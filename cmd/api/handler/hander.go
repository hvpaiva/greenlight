package handler

import (
	"database/sql"

	"github.com/hvpaiva/greenlight/cmd/api/app"
	"github.com/hvpaiva/greenlight/cmd/api/middleware"
	"github.com/hvpaiva/greenlight/internal/data"
)

type Handler struct {
	App        *app.Application
	Middleware *middleware.Middleware
	Models     *data.Models
}

func New(app *app.Application, db *sql.DB, limiter *middleware.Limiter) *Handler {
	models := data.New(db)
	return &Handler{
		App:        app,
		Middleware: middleware.New(app, models, limiter),
		Models:     models,
	}
}
