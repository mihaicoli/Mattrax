package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Logging logs each request
func Logging() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug().Str("path", r.URL.Path).Str("method", r.Method).Str("remote", r.RemoteAddr).Msg("req")
			next.ServeHTTP(w, r)
		})
	}
}
