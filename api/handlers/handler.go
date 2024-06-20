package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

func RenderTempl(w http.ResponseWriter, r *http.Request, c templ.Component) {
	c.Render(r.Context(), w)
}

func HTTPSendHTML(w http.ResponseWriter, c templ.Component) {
	tmpl, err := templ.ToGoHTML(context.Background(), c)
	if err != nil {
		ErrorLog(w, http.StatusInternalServerError, "internal server err")
	}

	fmt.Fprint(w, tmpl)
}

func ErrorLog(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	slog.Error(msg, "status", status)
	fmt.Fprint(w, msg)
}

func JsonResponse(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func TextResonse(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	fmt.Fprint(w, msg)
}
