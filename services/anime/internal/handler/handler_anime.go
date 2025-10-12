package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/JustinLi007/whatdoing/libs/go/util"
	"github.com/JustinLi007/whatdoing/services/anime/internal/database"

	"github.com/google/uuid"
)

type HandlerAnime interface {
	CreateAnime(w http.ResponseWriter, r *http.Request)
	GetAnime(w http.ResponseWriter, r *http.Request)
	UpdateAnime(w http.ResponseWriter, r *http.Request)
	DeleteAnime(w http.ResponseWriter, r *http.Request)
}

type handlerAnime struct {
	animeService database.ServiceAnime
}

var handlerAnimeInstance *handlerAnime

func NewHandlerAnime(animeService database.ServiceAnime) HandlerAnime {
	if handlerAnimeInstance != nil {
		return handlerAnimeInstance
	}

	newHandlerAnime := &handlerAnime{
		animeService: animeService,
	}
	handlerAnimeInstance = newHandlerAnime

	return handlerAnimeInstance
}

func (h *handlerAnime) CreateAnime(w http.ResponseWriter, r *http.Request) {
	type CreateAnimeRequest struct {
		Name     string `json:"name"`
		Episodes int    `json:"episodes"`
	}

	var req CreateAnimeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error: %v", err)
		util.WriteJson(w, http.StatusBadRequest, util.Envelope{
			"message": "bad request",
		})
		return
	}

	name := strings.TrimSpace(req.Name)
	episodes := req.Episodes

	if name == "" {
		log.Printf("error: %v", "missing name")
		util.WriteJson(w, http.StatusBadRequest, util.Envelope{
			"message": "bad request",
		})
		return
	}

	if episodes <= 0 {
		log.Printf("error: %v", "episodes <= 0")
		util.WriteJson(w, http.StatusBadRequest, util.Envelope{
			"message": "bad request",
		})
		return
	}

	reqAnime := &database.Anime{
		Name:     name,
		Episodes: episodes,
	}
	dbAnime, err := h.animeService.CreateAnime(reqAnime)
	if err != nil {
		log.Printf("error: %v", err)
		util.WriteJson(w, http.StatusInternalServerError, util.Envelope{
			"message": "internal server error",
		})
		return
	}

	util.WriteJson(w, http.StatusCreated, util.Envelope{
		"anime": dbAnime,
	})
}

func (h *handlerAnime) GetAnime(w http.ResponseWriter, r *http.Request) {
	type GetAnimeRequest struct {
		Id string `json:"id"`
	}

	var req GetAnimeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error: %v", err)
		util.WriteJson(w, http.StatusBadRequest, util.Envelope{
			"message": "bad request",
		})
		return
	}

	if err := uuid.Validate(req.Id); err != nil {
		log.Printf("error: %v", err)
		util.WriteJson(w, http.StatusBadRequest, util.Envelope{
			"message": "bad request",
		})
		return
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("error: %v", err)
		util.WriteJson(w, http.StatusBadRequest, util.Envelope{
			"message": "bad request",
		})
		return
	}

	reqAnime := &database.Anime{
		Id: id,
	}
	dbAnime, err := h.animeService.GetAnimeById(reqAnime)
	if err != nil {
		log.Printf("error: %v", err)
		util.WriteJson(w, http.StatusInternalServerError, util.Envelope{
			"message": "internal server error",
		})
		return
	}

	util.WriteJson(w, http.StatusOK, util.Envelope{
		"anime": dbAnime,
	})
}

func (h *handlerAnime) UpdateAnime(w http.ResponseWriter, r *http.Request) {
	type UpdateAnimeRequest struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Episodes int    `json:"episodes"`
	}

	var req UpdateAnimeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error: %v", err)
		util.WriteJson(w, http.StatusBadRequest, util.Envelope{
			"message": "bad request",
		})
		return
	}

	if err := uuid.Validate(req.Id); err != nil {
		log.Printf("error: %v", err)
		util.WriteJson(w, http.StatusBadRequest, util.Envelope{
			"message": "bad request",
		})
		return
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("error: %v", err)
		util.WriteJson(w, http.StatusBadRequest, util.Envelope{
			"message": "bad request",
		})
		return
	}

	name := strings.TrimSpace(req.Name)
	episodes := req.Episodes

	if name == "" {
		log.Printf("error: %v", "missing name")
		util.WriteJson(w, http.StatusBadRequest, util.Envelope{
			"message": "bad request",
		})
		return
	}

	if episodes <= 0 {
		log.Printf("error: %v", "episodes <= 0")
		util.WriteJson(w, http.StatusBadRequest, util.Envelope{
			"message": "bad request",
		})
		return
	}

	reqAnime := &database.Anime{
		Id:       id,
		Name:     name,
		Episodes: episodes,
	}
	dbAnime, err := h.animeService.UpdateAnime(reqAnime)
	if err != nil {
		log.Printf("error: %v", err)
		util.WriteJson(w, http.StatusInternalServerError, util.Envelope{
			"message": "internal server error",
		})
		return
	}

	util.WriteJson(w, http.StatusOK, util.Envelope{
		"anime": dbAnime,
	})
}

func (h *handlerAnime) DeleteAnime(w http.ResponseWriter, r *http.Request) {
	type DeleteAnimeRequest struct {
		Id string `json:"id"`
	}

	var req DeleteAnimeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error: %v", err)
		util.WriteJson(w, http.StatusBadRequest, util.Envelope{
			"message": "bad request",
		})
		return
	}

	if err := uuid.Validate(req.Id); err != nil {
		log.Printf("error: %v", err)
		util.WriteJson(w, http.StatusBadRequest, util.Envelope{
			"message": "bad request",
		})
		return
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("error: %v", err)
		util.WriteJson(w, http.StatusBadRequest, util.Envelope{
			"message": "bad request",
		})
		return
	}

	reqAnime := &database.Anime{
		Id: id,
	}
	if err := h.animeService.DeleteAnimeById(reqAnime); err != nil {
		log.Printf("error: %v", err)
		util.WriteJson(w, http.StatusInternalServerError, util.Envelope{
			"message": "internal server error",
		})
		return
	}

	util.WriteJson(w, http.StatusNoContent, util.Envelope{})
}
