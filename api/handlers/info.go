package handlers

import (
	"net/http"

	"github.com/condemo/raspi-test/store"
)

type InfoHandler struct {
	store store.Store
}

func NewInfoHandler(s store.Store) *InfoHandler {
	return &InfoHandler{store: s}
}

func (h *InfoHandler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /", h.homeHandler)
}

func (h *InfoHandler) homeHandler(w http.ResponseWriter, r *http.Request) {
	TextResonse(w, http.StatusOK, "Info route")
}
