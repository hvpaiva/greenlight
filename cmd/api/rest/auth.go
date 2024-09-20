package rest

import (
	"errors"
	"net/http"
	"strings"

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
