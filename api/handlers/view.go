package handlers

import (
	"net/http"
	"strconv"

	"github.com/condemo/raspi-home-service/config"
	"github.com/condemo/raspi-home-service/store"
	"github.com/condemo/raspi-home-service/tools"
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
	r.HandleFunc("POST /theme", h.themeChange)
}

func (h *ViewHandler) homeHandler(w http.ResponseWriter, r *http.Request) {
	sysInfo := tools.NewSysInfo()
	RenderTempl(w, r, core.Home(sysInfo))
}

func (h *ViewHandler) themeChange(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("t")
	i, _ := strconv.Atoi(t)

	config.UIConf.ChangeTheme(i)

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}
