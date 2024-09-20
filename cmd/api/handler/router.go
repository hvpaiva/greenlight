package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/hvpaiva/greenlight/cmd/api/erro"
	"github.com/hvpaiva/greenlight/internal/data"
)

func (h *Handler) Router() http.Handler {
	r := httprouter.New()

	r.GET("/v1/healthcheck", h.adapt(h.healthcheckHandler))

	r.GET("/v1/movies", h.Middleware.Authorize(data.PermissionMovieRead)(h.adapt(h.showMoviesHandler)))
	r.GET("/v1/movies/:id", h.Middleware.Authorize(data.PermissionMovieRead)(h.adapt(h.getMovieHandler)))
	r.POST("/v1/movies", h.Middleware.Authorize(data.PermissionMovieWrite)(h.adapt(h.createMovieHandler)))
	r.PUT("/v1/movies/:id", h.Middleware.Authorize(data.PermissionMovieWrite)(h.adapt(h.updateMovieHandler)))
	r.DELETE("/v1/movies/:id", h.Middleware.Authorize(data.PermissionMovieWrite)(h.adapt(h.deleteMovieHandler)))
	r.PATCH("/v1/movies/:id", h.Middleware.Authorize(data.PermissionMovieWrite)(h.adapt(h.patchMovieHandler)))

	r.POST("/v1/users", h.adapt(h.registerUserHandler))
	r.PATCH("/v1/users/activated", h.adapt(h.activateUserHandler))

	r.POST("/v1/tokens/authentication", h.adapt(h.CreteAuthTokenHandler))

	r.NotFound = notFoundFunc(h)
	r.MethodNotAllowed = methodNotAllowedFunc(h)

	return h.Middleware.RecoverPanic(h.Middleware.RateLimit(h.Middleware.Authenticate(r)))
}

func (h *Handler) register(r *httprouter.Router, method string, path string, handle Func, middlewares ...func(next httprouter.Handle) httprouter.Handle) {
	aggregated := h.adapt(handle)

	for _, m := range middlewares {
		aggregated = m(aggregated)
	}

	r.Handle(method, path, aggregated)
}

var methodNotAllowedFunc = func(h *Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		erro.Handle(h.App, w, r, erro.MethodNotAllowed)
	})
}

var notFoundFunc = func(h *Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		erro.Handle(h.App, w, r, erro.NotFound)
	})
}
