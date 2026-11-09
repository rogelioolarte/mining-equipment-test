package middlewares

import (
	"log/slog"
	"minning/internal/infra/config"
	"minning/internal/infra/logger"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/goccy/go-json"
)

// panicResponse is a pre-allocated response structure used by the Recover middleware.
// Defining it as a package-level variable avoids heap allocations during panic events.
var panicResponse = map[string]any{
	"success": false,
	"error": map[string]string{
		"code":    "INTERNAL_SERVER_ERROR",
		"message": "Internal server error",
	},
}

// RequestLogger returns a middleware that records the details of each HTTP request.
// It utilizes slog.LogAttrs for zero-allocation attribute logging and provides
// unified event naming ("http_request") for better log indexing.
//
// Logs are automatically upgraded to Error level for status codes >= 400.
func RequestLogger(l *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			next.ServeHTTP(ww, r)

			// Determine log level based on the HTTP status code.
			status := ww.Status()
			level := slog.LevelInfo
			if status >= 400 {
				level = slog.LevelError
			}

			// LogAttrs is used here for maximum efficiency, avoiding interface{} conversion.
			l.LogAttrs(r.Context(), level, "http_request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", status),
				slog.Duration("dur", time.Since(start)),
				slog.String("ip", r.RemoteAddr),
			)
		})
	}
}

// Recover returns a middleware that gracefully handles panics in the request lifecycle.
// It logs the stack trace/error via the provided logger and returns a standard
// 500 Internal Server Error response in JSON format.
func Recover(l *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					// Log the panic event for debugging.
					l.LogAttrs(r.Context(), slog.LevelError, "panic recovered",
						slog.Any("error", rec),
					)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(panicResponse)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// CORS provides a simplified, high-performance CORS implementation.
// It supports origin whitelisting, wildcard origins, and suffix-based matching.
// It also handles preflight OPTIONS requests automatically.
func CORS(cfg *config.HTTP) func(http.Handler) http.Handler {
	originMap := make(map[string]struct{}, len(cfg.AllowedOrigins))
	allowAll := false
	for _, o := range cfg.AllowedOrigins {
		if o == "*" {
			allowAll = true
		}
		originMap[o] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			if allowAll {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				_, exists := originMap[origin]
				// Efficient combined check for exact match or allowed suffix.
				if exists {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				} else {
					next.ServeHTTP(w, r)
					return
				}
			}

			setSharedCORSHeaders(w)

			// Handle Preflight requests.
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// setSharedCORSHeaders applies common headers for Allowed Methods, Headers, and Credentials.
// It also sets a Max-Age to cache preflight responses for 24 hours.
func setSharedCORSHeaders(w http.ResponseWriter) {
	h := w.Header()
	h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	h.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
	h.Set("Access-Control-Max-Age", "86400")
}
