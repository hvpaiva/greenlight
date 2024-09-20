package middleware

import (
	"net/http"

	"github.com/hvpaiva/greenlight/cmd/api/app"
	"github.com/hvpaiva/greenlight/internal/data"
)

type Middleware struct {
	App     *app.Application
	Models  *data.Models
	Limiter *Limiter
}

type Func func(next http.Handler) http.Handler

func New(app *app.Application, models *data.Models, limiter *Limiter) *Middleware {
	return &Middleware{
		App:     app,
		Models:  models,
		Limiter: limiter,
	}
}
