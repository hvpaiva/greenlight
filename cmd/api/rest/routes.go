package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/hvpaiva/greenlight/internal/data"
)

func (a *Application) Routes() http.Handler {

	router := httprouter.New()

	router.GET("/v1/healthcheck", a.healthcheckHandler)

	router.GET("/v1/movies", a.Authorize(data.PermissionMovieRead, a.showMoviesHandler))
	router.GET("/v1/movies/:id", a.Authorize(data.PermissionMovieRead, a.getMovieHandler))
	router.POST("/v1/movies", a.Authorize(data.PermissionMovieWrite, a.createMovieHandler))
	router.PUT("/v1/movies/:id", a.Authorize(data.PermissionMovieWrite, a.updateMovieHandler))
	router.DELETE("/v1/movies/:id", a.Authorize(data.PermissionMovieWrite, a.deleteMovieHandler))
	router.PATCH("/v1/movies/:id", a.Authorize(data.PermissionMovieWrite, a.patchMovieHandler))

	router.POST("/v1/users", a.registerUserHandler)
	router.PATCH("/v1/users/activated", a.activateUserHandler)

	router.POST("/v1/tokens/authentication", a.creteAuthTokenHandler)

	router.NotFound = http.HandlerFunc(a.NotFoundFunc)
	router.MethodNotAllowed = http.HandlerFunc(a.MethodNotAllowedFunc)

	return a.recoverPanic(a.RateLimit(a.Authenticate(router)))
}
