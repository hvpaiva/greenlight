package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/hvpaiva/greenlight/cmd/api/app"
)

type Limiter struct {
	Rps     float64
	Burst   int
	Enabled bool
}

func (m *Middleware) RateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()

		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.Limiter.Enabled {
			return
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			m.App.HandleError(w, r, err)
			return
		}

		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(m.Limiter.Rps), m.Limiter.Burst)}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			m.App.HandleError(w, r, app.NewError("rate limit exceeded", http.StatusTooManyRequests))
			return
		}

		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
