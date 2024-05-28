package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

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
