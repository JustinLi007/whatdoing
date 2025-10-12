package server

import (
	"net/http"

	"github.com/JustinLi007/whatdoing/libs/go/util"
	"github.com/go-chi/chi/v5"
)

func (s *Server) RegisterRoutes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Get("/healthz", s.Healthz)

	return mux
}

func (s *Server) Healthz(w http.ResponseWriter, r *http.Request) {
	util.WriteJson(w, http.StatusOK, util.Envelope{
		"message": "service progress good",
	})
}
