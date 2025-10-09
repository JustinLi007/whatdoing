package server

import (
	"gateway/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(s.Middleware.Cors)

	r.Get("/healthz", s.Healthz)

	r.Group(func(r chi.Router) {
	})

	return r
}

func (s *Server) Healthz(w http.ResponseWriter, r *http.Request) {
	utils.WriteJson(w, http.StatusOK, utils.Envelope{
		"message": "gateway good",
	})
}

func (s *Server) RegisterServices() http.Handler {
	rp := s.NewReverseProxy()
	s.ServiceMap.AddEndpoint("http://auth-service", "auth", "", true)
	s.ServiceMap.AddEndpoint("http://anime-service", "anime", "", false)
	s.ServiceMap.AddEndpoint("", "test", "ohfk", false)
	handler := s.Middleware.Cors(s.Middleware.VerifyJwt(rp))
	return handler
}
