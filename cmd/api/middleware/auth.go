package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/hvpaiva/greenlight/internal/data"
	"github.com/hvpaiva/greenlight/pkg/validator"
)

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = m.App.ContextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			m.App.HandleUnauthorized(w, r, "invalid or malformed authorization header")
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateToken(v, token); !v.Valid() {
			m.App.HandleUnauthorized(w, r, "invalid or malformed authorization token")
			return
		}

		user, err := m.Models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				m.App.HandleUnauthorized(w, r, "invalid or expired authorization token")
			default:
				m.App.HandleError(w, r, err)
			}
			return
		}

		r = m.App.ContextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Authorize(permission string, next httprouter.Handle) httprouter.Handle {
	fn := func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		user := m.App.ContextGetUser(r)

		if user.IsAnonymous() {
			m.App.HandleUnauthorized(w, r, "authentication required")
			return
		}

		if !user.Activated {
			m.App.HandleForbidden(w, r, "user not activated")
			return
		}

		next(w, r, param)
	}

	return m.CheckPermissions(permission, fn)
}

func (m *Middleware) CheckPermissions(permission string, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		user := m.App.ContextGetUser(r)

		permissions, err := m.Models.Permission.GetAllForUser(user.ID)
		if err != nil {
			m.App.HandleError(w, r, err)
			return
		}

		if !permissions.Contains(permission) {
			m.App.HandleForbidden(w, r, "user does not have the required permission")
			return
		}

		next(w, r, param)
	}
}
