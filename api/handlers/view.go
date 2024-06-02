package handlers

import (
	"net/http"

	"github.com/condemo/raspi-home-service/store"
	"github.com/condemo/raspi-home-service/views/core"
)

type ViewHandler struct {
	store store.Store
}

func NewViewHanlder(s store.Store) *ViewHandler {
	return &ViewHandler{store: s}
}

func (h *ViewHandler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /", h.homeHandler)
}

func (h *ViewHandler) homeHandler(w http.ResponseWriter, r *http.Request) {
	RenderTempl(w, r, core.Home())
}
