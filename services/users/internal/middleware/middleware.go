package middleware

import (
	"net/http"
)

type Middleware interface {
	Cors(next http.Handler) http.Handler
}

type middleware struct {
}

var middlewareInstance *middleware

var allowedOrigins = map[string]bool{
	"http://localhost:5173": true,
}

func NewMiddleware() Middleware {
	if middlewareInstance != nil {
		return middlewareInstance
	}
	newMiddleware := &middleware{}
	middlewareInstance = newMiddleware
	return middlewareInstance
}

func (m *middleware) Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
