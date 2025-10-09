package handlers

import (
	"encoding/json"
	"net/http"
	"service-user/internal/signer"
)

type HandlerSigner interface {
	GetJwks(w http.ResponseWriter, r *http.Request)
}

type handlerSigner struct {
	signer signer.Signer
}

var handlerSignerInstance *handlerSigner

func NewHandlerSigner(signer signer.Signer) HandlerSigner {
	if handlerSignerInstance != nil {
		return handlerSignerInstance
	}

	newHandlerSigner := &handlerSigner{
		signer: signer,
	}
	handlerSignerInstance = newHandlerSigner

	return handlerSignerInstance
}

func (h *handlerSigner) GetJwks(w http.ResponseWriter, r *http.Request) {
	set := h.signer.GetJwkSet()

	js, err := json.MarshalIndent(set, "", " ")
	if err != nil {
		return
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(js)
}
