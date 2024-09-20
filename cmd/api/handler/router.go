package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/hvpaiva/greenlight/internal/data"
)

func (h *Handler) Router() http.Handler {
	router := httprouter.New()

	router.GET("/v1/healthcheck", h.healthcheckHandler)

	router.GET("/v1/movies", h.Middleware.Authorize(data.PermissionMovieRead, h.showMoviesHandler))
	router.GET("/v1/movies/:id", h.Middleware.Authorize(data.PermissionMovieRead, h.getMovieHandler))
	router.POST("/v1/movies", h.Middleware.Authorize(data.PermissionMovieWrite, h.createMovieHandler))
	router.PUT("/v1/movies/:id", h.Middleware.Authorize(data.PermissionMovieWrite, h.updateMovieHandler))
	router.DELETE("/v1/movies/:id", h.Middleware.Authorize(data.PermissionMovieWrite, h.deleteMovieHandler))
	router.PATCH("/v1/movies/:id", h.Middleware.Authorize(data.PermissionMovieWrite, h.patchMovieHandler))

	router.POST("/v1/users", h.registerUserHandler)
	router.PATCH("/v1/users/activated", h.activateUserHandler)

	router.POST("/v1/tokens/authentication", h.CreteAuthTokenHandler)

	router.NotFound = http.HandlerFunc(h.App.NotFoundFunc)
	router.MethodNotAllowed = http.HandlerFunc(h.App.MethodNotAllowedFunc)

	return h.Middleware.RecoverPanic(h.Middleware.RateLimit(h.Middleware.Authenticate(router)))
}
