package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

// FrontendHeaders sets the headers for the frontend UI routes
func FrontendHeaders() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Security-Policy", "default-src 'self' 'unsafe-inline' cdn.jsdelivr.net cdnjs.cloudflare.com") // TODO: Remove unsafe-inline and cdns, use trusted types
			w.Header().Add("X-XSS-Protection", "1; mode=block")
			w.Header().Add("X-Frame-Options", "DENY")
			w.Header().Add("X-Content-Type-Options", "nosniff")

			// w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload") // TODO: If configured by user
			// TODO: Also Expect-CT with user config

			next.ServeHTTP(w, r)
		})
	}
}

// Headers sets the global headers for the server
func Headers() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Server", "MattraxMDM")
			next.ServeHTTP(w, r)
		})
	}
}
