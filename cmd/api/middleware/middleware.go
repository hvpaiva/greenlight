package middleware

import (
	"github.com/hvpaiva/greenlight/cmd/api/app"
	"github.com/hvpaiva/greenlight/internal/data"
)

type Middleware struct {
	App     *app.Application
	Models  *data.Models
	Limiter *Limiter
}

func New(app *app.Application, models *data.Models, limiter *Limiter) *Middleware {
	return &Middleware{
		App:     app,
		Models:  models,
		Limiter: limiter,
	}
}
