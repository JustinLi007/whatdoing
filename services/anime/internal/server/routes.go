package server

import (
	"net/http"

	"github.com/JustinLi007/whatdoing/libs/go/utils"

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
	utils.WriteJson(w, http.StatusOK, utils.Envelope{
		"message": "service anime good",
	})
}
