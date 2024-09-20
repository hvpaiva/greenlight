package rest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/hvpaiva/greenlight/internal/data"
	"github.com/hvpaiva/greenlight/pkg/validator"
)

func (a *Application) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = a.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			a.HandleUnauthorized(w, r, "invalid or malformed authorization header")
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateToken(v, token); !v.Valid() {
			a.HandleUnauthorized(w, r, "invalid or malformed authorization token")
			return
		}

		user, err := a.Models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				a.HandleUnauthorized(w, r, "invalid or expired authorization token")
			default:
				a.HandleError(w, r, err)
			}
			return
		}

		r = a.contextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}

func (a *Application) Authorize(permission string, next httprouter.Handle) httprouter.Handle {
	fn := func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		user := a.contextGetUser(r)

		if user.IsAnonymous() {
			a.HandleUnauthorized(w, r, "authentication required")
			return
		}

		if !user.Activated {
			a.HandleForbidden(w, r, "user not activated")
			return
		}

		next(w, r, param)
	}

	return a.CheckPermissions(permission, fn)
}

func (a *Application) CheckPermissions(permission string, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		user := a.contextGetUser(r)

		permissions, err := a.Models.Permission.GetAllForUser(user.ID)
		if err != nil {
			a.HandleError(w, r, err)
			return
		}

		if !permissions.Contains(permission) {
			a.HandleForbidden(w, r, "user does not have the required permission")
			return
		}

		next(w, r, param)
	}
}
