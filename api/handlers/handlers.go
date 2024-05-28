package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func errorLog(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	slog.Error(msg, "status", status)
	json.NewEncoder(w).Encode(msg)
}

func jsonResponse(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
