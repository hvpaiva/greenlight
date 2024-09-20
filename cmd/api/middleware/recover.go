package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hvpaiva/greenlight/cmd/api/erro"
)

func (m *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				if e, ok := err.(error); ok {
					erro.Handle(m.App, w, r, erro.ThrowInternalServer(fmt.Sprintf("panic: %s", e.Error()), e))
					return
				}

				if s, ok := err.(string); ok {
					erro.Handle(m.App, w, r, erro.ThrowInternalServer(fmt.Sprintf("panic: %s", s), errors.New(s)))
					return
				}

				erro.Handle(m.App, w, r, erro.InternalServer)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
