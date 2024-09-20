package rest

import (
	"context"
	"net/http"

	"github.com/hvpaiva/greenlight/internal/data"
)

type contextKey string

const userContextKey = contextKey("user")

func (a *Application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (a *Application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("user not found in request context")
	}

	return user
}
