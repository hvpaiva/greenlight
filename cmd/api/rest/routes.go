package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *Application) Routes() http.Handler {

	router := httprouter.New()

	router.GET("/v1/healthcheck", a.healthcheckHandler)

	router.GET("/v1/movies", a.showMoviesHandler)
	router.GET("/v1/movies/:id", a.getMovieHandler)
	router.POST("/v1/movies", a.createMovieHandler)
	router.PUT("/v1/movies/:id", a.updateMovieHandler)
	router.DELETE("/v1/movies/:id", a.deleteMovieHandler)
	router.PATCH("/v1/movies/:id", a.patchMovieHandler)

	router.POST("/v1/users", a.registerUserHandler)

	router.NotFound = http.HandlerFunc(a.NotFoundFunc)
	router.MethodNotAllowed = http.HandlerFunc(a.MethodNotAllowedFunc)

	return a.recoverPanic(a.RateLimit(router))
}
