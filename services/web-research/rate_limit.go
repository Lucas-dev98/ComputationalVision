package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type rateLimitConfig struct {
	enabled bool
	limit   int
	window  time.Duration
}

type rateLimitEntry struct {
	count       int
	windowStart time.Time
}

type rateLimiter struct {
	config  rateLimitConfig
	mu      sync.Mutex
	clients map[string]*rateLimitEntry
}

func loadRateLimitConfig() rateLimitConfig {
	return rateLimitConfig{
		enabled: getEnv("RATE_LIMIT_ENABLED", "true") != "false",
		limit:   getEnvAsInt("RATE_LIMIT_REQUESTS", 120),
		window:  time.Duration(getEnvAsInt("RATE_LIMIT_WINDOW_SECONDS", 60)) * time.Second,
	}
}

func newRateLimitMiddleware(config rateLimitConfig) mux.MiddlewareFunc {
	limiter := &rateLimiter{
		config:  config,
		clients: make(map[string]*rateLimitEntry),
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.config.enabled || r.Method == http.MethodOptions || r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			allowed, remaining, retryAfter := limiter.allow(clientIP(r), time.Now())
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limiter.config.limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
				w.WriteHeader(http.StatusTooManyRequests)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": "Rate limit excedido. Tente novamente em alguns segundos.",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (l *rateLimiter) allow(client string, now time.Time) (bool, int, int) {
	if client == "" {
		client = "unknown"
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	entry, ok := l.clients[client]
	if !ok || now.Sub(entry.windowStart) >= l.config.window {
		l.clients[client] = &rateLimitEntry{count: 1, windowStart: now}
		return true, max(l.config.limit-1, 0), 0
	}

	if entry.count >= l.config.limit {
		retryAfter := int(l.config.window.Seconds() - now.Sub(entry.windowStart).Seconds())
		if retryAfter < 1 {
			retryAfter = 1
		}
		return false, 0, retryAfter
	}

	entry.count++
	return true, max(l.config.limit-entry.count, 0), 0
}

func clientIP(r *http.Request) string {
	forwarded := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if forwarded != "" {
		return forwarded
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}

func getEnvAsInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}

	return value
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
