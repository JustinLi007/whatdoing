package server

import (
	"net/http"

	"github.com/JustinLi007/whatdoing/libs/go/utils"
	"github.com/go-chi/chi/v5"
)

func (s *Server) RegisterRoutes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Get("/healthz", s.Healthz)

	return mux
}

func (s *Server) Healthz(w http.ResponseWriter, r *http.Request) {
	utils.WriteJson(w, http.StatusOK, utils.Envelope{
		"message": "service progress good",
	})
}
