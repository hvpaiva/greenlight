package rest

import (
	"net/http"

	"golang.org/x/time/rate"
)

func (a Application) RateLimit(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			a.HandleError(w, r, NewHTTPError("rate limit exceeded", http.StatusTooManyRequests))
			return
		}

		next.ServeHTTP(w, r)
	})
}
