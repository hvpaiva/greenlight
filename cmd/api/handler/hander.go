package handler

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/hvpaiva/greenlight/cmd/api/app"
	"github.com/hvpaiva/greenlight/cmd/api/erro"
	"github.com/hvpaiva/greenlight/cmd/api/middleware"
	"github.com/hvpaiva/greenlight/internal/data"
)

type Func func(http.ResponseWriter, *http.Request, httprouter.Params) error

type MFunc func(http.Handler) (http.Handler, error)

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

func (h *Handler) adapt(handler Func) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if e := handler(w, r, p); e != nil {
			erro.Handle(h.App, w, r, e)
		}
	}
}
