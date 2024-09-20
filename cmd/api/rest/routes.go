package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *Application) Routes() http.Handler {

	router := httprouter.New()

	router.GET("/v1/healthcheck", a.healthcheckHandler)

	router.GET("/v1/movies", a.Authorize(a.showMoviesHandler))
	router.GET("/v1/movies/:id", a.Authorize(a.getMovieHandler))
	router.POST("/v1/movies", a.Authorize(a.createMovieHandler))
	router.PUT("/v1/movies/:id", a.Authorize(a.updateMovieHandler))
	router.DELETE("/v1/movies/:id", a.Authorize(a.deleteMovieHandler))
	router.PATCH("/v1/movies/:id", a.Authorize(a.patchMovieHandler))

	router.POST("/v1/users", a.registerUserHandler)
	router.PATCH("/v1/users/activated", a.activateUserHandler)

	router.POST("/v1/tokens/authentication", a.creteAuthTokenHandler)

	router.NotFound = http.HandlerFunc(a.NotFoundFunc)
	router.MethodNotAllowed = http.HandlerFunc(a.MethodNotAllowedFunc)

	return a.recoverPanic(a.RateLimit(a.Authenticate(router)))
}
