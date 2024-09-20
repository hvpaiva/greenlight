package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/hvpaiva/greenlight/cmd/api/erro"
	"github.com/hvpaiva/greenlight/cmd/api/middleware"
	"github.com/hvpaiva/greenlight/internal/data"
)

func (h *Handler) Router() http.Handler {
	r := httprouter.New()

	h.register(r, http.MethodGet, "/v1/healthcheck", h.healthcheckHandler)

	h.register(r, http.MethodGet, "/v1/movies", h.showMoviesHandler, h.Middleware.Authorize(data.PermissionMovieRead))
	h.register(r, http.MethodGet, "/v1/movies/:id", h.getMovieHandler, h.Middleware.Authorize(data.PermissionMovieRead))
	h.register(r, http.MethodPost, "/v1/movies", h.createMovieHandler, h.Middleware.Authorize(data.PermissionMovieWrite))
	h.register(r, http.MethodPut, "/v1/movies/:id", h.updateMovieHandler, h.Middleware.Authorize(data.PermissionMovieWrite))
	h.register(r, http.MethodDelete, "/v1/movies/:id", h.deleteMovieHandler, h.Middleware.Authorize(data.PermissionMovieWrite))
	h.register(r, http.MethodPatch, "/v1/movies/:id", h.patchMovieHandler, h.Middleware.Authorize(data.PermissionMovieWrite))

	h.register(r, http.MethodPost, "/v1/users", h.registerUserHandler)
	h.register(r, http.MethodPatch, "/v1/users/activated", h.activateUserHandler)

	h.register(r, http.MethodPost, "/v1/tokens/authentication", h.creteAuthTokenHandler)

	r.NotFound = notFoundFunc(h)
	r.MethodNotAllowed = methodNotAllowedFunc(h)

	return addMiddlewares(r, h.Middleware.RecoverPanic, h.Middleware.RateLimit, h.Middleware.Authenticate)
}

func (h *Handler) register(r *httprouter.Router, method string, path string, handle handlerFunc, middlewares ...middleware.Func) {
	var aggregated = h.adapt(handle)

	for _, m := range middlewares {
		aggregated = m(aggregated)
	}

	r.Handler(method, path, aggregated)
}

func addMiddlewares(r *httprouter.Router, middlewares ...middleware.Func) http.Handler {
	var aggregated http.Handler = r
	for _, m := range middlewares {
		aggregated = m(aggregated)
	}

	return aggregated
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
