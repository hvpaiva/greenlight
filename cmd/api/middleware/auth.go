package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/hvpaiva/greenlight/cmd/api/erro"
	"github.com/hvpaiva/greenlight/internal/data"
	"github.com/hvpaiva/greenlight/pkg/validator"
)

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		w.Header().Set("WWW-Authenticate", "Bearer")

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = m.App.ContextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			erro.Handle(m.App, w, r, erro.Unauthorized.WithMessage("invalid or malformed authorization header"))
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateToken(v, token); !v.Valid() {
			erro.Handle(m.App, w, r, erro.Unauthorized.WithMessage("invalid or malformed authorization token"))
			return
		}

		user, err := m.Models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				erro.Handle(m.App, w, r, erro.Unauthorized.WithMessage("invalid or expired authorization token"))
			default:
				erro.Handle(m.App, w, r, erro.ThrowInternalServer("get user token", err))
			}
			return
		}

		r = m.App.ContextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Authorize(permission string) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := m.App.ContextGetUser(r)

			if user.IsAnonymous() {
				erro.Handle(m.App, w, r, erro.Unauthorized.WithMessage("authentication required"))
				return
			}

			if !user.Activated {
				erro.Handle(m.App, w, r, erro.Forbidden.WithMessage("user not activated"))
				return
			}

			handler.ServeHTTP(w, r)
		})

		return m.CheckPermissions(permission)(fn)
	}
}

func (m *Middleware) CheckPermissions(permission string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := m.App.ContextGetUser(r)

			permissions, err := m.Models.Permission.GetAllForUser(user.ID)
			if err != nil {
				erro.Handle(m.App, w, r, erro.ThrowInternalServer("get user permissions", err))
				return
			}

			if !permissions.Contains(permission) {
				erro.Handle(m.App, w, r, erro.Forbidden.WithMessage("user does not have the required permission"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
