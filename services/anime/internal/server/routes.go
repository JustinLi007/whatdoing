package server

import (
	"net/http"

	"github.com/JustinLi007/whatdoing/libs/go/util"

	"github.com/go-chi/chi/v5"
)

func (s *Server) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("/anime", s.animeHandler.CreateAnime)
		r.Get("/anime", s.animeHandler.GetAnime)
		r.Put("/anime", s.animeHandler.UpdateAnime)
		r.Delete("/anime", s.animeHandler.DeleteAnime)
	})

	r.Get("/healthz", s.Healthz)

	return r
}

func (s *Server) Healthz(w http.ResponseWriter, r *http.Request) {
	util.WriteJson(w, http.StatusOK, util.Envelope{
		"message": "service anime good",
	})
}
