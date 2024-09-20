package handler

import (
	"database/sql"
	"net/http"

	"github.com/hvpaiva/greenlight/cmd/api/app"
	"github.com/hvpaiva/greenlight/cmd/api/erro"
	"github.com/hvpaiva/greenlight/cmd/api/middleware"
	"github.com/hvpaiva/greenlight/internal/data"
)

type handlerFunc func(http.ResponseWriter, *http.Request) error

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

func (h *Handler) adapt(handler handlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if e := handler(w, r); e != nil {
			erro.Handle(h.App, w, r, e)
		}
	})
}
