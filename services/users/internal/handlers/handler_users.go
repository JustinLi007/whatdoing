package handlers

import (
	"encoding/json"
	"net/http"
	"service-user/internal/database"
	"service-user/internal/password"
	"service-user/internal/signer"
	"service-user/internal/token"
	"service-user/internal/utils"
	"time"
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
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{})
		return
	}

	if !utils.IsValidEmail(req.Email) {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{})
		return
	}

	if !utils.IsValidPassword(req.Password) {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{})
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
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{})
		return
	}

	jwt, err := h.signer.NewJwt(dbUser.Id.String(), "", time.Hour)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{})
		return
	}

	utils.SetCookie(w, "jwt", jwt)
	utils.SetCookie(w, "refresh-token", dbUser.RefreshToken.GetPlainText())
	utils.WriteJson(w, http.StatusCreated, utils.Envelope{
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
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{})
		return
	}

	if !utils.IsValidEmail(req.Email) {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{})
		return
	}

	if !utils.IsValidPassword(req.Password) {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{})
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
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{})
		return
	}

	jwt, err := h.signer.NewJwt(dbUser.Id.String(), "ohfk,", time.Hour)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{})
		return
	}

	utils.SetCookie(w, "jwt", jwt)
	utils.SetCookie(w, "refresh-token", dbUser.RefreshToken.GetPlainText())
	utils.WriteJson(w, http.StatusOK, utils.Envelope{
		"user": dbUser,
	})
}

func (h *handlerUsers) Logout(w http.ResponseWriter, r *http.Request) {
	utils.WriteJson(w, http.StatusNotImplemented, utils.Envelope{
		"message": "logout not yet implemented",
	})
}

func (h *handlerUsers) Refresh(w http.ResponseWriter, r *http.Request) {
	utils.WriteJson(w, http.StatusNotImplemented, utils.Envelope{
		"message": "refresh not yet implemented",
	})
}
