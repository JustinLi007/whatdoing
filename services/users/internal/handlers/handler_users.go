package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	libutils "github.com/JustinLi007/whatdoing/libs/go/util"
	"github.com/JustinLi007/whatdoing/services/users/internal/database"
	"github.com/JustinLi007/whatdoing/services/users/internal/password"
	"github.com/JustinLi007/whatdoing/services/users/internal/signer"
	"github.com/JustinLi007/whatdoing/services/users/internal/token"
	"github.com/JustinLi007/whatdoing/services/users/internal/utils"
)

type HandlerUsers interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
}

type handlerUsers struct {
	signer      signer.Signer
	userService database.ServiceUsers
}

var handlerUsersInstance *handlerUsers

func NewHandlerUsers(signer signer.Signer, userService database.ServiceUsers) HandlerUsers {
	if handlerUsersInstance != nil {
		return handlerUsersInstance
	}
	newHandlerUsers := &handlerUsers{
		signer:      signer,
		userService: userService,
	}
	handlerUsersInstance = newHandlerUsers
	return handlerUsersInstance
}

func (h *handlerUsers) SignUp(w http.ResponseWriter, r *http.Request) {
	type SignUpRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		libutils.WriteJson(w, http.StatusBadRequest, libutils.Envelope{})
		return
	}

	if !utils.IsValidEmail(req.Email) {
		libutils.WriteJson(w, http.StatusBadRequest, libutils.Envelope{})
		return
	}

	if !utils.IsValidPassword(req.Password) {
		libutils.WriteJson(w, http.StatusBadRequest, libutils.Envelope{})
		return
	}

	reqUser := &database.User{
		Email:        req.Email,
		Password:     &password.Password{},
		RefreshToken: token.NewToken(token.REFRESH_TOKEN_TTL),
	}
	reqUser.Password.Set(req.Password)

	dbUser, err := h.userService.CreateUser(reqUser)
	if err != nil {
		libutils.WriteJson(w, http.StatusInternalServerError, libutils.Envelope{})
		return
	}

	jwt, err := h.signer.NewJwt(dbUser.Id.String(), "", time.Hour)
	if err != nil {
		libutils.WriteJson(w, http.StatusInternalServerError, libutils.Envelope{})
		return
	}

	libutils.SetCookie(w, "jwt", jwt)
	libutils.SetCookie(w, "refresh-token", dbUser.RefreshToken.GetPlainText())
	libutils.WriteJson(w, http.StatusCreated, libutils.Envelope{
		"user": dbUser,
	})
}

func (h *handlerUsers) Login(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		libutils.WriteJson(w, http.StatusBadRequest, libutils.Envelope{})
		return
	}

	if !utils.IsValidEmail(req.Email) {
		libutils.WriteJson(w, http.StatusBadRequest, libutils.Envelope{})
		return
	}

	if !utils.IsValidPassword(req.Password) {
		libutils.WriteJson(w, http.StatusBadRequest, libutils.Envelope{})
		return
	}

	reqUser := &database.User{
		Email:    req.Email,
		Password: &password.Password{},
	}
	reqUser.Password.Set(req.Password)

	//TODO: user should have scope?
	dbUser, err := h.userService.GetUserByEmailPassword(reqUser)
	if err != nil {
		libutils.WriteJson(w, http.StatusInternalServerError, libutils.Envelope{})
		return
	}

	jwt, err := h.signer.NewJwt(dbUser.Id.String(), "ohfk,", time.Hour)
	if err != nil {
		libutils.WriteJson(w, http.StatusInternalServerError, libutils.Envelope{})
		return
	}

	libutils.SetCookie(w, "jwt", jwt)
	libutils.SetCookie(w, "refresh-token", dbUser.RefreshToken.GetPlainText())
	libutils.WriteJson(w, http.StatusOK, libutils.Envelope{
		"user": dbUser,
	})
}

func (h *handlerUsers) Logout(w http.ResponseWriter, r *http.Request) {
	libutils.WriteJson(w, http.StatusNotImplemented, libutils.Envelope{
		"message": "logout not yet implemented",
	})
}

func (h *handlerUsers) Refresh(w http.ResponseWriter, r *http.Request) {
	libutils.WriteJson(w, http.StatusNotImplemented, libutils.Envelope{
		"message": "refresh not yet implemented",
	})
}
