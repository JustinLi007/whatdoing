package middleware

import (
	"gateway/internal/service"
	"gateway/internal/utils"
	"gateway/internal/verifier"
	"net/http"
)

type Middleware interface {
	Cors(next http.Handler) http.Handler
	VerifyJwt(next http.Handler) http.Handler
}

type middleware struct {
	verifier   verifier.Verifier
	serviceMap service.ServiceMap
}

var middlewareInstance *middleware

var allowedOrigins = map[string]bool{
	"http://localhost:5173": true,
}

func NewMiddleware(verifier verifier.Verifier, serviceMap service.ServiceMap) Middleware {
	if middlewareInstance != nil {
		return middlewareInstance
	}
	newMiddleware := &middleware{
		verifier:   verifier,
		serviceMap: serviceMap,
	}
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

func (m *middleware) VerifyJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prefix, _, ok := utils.ParseRequestUrl(r)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		endpoint, err := m.serviceMap.GetEndpoint(prefix)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if endpoint.Public {
			next.ServeHTTP(w, r)
			return
		}

		jwtCookie, err := r.Cookie("jwt")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sub, scope, err := m.verifier.ValidateJwt(jwtCookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !utils.HasScope(endpoint.Scope, scope) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		r.Header.Del("Whatdoing-User-Id")
		r.Header.Del("Whatdoing-Scope")

		r.Header.Set("Whatdoing-User-Id", sub)
		r.Header.Set("Whatdoing-Scope", scope)

		next.ServeHTTP(w, r)
	})
}
