package server

import (
	"net/http"

	"github.com/JustinLi007/whatdoing/libs/go/util"

	"github.com/go-chi/chi/v5"
)

func (s *Server) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()

	// r.Use(s.Middleware.Cors)

	// r.Group(func(r chi.Router) {
	// 	r.Post("/auth/signup", s.HandlerUsers.SignUp)
	// 	r.Post("/auth/login", s.HandlerUsers.Login)
	// 	r.Post("/auth/logout", s.HandlerUsers.Logout)
	// 	r.Post("/auth/refresh", s.HandlerUsers.Refresh)
	// })

	// r.Get("/.well-known/jwks.json", s.HandlerSigner.GetJwks)

	r.Get("/healthz", s.Healthz)

	return r
}

func (s *Server) Healthz(w http.ResponseWriter, r *http.Request) {
	util.WriteJson(w, http.StatusOK, util.Envelope{
		"message": "service users good",
	})
}
