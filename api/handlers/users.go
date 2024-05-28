package handlers

import (
	"net/http"

	"github.com/condemo/raspi-test/api/util"
	"github.com/condemo/raspi-test/store"
	"github.com/condemo/raspi-test/types"
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
		errorLog(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if !util.VerifyPass(user.Password, pass) {
		errorLog(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := util.CreateJWT(user.ID)
	if err != nil {
		errorLog(w, http.StatusInternalServerError, "internal server error")
	}

	jsonResponse(w, http.StatusOK, map[string]string{
		"access_token": token,
		"type":         "bearer",
	})
}

func (h *UserHandler) signupHandler(w http.ResponseWriter, r *http.Request) {
	un := r.FormValue("username")
	pass, err := util.EncryptPass(r.FormValue("password"))
	if err != nil {
		errorLog(w, http.StatusInternalServerError, "internal server error")
	}

	user := &types.User{
		Username: un,
		Password: pass,
	}

	if ok := user.Validate(); !ok {
		errorLog(w, http.StatusBadRequest, "bad request")
		return
	}

	err = h.store.CreateUser(user)
	if err != nil {
		errorLog(w, http.StatusInternalServerError, "internal server error")
		return
	}

	jsonResponse(w, http.StatusCreated, user)
}
