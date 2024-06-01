package handlers

import (
	"net/http"

	"github.com/condemo/raspi-home-service/api/util"
	"github.com/condemo/raspi-home-service/store"
	"github.com/condemo/raspi-home-service/types"
)

type UserHandler struct {
	store store.Store
}

func NewUserHandler(s store.Store) *UserHandler {
	return &UserHandler{store: s}
}

func (h *UserHandler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("POST /login", h.loginHandler)
	r.HandleFunc("POST /signup", h.signupHandler)
}

func (h *UserHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	pass := r.FormValue("password")

	user, err := h.store.GetUserByUsername(username)
	if err != nil {
		ErrorLog(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if !util.VerifyPass(user.Password, pass) {
		ErrorLog(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := util.CreateJWT(user.ID)
	if err != nil {
		ErrorLog(w, http.StatusInternalServerError, "internal server error")
	}

	JsonResponse(w, http.StatusOK, map[string]string{
		"access_token": token,
		"type":         "bearer",
	})
}

func (h *UserHandler) signupHandler(w http.ResponseWriter, r *http.Request) {
	un := r.FormValue("username")
	pass, err := util.EncryptPass(r.FormValue("password"))
	if err != nil {
		ErrorLog(w, http.StatusInternalServerError, "internal server error")
		return
	}

	user := &types.User{
		Username: un,
		Password: pass,
	}

	if ok := user.Validate(); !ok {
		ErrorLog(w, http.StatusBadRequest, "bad request")
		return
	}

	err = h.store.CreateUser(user)
	if err != nil {
		ErrorLog(w, http.StatusInternalServerError, "internal server error")
		return
	}

	JsonResponse(w, http.StatusCreated, user)
}
